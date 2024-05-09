package tvchooser

import (
	"fmt"
	"os"

	"github.com/aerogu/tvchooser/tvclang"
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
