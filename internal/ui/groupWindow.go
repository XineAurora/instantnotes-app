package ui

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
	"github.com/XineAurora/instantnotes-app/internal/application"
)

type GroupWindow struct {
	Window fyne.Window
}

func NewGroupWindow(app *application.Application) GroupWindow {
	w := GroupWindow{}
	w.Window = app.App.NewWindow("Group")
	w.Window.SetContent(w.initUi())
	w.Window.Resize(fyne.NewSize(800, 600))
	return w
}

func (w *GroupWindow) initUi() fyne.CanvasObject {
	groupName := widget.NewLabel("Test Group 1")
	groupName.Alignment = fyne.TextAlignCenter
	groupName.TextStyle.Bold = true

	groupMembers := container.New(layout.NewGridLayout(3),
		widget.NewLabel("Stewie"), widget.NewLabel("test@email.com"), widget.NewLabel("Owner"),
		widget.NewLabel("Brian"), widget.NewLabel("smartiest@email.com"), widget.NewButton("Kick", func() {}),
		widget.NewLabel("Peter"), widget.NewLabel("griffin@email.com"), widget.NewButton("Kick", func() {}),
		widget.NewLabel("Meg"), widget.NewLabel("shut@up.meg"), widget.NewButton("Kick", func() {}),
	)
	label := widget.NewLabel("Group Members")
	label.Alignment = fyne.TextAlignCenter
	return container.NewVBox(groupName, label, groupMembers, widget.NewButton("Add new member", func() {}),
		container.NewHBox(layout.NewSpacer(), widget.NewButton("Delete Group", func() {})))
}
