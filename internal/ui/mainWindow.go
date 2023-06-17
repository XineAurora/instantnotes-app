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
	"github.com/XineAurora/instantnotes-app/internal/types"
)

type MainWindow struct {
	parent fyne.Window
	Window fyne.CanvasObject

	currentNote   *types.Note
	CurrentFolder *types.Folder
	folderContent *fyne.Container
	noteTitle     *widget.Entry
	noteContent   *widget.Entry

	OpenGroupChan chan uint

	api *api.ApiConnector
}

func NewMainWindow(p fyne.Window, api *api.ApiConnector) MainWindow {
	mw := MainWindow{}
	mw.api = api
	mw.Window = mw.initUi()
	mw.Window.Resize(fyne.NewSize(800, 600))
	mw.parent = p
	//setup shortcuts
	mw.setupShortcuts()
	mw.currentNote = &types.Note{ID: 0}
	mw.CurrentFolder = &types.Folder{ID: 0}

	mw.OpenGroupChan = make(chan uint)
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

		//groups button
		widget.NewToolbarAction(theme.AccountIcon(), func() { m.newGroupsDialog().Show() }),
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
		widget.NewToolbarAction(theme.ViewRefreshIcon(), func() { go m.UpdateFolderContent(m.CurrentFolder.ID) }),
	)

	side := container.New(layout.NewBorderLayout(toolBar, nil, nil, nil), toolBar, scrollArea)

	split := container.NewHSplit(side, content)
	split.Offset = 0.25
	return split
}

func (m *MainWindow) newGroupsDialog() dialog.Dialog {

	//get all groups
	groups, err := m.api.GetAllGroups()
	if err != nil {
		return dialog.NewInformation("Error", "Cannot get groups: "+err.Error(), m.parent)
	}

	// create table with name, view group button
	grid := container.New(layout.NewGridLayout(3))
	var dial dialog.Dialog
	for _, group := range groups {
		thisGroup := group
		grid.Add(widget.NewLabel(group.Name))
		grid.Add(widget.NewButton("View", func() {
			// open groupWindow with info about this group
			m.OpenGroupChan <- thisGroup.ID
			dial.Hide()
		}))
	}
	grid.Resize(fyne.NewSize(100, 200))
	dial = dialog.NewCustom("Groups", "Back", grid, m.parent)
	return dial
}

func (m *MainWindow) setupShortcuts() {
	ctrlS := desktop.CustomShortcut{KeyName: fyne.KeyS, Modifier: fyne.KeyModifierControl}
	m.parent.Canvas().AddShortcut(&ctrlS, func(shortcut fyne.Shortcut) { m.saveNote() })
}

func (m *MainWindow) UpdateFolderContent(folderId uint) error {
	//clear prev content
	m.folderContent.RemoveAll()

	//add return back if it has parent folder
	if folderId != 0 {
		b := widget.NewButtonWithIcon("...", theme.NavigateBackIcon(), func() {
			pid, err := m.api.GetParentFolder(m.CurrentFolder.ID)
			if err != nil {
				return
			}
			go m.UpdateFolderContent(pid.ParentFolderID)
		})
		m.folderContent.Add(b)
	}

	//get folder content from api
	curFolder, err := m.api.GetFolderInfo(folderId)
	if err != nil {
		m.folderContent.Add(widget.NewLabel("Not Logged In"))
		return err
	}
	m.CurrentFolder = &curFolder
	folders, notes, err := m.api.GetFolderContent(folderId)
	if err != nil {
		m.folderContent.Add(widget.NewLabel("Not Logged In"))
		return err
	}

	//add folders
	for _, f := range folders {
		if f.GroupId == 0 || m.CurrentFolder.GroupId != 0 {
			thisFolder := f
			b := widget.NewButtonWithIcon(thisFolder.Name, theme.FolderIcon(), func() {
				m.UpdateFolderContent(thisFolder.ID)
			})
			m.folderContent.Add(b)
		}
	}

	//add notes
	for _, n := range notes {
		thisNote := n
		b := widget.NewButtonWithIcon(thisNote.Title, theme.FileIcon(), func() {
			m.loadNote(&thisNote)
		})
		m.folderContent.Add(b)
	}

	//add group folder
	if m.CurrentFolder.GroupId == 0 && m.CurrentFolder.ID == 0 {
		m.folderContent.Add(widget.NewLabel("Groups"))
		for _, f := range folders {
			if f.GroupId != 0 {
				thisFolder := f
				b := widget.NewButtonWithIcon(thisFolder.Name, theme.FolderIcon(), func() {
					m.UpdateFolderContent(thisFolder.ID)
				})
				m.folderContent.Add(b)
			}
		}
	}
	return nil
}

func (m *MainWindow) loadNote(n *types.Note) {
	m.currentNote = n
	m.noteTitle.SetText(n.Title)
	m.noteContent.SetText(n.Content)
}

func (m *MainWindow) createNote() {
	note, err := m.api.CreateNote("Untitled", "", m.CurrentFolder.ID, m.CurrentFolder.GroupId)
	if err != nil {
		log.Fatal(err)
	}
	m.loadNote(&note)
	go m.UpdateFolderContent(m.CurrentFolder.ID)
}

func (m *MainWindow) saveNote() {
	if m.currentNote.ID == 0 {
		return
	}
	err := m.api.UpdateNote(m.currentNote.ID, m.noteTitle.Text, m.noteContent.Text, m.CurrentFolder.ID)
	if err != nil {
		log.Fatal(err)
	}
	go m.UpdateFolderContent(m.CurrentFolder.ID)
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
	go m.UpdateFolderContent(m.CurrentFolder.ID)
}

func (m *MainWindow) createFolder() {
	folderName := widget.NewEntry()
	items := []*widget.FormItem{
		widget.NewFormItem("Folder Name:", folderName),
	}
	dialog.ShowForm("New Folder", "Create", "Cancel", items, func(b bool) {
		if b {
			err := m.api.CreateFolder(folderName.Text, m.CurrentFolder.GroupId, m.CurrentFolder.ID)
			if err != nil {
				widget.NewPopUp(container.New(layout.NewCenterLayout(), widget.NewLabel("Error occured during creating folder")), m.parent.Canvas()).Show()
			}
			go m.UpdateFolderContent(m.CurrentFolder.ID)
		}
	}, m.parent)
}
