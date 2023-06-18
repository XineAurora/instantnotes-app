package ui

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
	"github.com/XineAurora/instantnotes-app/internal/api"
	"github.com/XineAurora/instantnotes-app/internal/types"
)

type QuickCreateWindow struct {
	Window     fyne.Window
	saveFolder *types.Folder
	api        *api.ApiConnector
}

func NewQuickCreateWindow(app fyne.App, api *api.ApiConnector) *QuickCreateWindow {
	q := QuickCreateWindow{
		Window:     app.NewWindow("Create Note"),
		saveFolder: &types.Folder{},
		api:        api,
	}
	q.Window.SetContent(q.initUi())
	q.Window.Resize(fyne.NewSize(800, 600))
	return &q
}

func (w *QuickCreateWindow) initUi() fyne.CanvasObject {
	title := widget.NewEntry()
	title.SetPlaceHolder("Type a Title")
	content := widget.NewMultiLineEntry()
	content.SetPlaceHolder("Type a Note here")
	main := container.New(layout.NewBorderLayout(title, nil, nil, nil), title, content)

	saveFoldButton := widget.NewButton("root", func() {})
	saveButton := widget.NewButton("save", func() {})

	saveArea := container.New(layout.NewBorderLayout(nil, nil, nil, saveButton), saveFoldButton, saveButton)

	return container.New(layout.NewBorderLayout(nil, saveArea, nil, nil), main, saveArea)
}
