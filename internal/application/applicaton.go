package application

import (
	"fyne.io/fyne/v2"
	"github.com/XineAurora/instantnotes-app/internal/api"
	"github.com/XineAurora/instantnotes-app/internal/ui"
)

type Application struct {
	App fyne.App
	Api *api.ApiConnector

	window fyne.Window

	mainWindow  ui.MainWindow
	loginWidnow ui.LoginWindow
	groupWindow ui.GroupWindow
	quickCreate ui.QuickCreateWindow
}

func New(fyneApp fyne.App, api *api.ApiConnector) *Application {
	app := Application{App: fyneApp, Api: api}

	app.window = fyneApp.NewWindow("Instant Notes")

	app.mainWindow = ui.NewMainWindow(app.window, app.Api)
	app.loginWidnow = ui.NewLoginWindow(app.App, app.Api)
	app.groupWindow = ui.NewGroupWindow(app.App, app.Api)
	app.quickCreate = ui.NewQuickCreateWindow(app.App, app.Api)

	app.window.SetContent(app.loginWidnow.SignInW)
	// open main window after logging in
	go func() {
		for range app.loginWidnow.LogInChan {
			app.openMainWindow()
			go app.mainWindow.UpdateFolderContent(app.mainWindow.CurrentFolder.ID)
		}
	}()

	// open group window after clicking "View"
	go func() {
		for groupId := range app.mainWindow.OpenGroupChan {
			app.openGroupWindow()
			app.groupWindow.LoadGroup(groupId)
		}
	}()

	return &app
}

// this function blocks calling thread until all windows is closed or quit
func (a *Application) Run() {
	//TODO: check saved data if there is jwt token open main window, else open login window
	a.openSignInWindow()
	a.window.Show()
	a.App.Run()
}
