package ui

import (
	"fmt"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"github.com/XineAurora/instantnotes-app/internal/api"
	"github.com/XineAurora/instantnotes-app/internal/types"
)

type GroupWindow struct {
	Window fyne.CanvasObject

	api    *api.ApiConnector
	parent fyne.Window

	LoadMainChan chan bool

	groupName    *widget.Label
	groupMembers *fyne.Container
	currentGroup *types.Group
}

func NewGroupWindow(parent fyne.Window, api *api.ApiConnector) *GroupWindow {
	w := GroupWindow{}
	w.Window = w.initUi()
	w.Window.Resize(fyne.NewSize(800, 600))
	w.api = api
	w.parent = parent
	w.LoadMainChan = make(chan bool)
	w.currentGroup = &types.Group{}
	return &w
}

func (w *GroupWindow) initUi() fyne.CanvasObject {
	// groupMembers := container.New(layout.NewGridLayout(3),
	// 	widget.NewLabel("Stewie"), widget.NewLabel("test@email.com"), widget.NewLabel("Owner"),
	// 	widget.NewLabel("Brian"), widget.NewLabel("smartiest@email.com"), widget.NewButton("Kick", func() {}),
	// 	widget.NewLabel("Peter"), widget.NewLabel("griffin@email.com"), widget.NewButton("Kick", func() {}),
	// 	widget.NewLabel("Meg"), widget.NewLabel("shut@up.meg"), widget.NewButton("Kick", func() {}),

	backButton := widget.NewButtonWithIcon("Back", theme.NavigateBackIcon(), func() { w.LoadMainChan <- true })

	w.groupName = widget.NewLabel("")
	w.groupName.Alignment = fyne.TextAlignCenter
	w.groupName.TextStyle.Bold = true

	label := widget.NewLabel("Group Members")
	label.Alignment = fyne.TextAlignCenter

	w.groupMembers = container.New(layout.NewGridLayout(3))

	addNewMemberButton := widget.NewButton("Add new member", w.addNewMember)
	return container.NewVBox(backButton, w.groupName, label, w.groupMembers, addNewMemberButton,
		container.NewHBox(layout.NewSpacer(), widget.NewButton("Delete Group", func() {})))
}

func (w *GroupWindow) LoadGroup(groupId uint) {
	w.groupMembers.RemoveAll()

	group, err := w.api.GetGroup(groupId)
	w.currentGroup = &group

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
			w.addOwnerToList(u)
		}
	}
	for _, u := range members {
		if u.ID != group.OwnerID {
			w.addMemberToList(u)
		}
	}
}

func (w *GroupWindow) addOwnerToList(user types.User) {
	w.groupMembers.Add(widget.NewLabel(user.Name))
	w.groupMembers.Add(widget.NewLabel(user.Email))
	w.groupMembers.Add(widget.NewLabel("Owner"))
}

func (w *GroupWindow) addMemberToList(user types.User) {
	w.groupMembers.Add(widget.NewLabel(user.Name))
	w.groupMembers.Add(widget.NewLabel(user.Email))
	//TODO: check if user is owner to add this button
	w.groupMembers.Add(widget.NewButton("Kick", func() { w.removeMember(user.ID) }))
}

func (w *GroupWindow) addNewMember() {
	emailEntry := widget.NewEntry()
	dialog.ShowForm("Add New Member", "Add", "Cancel", []*widget.FormItem{widget.NewFormItem("Email", emailEntry)}, func(b bool) {
		if b {
			if err := w.api.AddGroupMember(emailEntry.Text, w.currentGroup.ID); err != nil {
				widget.NewPopUp(container.New(layout.NewCenterLayout(), widget.NewLabel("There is no user with this email. "+err.Error())), w.parent.Canvas()).Show()
			}
		}
		w.LoadGroup(w.currentGroup.ID)
	}, w.parent)
}

func (w *GroupWindow) removeMember(userId uint) {
	if err := w.api.RemoveGroupMember(userId, w.currentGroup.ID); err != nil {
		fmt.Println(err)
	}
	w.LoadGroup(w.currentGroup.ID)
}
