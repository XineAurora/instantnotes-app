package ui

import (
	"fmt"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
	"github.com/XineAurora/instantnotes-app/internal/api"
	"github.com/XineAurora/instantnotes-app/internal/types"
)

type GroupWindow struct {
	Window fyne.CanvasObject

	api *api.ApiConnector

	groupName    *widget.Label
	groupMembers *fyne.Container
}

func NewGroupWindow(app fyne.App, api *api.ApiConnector) GroupWindow {
	w := GroupWindow{}
	w.Window = w.initUi()
	w.Window.Resize(fyne.NewSize(800, 600))
	w.api = api
	return w
}

func (w *GroupWindow) initUi() fyne.CanvasObject {
	// groupMembers := container.New(layout.NewGridLayout(3),
	// 	widget.NewLabel("Stewie"), widget.NewLabel("test@email.com"), widget.NewLabel("Owner"),
	// 	widget.NewLabel("Brian"), widget.NewLabel("smartiest@email.com"), widget.NewButton("Kick", func() {}),
	// 	widget.NewLabel("Peter"), widget.NewLabel("griffin@email.com"), widget.NewButton("Kick", func() {}),
	// 	widget.NewLabel("Meg"), widget.NewLabel("shut@up.meg"), widget.NewButton("Kick", func() {}),

	w.groupName = widget.NewLabel("")
	w.groupName.Alignment = fyne.TextAlignCenter
	w.groupName.TextStyle.Bold = true

	label := widget.NewLabel("Group Members")
	label.Alignment = fyne.TextAlignCenter

	w.groupMembers = container.New(layout.NewGridLayout(3))

	return container.NewVBox(w.groupName, label, w.groupMembers, widget.NewButton("Add new member", func() {}),
		container.NewHBox(layout.NewSpacer(), widget.NewButton("Delete Group", func() {})))
}

func (w *GroupWindow) LoadGroup(groupId uint) {
	w.groupMembers.RemoveAll()

	group, err := w.api.GetGroup(groupId)
	if err != nil {
		// write an error
		fmt.Println(err)
	}
	w.groupName.SetText(group.Name)
	members, err := w.api.GetGroupMembers(groupId)
	if err != nil {
		// write an error
		fmt.Println(err)
	}
	for _, u := range members {
		if u.ID == group.OwnerID {
			w.addOwner(u)
		}
	}
	for _, u := range members {
		if u.ID != group.OwnerID {
			w.addMember(u)
		}
	}
}

func (w *GroupWindow) addOwner(user types.User) {
	w.groupMembers.Add(widget.NewLabel(user.Name))
	w.groupMembers.Add(widget.NewLabel(user.Email))
	w.groupMembers.Add(widget.NewLabel("Owner"))
}

func (w *GroupWindow) addMember(user types.User) {
	w.groupMembers.Add(widget.NewLabel(user.Name))
	w.groupMembers.Add(widget.NewLabel(user.Email))
	//TODO: check if user is owner to add this button
	w.groupMembers.Add(widget.NewButton("Kick", func() {}))
}
