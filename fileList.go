package tvchooser

import (
	"fmt"
	"os"

	"github.com/AEROGU/tvchooser/tvclang"
	"github.com/dustin/go-humanize"
	"github.com/rivo/tview"
)

type fileView struct {
	fileList         *tview.List
	directoryView    *directoryView
	rootPath         string
	selectedFileName string
	showHidden       bool
	textViewToUpdate *tview.TextView
}

// updatePath updates the root path of the fileView and refreshes the file list displayed.
// It clears the current file list and populates it with files from the new path,
// filtering out hidden files if the showHidden flag is false.
//
// Parameters:
//   - newPath: The new root path to set for the fileView.
//
// Behavior:
//   - If the newPath is the same as the current rootPath, the function returns immediately.
//   - If the newPath is empty, the function clears the file list and returns.
//   - Reads the directory at the newPath and iterates through its files:
//   - Skips files that are hidden (starting with '.', '~', or '$') if showHidden is false.
//   - Retrieves file information (modification time and size) and formats it for display.
//   - Adds each file to the file list with its name and formatted information.
//   - Sets a callback function to handle file selection events.
//
// Errors:
//   - If an error occurs while reading the directory or retrieving file information,
//     the file is added to the list without additional details.
func (fv *fileView) updatePath(newPath string) {
	if newPath == fv.rootPath {
		return
	}

	fv.rootPath = newPath
	fv.fileList.Clear()

	if newPath == "" {
		return
	}

	files, err := os.ReadDir(fv.rootPath)
	if err == nil {
		for _, file := range files {
			if !file.IsDir() {

				firstFileCharacter := file.Name()[0]
				if !fv.showHidden && (firstFileCharacter == '.' || firstFileCharacter == '~' || firstFileCharacter == '$') {
					continue
				}

				fileInfo, err := file.Info()

				if err == nil {
					fInf := fmt.Sprintf("  %s: %s", tvclang.GetTranslations().Modfied, fileInfo.ModTime().Format("2006-01-02 15:04:05")) // Modification date
					fInf += fmt.Sprintf(" | %s: %s", tvclang.GetTranslations().Size, humanize.Bytes(uint64(fileInfo.Size())))            // Size
					fv.fileList.AddItem(fileInfo.Name(), fInf, 0, nil)
				} else {
					fv.fileList.AddItem(file.Name(), "", 0, nil)
				}

				fv.fileList.SetSelectedFunc(func(int, string, string, rune) {
					go fv.onSelectedFunc()
				})
			}
		}
	}
}

// onSelectedFunc is a method of the fileView struct that handles the event
// when an item is selected in the file list. It retrieves the currently
// selected item's text from the file list and updates the selectedFileName
// field. If a valid file name is selected, it constructs the full file path
// by combining the selected directory path and the file name, and updates
// the associated text view with this path.
func (fv *fileView) onSelectedFunc() {
	curItem := fv.fileList.GetCurrentItem()
	fv.selectedFileName, _ = fv.fileList.GetItemText(curItem)
	if fv.selectedFileName != "" {
		path := fv.directoryView.selectedPath + fv.selectedFileName
		fv.textViewToUpdate.SetText(path)
	}
}

func newFileView(rootPath string, showHidden bool, textViewToUpdate *tview.TextView, diredirectoryView *directoryView) *fileView {
	fv := &fileView{
		fileList:         tview.NewList(),
		rootPath:         rootPath,
		showHidden:       showHidden,
		textViewToUpdate: textViewToUpdate,
		directoryView:    diredirectoryView,
	}

	fv.updatePath(rootPath)

	return fv
}
