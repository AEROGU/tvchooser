package tvchooser

import (
	"github.com/aerogu/tvchooser/tvclang"
	"github.com/rivo/tview"
)

func FileChooser(parentApp *tview.Application) string {
	selectedPath := ""

	app := tview.NewApplication()
	runApp := func() {
		if err := app.Run(); err != nil {
			panic(err)
		}
	}

	selectedPathView := tview.NewTextView()
	selectedPathView.SetBorder(true)

	dirView := newDirectoryView(false, selectedPathView, nil)
	fileView := newFileView("", dirView.showHidden, selectedPathView, dirView)
	dirView.onSelectedFunc = func(node *tview.TreeNode) {
		fileView.updatePath(node.GetReference().(nodeInfo).Path)
	}
	selectionPanel := tview.NewFlex().SetDirection(tview.FlexColumn).AddItem(dirView.dirView, 0, 1, true).AddItem(fileView.fileList, 0, 1, false)

	buttonsView := tview.NewForm()
	buttonsView.SetButtonsAlign(tview.AlignRight)
	// Cancel button
	buttonsView.AddButton(tvclang.GetTranslations().Cancel, func() {
		selectedPath = ""
		app.Stop()
		buttonsView.SetButtonsAlign(tview.AlignRight)
	})
	// Accept button
	buttonsView.AddButton(tvclang.GetTranslations().Accept, func() {
		selectedPath = dirView.selectedPath + fileView.selectedFileName
		app.Stop()
	})

	rootPanel := tview.NewFlex().SetDirection(tview.FlexRow)
	rootPanel.AddItem(selectedPathView, 3, 0, false)
	rootPanel.AddItem(selectionPanel, 0, 1, true)
	rootPanel.AddItem(buttonsView, 3, 0, false)

	app.SetRoot(rootPanel, true).EnableMouse(true).EnablePaste(true)
	if parentApp != nil {
		parentApp.Suspend(func() {
			runApp()
		})
	} else {
		runApp()
	}

	return selectedPath
}

func DirectoryChooser(parentApp *tview.Application) string {
	selectedPath := ""

	app := tview.NewApplication()
	runApp := func() {
		if err := app.Run(); err != nil {
			panic(err)
		}
	}

	selectedPathView := tview.NewTextView()
	selectedPathView.SetBorder(true)

	dirView := newDirectoryView(false, selectedPathView, nil)
	selectionPanel := tview.NewFlex().SetDirection(tview.FlexColumn).AddItem(dirView.dirView, 0, 2, true)

	buttonsView := tview.NewForm()
	buttonsView.SetButtonsAlign(tview.AlignRight)
	// Cancel button
	buttonsView.AddButton(tvclang.GetTranslations().Cancel, func() {
		selectedPath = ""
		app.Stop()
	})
	// Accept button
	buttonsView.AddButton(tvclang.GetTranslations().Accept, func() {
		selectedPath = dirView.selectedPath
		app.Stop()
	})

	rootPanel := tview.NewFlex().SetDirection(tview.FlexRow)
	rootPanel.AddItem(selectedPathView, 3, 0, false)
	rootPanel.AddItem(selectionPanel, 0, 1, true)
	rootPanel.AddItem(buttonsView, 3, 0, false)

	app.SetRoot(rootPanel, true).EnableMouse(true).EnablePaste(true)
	if parentApp != nil {
		parentApp.Suspend(func() {
			runApp()
		})
	} else {
		runApp()
	}

	return selectedPath
}
