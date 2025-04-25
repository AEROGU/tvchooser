package tvchooser

import (
	"os"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/AEROGU/tvchooser/mounted"
	"github.com/AEROGU/tvchooser/tvclang"
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

const (
	expandedPrefix  = "▼ "
	collapsedPrefix = "► "
)

type nodeInfo struct {
	Path     string
	IsRoot   bool
	IsCustom bool // Is a custom path with no path, just for hold custom shortcuts
	IsFinal  bool // Is the last node in the path, if true, onNodeSelected is not called anymore on this node
}

type directoryView struct {
	dirView          *tview.TreeView // Internal tree view
	fileView         *tview.List     // Reference to update when a directory is selected
	selectedPath     string          // Last selected directory path
	showHidden       bool
	textViewToUpdate *tview.TextView
	onSelectedFunc   func(node *tview.TreeNode)
	fastAccessPaths  []string // Paths to be added to the tree view as fast access paths
}

// A helper function which adds the directories of the given path
// to the given target node.
//
// Parameters:
// - target: The target node to which child nodes are added.
// - path: The path of the directory to process.
func (dv *directoryView) addChildDirectories(target *tview.TreeNode, path string) {
	files, err := os.ReadDir(path)
	if err != nil {
		txt := target.GetText()
		target.SetText(txt + " - " + tvclang.GetTranslations().AccessDenied)

		// Avoid execution of most of onNodeSelected on this node by setting IsFinal to true
		{
			info := target.GetReference().(nodeInfo)
			info.IsFinal = true
			target.SetReference(info)
		}

		target.SetColor(tcell.ColorRed)
		// target.SetSelectable(false)
		return
	}

	for _, file := range files {
		if !file.IsDir() {
			continue // Skip files
		}

		firstDirectoryCharacter := file.Name()[0]
		if !dv.showHidden && (firstDirectoryCharacter == '.' || firstDirectoryCharacter == '~' || firstDirectoryCharacter == '$') {
			continue
		}

		node := tview.NewTreeNode(collapsedPrefix + file.Name())
		node.SetReference(nodeInfo{
			Path:     filepath.Join(path, file.Name()),
			IsRoot:   false,
			IsCustom: false,
		})

		// node.SetColor(tcell.ColorGreen)
		node.SetColor(tcell.ColorTeal)
		node.SetIndent(1)

		target.AddChild(node)
	}
}

// onNodeSelected handles the selection of a tree node.
//
// It takes a pointer to a tview.TreeNode as a parameter.
// The function retrieves the reference of the node and typecasts it to nodeInfo.
// It updates the selectedPath field of the directoryView.
// If the node is marked as final, the function returns early.
// The function then retrieves the text of the node and checks if it has any children.
// If the node has no children and is not a custom node, the function loads and shows files in the directory.
// If there are no directories left, the function removes the arrow at the beginning.
// If the node is not a custom node, the function updates the selectedPath field of the directoryView.
// If the node is a custom node, the function clears the selectedPath field.
// Finally, the function updates the text of the node based on its expansion state (Updates the arrow).
//
// Parameters:
// - node: The selected tree node.
//
// Returns: None.
func (dv *directoryView) onNodeSelected(node *tview.TreeNode) {
	info := node.GetReference().(nodeInfo)

	if !info.IsCustom {
		// Sets dv.selectedPath to info.Path, if it has a trailing path separator do not add it.
		if strings.HasSuffix(info.Path, string(os.PathSeparator)) {
			dv.selectedPath = info.Path
		} else {
			if info.Path != "" {
				dv.selectedPath = info.Path + string(os.PathSeparator)
			} else {
				dv.selectedPath = info.Path
			}
		}
	} else {
		dv.selectedPath = ""
	}

	if dv.textViewToUpdate != nil {
		go dv.textViewToUpdate.SetText(dv.selectedPath)
	}

	if dv.onSelectedFunc != nil {
		dv.onSelectedFunc(node)
	}

	if info.IsFinal {
		return
	}

	txt := node.GetText()

	children := node.GetChildren()
	if len(children) == 0 && !info.IsCustom {
		// Load and show files in this directory.
		path := info.Path
		dv.addChildDirectories(node, path)
		txt = node.GetText() // Get the new text here, in case it was changed in the add function.

		// If no directories left, erase the arrot at the beginning. (Only for not custom directories)
		if len(node.GetChildren()) == 0 && !info.IsCustom {
			if strings.HasPrefix(txt, collapsedPrefix) {
				txt = strings.TrimLeft(txt, collapsedPrefix)
			} else {
				txt = strings.TrimLeft(txt, expandedPrefix)
			}
			node.SetText("  " + txt)

			info.IsFinal = true
			node.SetReference(info)

			return
		}

	} else {
		// Collapse if visible, expand if collapsed.
		node.SetExpanded(!node.IsExpanded())
	}

	if strings.HasPrefix(txt, collapsedPrefix) {
		txt = strings.TrimLeft(txt, collapsedPrefix)
	} else {
		txt = strings.TrimLeft(txt, expandedPrefix)
	}

	if node.IsExpanded() {
		node.SetText(expandedPrefix + txt)
	} else {
		node.SetText(collapsedPrefix + txt)
	}
}

// newDirectoryView creates a new directory view with the specified showHidden flag.
//
// Parameters:
// - showHidden: A boolean indicating whether to show hidden directories.
//
// Returns:
// - A pointer to the created directoryView.
func newDirectoryView(showHidden bool, textViewToUpdate *tview.TextView, onSelectedFunc func(node *tview.TreeNode), fastAccessPaths []string) *directoryView {
	tree := tview.NewTreeView()

	// Add rootNode node.
	rootNode := tview.NewTreeNode(expandedPrefix + tvclang.GetTranslations().ThisPC).SetColor(tcell.ColorWhite).SetIndent(0)
	rootNode.SetReference(nodeInfo{
		Path:     "",
		IsRoot:   true,
		IsCustom: true,
	})

	// Add favorites node if any.
	if len(fastAccessPaths) > 0 {
		favoritesNode := addChildNode(rootNode, tvclang.GetTranslations().Favorites, true, nodeInfo{
			Path:     "",
			IsRoot:   true,
			IsCustom: true,
		})
		for _, path := range fastAccessPaths {

			curFastAccessDirPath, curFastAccessDirName := splitPathAndName(path, '|')
			if curFastAccessDirPath == "" { // Check if the path is empty (just in case).
				continue // Skip empty paths.
			}

			// curFastAccessDirName := filepath.Base(path)
			// if curFastAccessDirName == "" { // Check if the directory name is empty (just in case).
			// 	continue // Skip empty directory names.
			// }

			// Add the path to the favorites node.
			addChildNode(favoritesNode, curFastAccessDirName, false, nodeInfo{
				Path:     curFastAccessDirPath,
				IsRoot:   false,
				IsCustom: false,
			})
		}
	}

	// Add userprofile node.
	userHomeDir, _ := os.UserHomeDir()
	addChildNode(rootNode, tvclang.GetTranslations().HomeDir, false, nodeInfo{
		Path:     userHomeDir,
		IsRoot:   true,
		IsCustom: false,
	})

	// Add devices node.
	devicesNode := addChildNode(rootNode, tvclang.GetTranslations().Devices, true, nodeInfo{
		Path:     "",
		IsRoot:   true,
		IsCustom: true,
	})

	if runtime.GOOS == "windows" {
		devices, err := mounted.GetWindowsDriveLetters()
		if err != nil {
			devicesNode.SetColor(tcell.ColorRed)
			devicesNode.SetSelectable(false)
		} else {
			for _, drive := range devices {
				driveRoot := drive + ":" + string(os.PathSeparator)
				addChildNode(devicesNode, driveRoot, false, nodeInfo{
					Path:     driveRoot,
					IsRoot:   true,
					IsCustom: false,
				})

			}
		}
	}

	tree.SetRoot(rootNode).SetCurrentNode(rootNode)

	dv := &directoryView{
		dirView:          tree,
		showHidden:       showHidden,
		textViewToUpdate: textViewToUpdate,
		onSelectedFunc:   onSelectedFunc,
		fastAccessPaths:  fastAccessPaths,
	}

	// If a directory was selected, open it.
	tree.SetSelectedFunc(dv.onNodeSelected)

	return dv
}

func addChildNode(rootNode *tview.TreeNode, nodeName string, expanded bool, info nodeInfo) *tview.TreeNode {
	var prefix string
	if expanded {
		prefix = expandedPrefix
	} else {
		prefix = collapsedPrefix
	}

	newNode := tview.NewTreeNode(prefix + nodeName).SetReference(info)
	rootNode.AddChild(newNode)
	return newNode
}
