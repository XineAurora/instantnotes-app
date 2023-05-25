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
		widget.NewToolbarAction(theme.DocumentSaveIcon(), func() {
			err := m.api.UpdateNote(m.currentNote.ID, m.noteTitle.Text, m.noteContent.Text, m.currentFolder.ID)
			if err != nil {
				log.Fatal(err)
			}
		}),
		//delete current note
		widget.NewToolbarAction(theme.ContentRemoveIcon(), func() {}),
	)
	content := container.New(layout.NewBorderLayout(noteToolbar, nil, nil, nil), noteToolbar)
	content.Add(container.New(layout.NewBorderLayout(m.noteTitle, nil, nil, nil), m.noteTitle, m.noteContent))

	//init left side

	m.folderContent = container.NewVBox()
	m.updateFolderContent(0)

	//fill list

	toolBar := widget.NewToolbar(
		//create new note
		widget.NewToolbarAction(theme.ContentAddIcon(), func() {}),
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
