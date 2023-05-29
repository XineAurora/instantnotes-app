package ui

import (
	"log"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/driver/desktop"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"github.com/XineAurora/instantnotes-app/internal/api"
	"github.com/XineAurora/instantnotes-app/internal/application"
	"github.com/XineAurora/instantnotes-app/internal/types"
)

type MainWindow struct {
	Window fyne.Window

	currentNote   *types.Note
	currentFolder *types.Folder
	folderContent *fyne.Container
	noteTitle     *widget.Entry
	noteContent   *widget.Entry

	api *api.ApiConnector
}

func NewMainWindow(a *application.Application) MainWindow {
	mw := MainWindow{}
	mw.api = a.Api
	mw.Window = a.App.NewWindow("Instant Notes")
	mw.Window.SetContent(mw.initUi())
	mw.Window.Resize(fyne.NewSize(600, 400))
	//setup shortcuts
	mw.setupShortcuts()
	mw.currentNote = &types.Note{ID: 0}
	mw.currentFolder = &types.Folder{ID: 0}
	return mw
}

func (m *MainWindow) initUi() fyne.CanvasObject {
	//init right side

	m.noteTitle = widget.NewEntry()
	m.noteContent = widget.NewMultiLineEntry()
	noteToolbar := widget.NewToolbar(
		//save current note
		widget.NewToolbarAction(theme.DocumentSaveIcon(), m.saveNote),
		//delete current note
		widget.NewToolbarAction(theme.ContentRemoveIcon(), m.deleteNote),

		//speech to text
		widget.NewToolbarAction(theme.MediaRecordIcon(), func() {}),
	)
	content := container.New(layout.NewBorderLayout(noteToolbar, nil, nil, nil), noteToolbar)
	content.Add(container.New(layout.NewBorderLayout(m.noteTitle, nil, nil, nil), m.noteTitle, m.noteContent))

	//init left side

	m.folderContent = container.NewVBox()
	scrollArea := container.NewVScroll(m.folderContent)

	//fill list

	toolBar := widget.NewToolbar(
		//create new note
		widget.NewToolbarAction(theme.ContentAddIcon(), m.createNote),
		//create new Folder
		widget.NewToolbarAction(theme.FolderNewIcon(), m.createFolder),

		//refresh folder content
		widget.NewToolbarAction(theme.ViewRefreshIcon(), func() { go m.updateFolderContent(m.currentFolder.ID) }),
	)

	side := container.New(layout.NewBorderLayout(toolBar, nil, nil, nil), toolBar, scrollArea)

	split := container.NewHSplit(side, content)
	split.Offset = 0.25
	return split
}

func (m *MainWindow) setupShortcuts() {
	ctrlS := desktop.CustomShortcut{KeyName: fyne.KeyS, Modifier: fyne.KeyModifierControl}
	m.Window.Canvas().AddShortcut(&ctrlS, func(shortcut fyne.Shortcut) { m.saveNote() })
}

func (m *MainWindow) updateFolderContent(folderId uint) error {
	//clear prev content
	m.folderContent.RemoveAll()

	//add return back if it has parent folder
	if folderId != 0 {
		b := widget.NewButtonWithIcon("...", theme.NavigateBackIcon(), func() {
			pid, err := m.api.GetParentFolder(m.currentFolder.ID)
			if err != nil {
				return
			}
			go m.updateFolderContent(pid.ParentFolderID)
		})
		m.folderContent.Add(b)
	}

	//get folder content from api
	curFolder, err := m.api.GetFolderInfo(folderId)
	if err != nil {
		m.folderContent.Add(widget.NewLabel("Not Logged In"))
		return err
	}
	m.currentFolder = &curFolder
	folders, notes, err := m.api.GetFolderContent(folderId)
	if err != nil {
		m.folderContent.Add(widget.NewLabel("Not Logged In"))
		return err
	}

	//add folders
	for _, f := range folders {
		thisFolder := f
		b := widget.NewButtonWithIcon(thisFolder.Name, theme.FolderIcon(), func() {
			m.updateFolderContent(thisFolder.ID)
		})
		m.folderContent.Add(b)
	}

	//add notes
	for _, n := range notes {
		thisNote := n
		b := widget.NewButtonWithIcon(thisNote.Title, theme.FileIcon(), func() {
			m.loadNote(&thisNote)
		})
		m.folderContent.Add(b)
	}

	return nil
}

func (m *MainWindow) loadNote(n *types.Note) {
	m.currentNote = n
	m.noteTitle.SetText(n.Title)
	m.noteContent.SetText(n.Content)
}

func (m *MainWindow) createNote() {
	note, err := m.api.CreateNote("Untitled", "", m.currentFolder.ID, m.currentFolder.GroupId)
	if err != nil {
		log.Fatal(err)
	}
	m.loadNote(&note)
	go m.updateFolderContent(m.currentFolder.ID)
}

func (m *MainWindow) saveNote() {
	if m.currentNote.ID == 0 {
		return
	}
	err := m.api.UpdateNote(m.currentNote.ID, m.noteTitle.Text, m.noteContent.Text, m.currentFolder.ID)
	if err != nil {
		log.Fatal(err)
	}
	go m.updateFolderContent(m.currentFolder.ID)
}

func (m *MainWindow) deleteNote() {
	if m.currentNote.ID == 0 {
		return
	}
	err := m.api.DeleteNote(m.currentNote.ID)
	if err != nil {
		log.Fatal(err)
	}
	m.currentNote = &types.Note{ID: 0}
	m.loadNote(m.currentNote)
	go m.updateFolderContent(m.currentFolder.ID)
}

func (m *MainWindow) createFolder() {
	folderName := widget.NewEntry()
	items := []*widget.FormItem{
		widget.NewFormItem("Folder Name:", folderName),
	}
	dialog.ShowForm("New Folder", "Create", "Cancel", items, func(b bool) {
		if b {
			err := m.api.CreateFolder(folderName.Text, 0)
			if err != nil {
				widget.NewPopUp(container.New(layout.NewCenterLayout(), widget.NewLabel("Error occured during creating folder")), m.Window.Canvas()).Show()
			}
		}
	}, m.Window)
}
