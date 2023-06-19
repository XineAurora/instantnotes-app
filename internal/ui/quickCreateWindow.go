package ui

import (
	"path/filepath"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"github.com/XineAurora/instantnotes-app/internal/api"
	"github.com/XineAurora/instantnotes-app/internal/types"
)

type QuickCreateWindow struct {
	Window     fyne.Window
	api        *api.ApiConnector
	saveFolder *types.Folder
	folderPath string

	title   *widget.Entry
	content *widget.Entry

	saveFoldButton *widget.Button
}

func NewQuickCreateWindow(app fyne.App, api *api.ApiConnector) *QuickCreateWindow {
	q := QuickCreateWindow{
		Window:     app.NewWindow("Create Note"),
		api:        api,
		saveFolder: &types.Folder{},
		folderPath: "\\",
	}
	q.Window.SetContent(q.initUi())
	q.Window.Resize(fyne.NewSize(800, 600))
	q.Window.SetCloseIntercept(func() {
		q.Window.Hide()
	})
	return &q
}

func (w *QuickCreateWindow) initUi() fyne.CanvasObject {
	w.title = widget.NewEntry()
	w.title.SetPlaceHolder("Type a Title")
	w.content = widget.NewMultiLineEntry()
	w.content.SetPlaceHolder("Type a Note here")
	main := container.New(layout.NewBorderLayout(w.title, nil, nil, nil), w.title, w.content)

	// open dialog with folders, by default use root, otherwise previous location
	w.saveFoldButton = widget.NewButton(w.folderPath, func() {
		w.chooseFolderDial().Show()
	})
	// if note created, hide window
	saveButton := widget.NewButton("Save", w.saveNote)
	clearButton := widget.NewButton("Clear", w.clear)
	saveArea := container.New(layout.NewBorderLayout(nil, nil, clearButton, saveButton), clearButton, w.saveFoldButton, saveButton)

	return container.New(layout.NewBorderLayout(nil, saveArea, nil, nil), main, saveArea)
}

func (w *QuickCreateWindow) saveNote() {
	if _, err := w.api.CreateNote(w.title.Text, w.content.Text, w.saveFolder.ID, w.saveFolder.GroupId); err != nil {
		//show popup
		widget.NewPopUp(widget.NewLabel("Cannot create note"), nil).Show()
	} else {
		w.clear()
		w.Window.Hide()
	}
}

func (w *QuickCreateWindow) clear() {
	w.title.SetText("")
	w.content.SetText("")
}

func (w *QuickCreateWindow) chooseFolderDial() dialog.Dialog {
	folderContent := container.NewVBox()
	scrollArea := container.NewVScroll(folderContent)
	scrollArea.SetMinSize(fyne.NewSize(w.Window.Canvas().Size().Width/3, w.Window.Canvas().Size().Height/2))
	path := widget.NewLabel("")
	tmpFolder := w.saveFolder
	tmpPath := w.folderPath
	var loadFolder func(id uint)
	loadFolder = func(id uint) {
		//clear area
		folderContent.RemoveAll()
		//get loading folder info
		curFolder, err := w.api.GetFolderInfo(id)
		if err != nil {
			folderContent.Add(widget.NewLabel("Not Logged In"))
			return
		}
		//set this as save folder
		tmpFolder = &curFolder
		//add nav back button
		if id != 0 {
			b := widget.NewButtonWithIcon("...", theme.NavigateBackIcon(), func() {
				pid, err := w.api.GetParentFolder(id)
				if err != nil {
					return
				}
				tmpPath = filepath.Dir(tmpPath)
				loadFolder(pid.ParentFolderID)
			})
			folderContent.Add(b)
		}
		//get all folders inside loaded
		folders, _, _ := w.api.GetFolderContent(id)
		// add non-group folders if it's root folder, add all otherwise
		for _, f := range folders {
			if f.GroupId == 0 || tmpFolder.GroupId != 0 {
				thisFolder := f
				b := widget.NewButtonWithIcon(thisFolder.Name, theme.FolderIcon(), func() {
					tmpPath = filepath.Join(tmpPath, thisFolder.Name)
					loadFolder(thisFolder.ID)
				})
				folderContent.Add(b)
			}
		}
		// add section of group folder if it's root folder
		if tmpFolder.GroupId == 0 && tmpFolder.ID == 0 {
			folderContent.Add(widget.NewLabel("Groups"))
			for _, f := range folders {
				if f.GroupId != 0 {
					thisFolder := f
					b := widget.NewButtonWithIcon(thisFolder.Name, theme.FolderIcon(), func() {
						tmpPath = filepath.Join(tmpPath, thisFolder.Name)
						loadFolder(thisFolder.ID)
					})
					folderContent.Add(b)
				}
			}
		}
	}
	loadFolder(w.saveFolder.ID)

	canv := container.New(layout.NewBorderLayout(path, nil, nil, nil), path, scrollArea)

	return dialog.NewCustomConfirm("Choose folder", "Confirm", "Cancel", canv, func(b bool) {
		if b {
			w.saveFolder = tmpFolder
			w.folderPath = tmpPath
			w.saveFoldButton.SetText(tmpPath)
		}
	}, w.Window)
}
