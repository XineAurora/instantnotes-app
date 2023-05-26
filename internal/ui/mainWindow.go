package ui

import (
	"log"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"github.com/XineAurora/instantnotes-app/internal/api"
	"github.com/XineAurora/instantnotes-app/internal/types"
)

type MainWindow struct {
	currentNote   *types.Note
	currentFolder *types.Folder
	folderContent *fyne.Container
	noteTitle     *widget.Entry
	noteContent   *widget.Entry

	api *api.ApiConnector
}

func (m *MainWindow) InitUi() fyne.CanvasObject {
	//init right side
	m.api = &api.ApiConnector{}
	err := m.api.SignIn("test@email.com", "password")
	if err != nil {
		log.Fatal(err)
	}

	m.noteTitle = widget.NewEntry()
	m.noteContent = widget.NewMultiLineEntry()
	noteToolbar := widget.NewToolbar(
		//save current note
		widget.NewToolbarAction(theme.DocumentSaveIcon(), m.saveNote),
		//delete current note
		widget.NewToolbarAction(theme.ContentRemoveIcon(), m.deleteNote),
	)
	content := container.New(layout.NewBorderLayout(noteToolbar, nil, nil, nil), noteToolbar)
	content.Add(container.New(layout.NewBorderLayout(m.noteTitle, nil, nil, nil), m.noteTitle, m.noteContent))

	//init left side

	m.folderContent = container.NewVBox()
	m.updateFolderContent(0)

	//fill list

	toolBar := widget.NewToolbar(
		//create new note
		widget.NewToolbarAction(theme.ContentAddIcon(), m.createNote),
	)

	side := container.New(layout.NewBorderLayout(toolBar, nil, nil, nil), toolBar, m.folderContent)

	split := container.NewHSplit(side, content)
	split.Offset = 0.25
	return split
}

func (m *MainWindow) updateFolderContent(folderId uint) error {
	//get folder content from api
	curFolder, err := m.api.GetFolderInfo(folderId)
	if err != nil {
		return err
	}
	m.currentFolder = &curFolder
	folders, notes, err := m.api.GetFolderContent(folderId)
	if err != nil {
		return err
	}

	//clear prev content
	m.folderContent.RemoveAll()

	//add folders
	for _, f := range folders {
		thisFolder := f
		b := widget.NewButtonWithIcon(thisFolder.Name, theme.FolderIcon(), func() {
			// m.updateFolderContent(thisFolder.ID)
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
	m.updateFolderContent(m.currentFolder.ID)
}

func (m *MainWindow) saveNote() {
	if m.currentNote.ID == 0 {
		return
	}
	err := m.api.UpdateNote(m.currentNote.ID, m.noteTitle.Text, m.noteContent.Text, m.currentFolder.ID)
	if err != nil {
		log.Fatal(err)
	}
	err = m.updateFolderContent(m.currentFolder.ID)
	if err != nil {
		log.Fatal(err)
	}
}

func (m *MainWindow) deleteNote() {
	if m.currentNote.ID == 0 {
		return
	}
	err := m.api.DeleteNote(m.currentNote.ID)
	if err != nil {
		log.Fatal(err)
	}
	m.loadNote(&types.Note{})
	m.updateFolderContent(m.currentFolder.ID)
}
