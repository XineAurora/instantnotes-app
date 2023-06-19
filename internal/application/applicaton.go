package application

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/driver/desktop"
	"fyne.io/fyne/v2/theme"
	"github.com/XineAurora/instantnotes-app/internal/api"
	"github.com/XineAurora/instantnotes-app/internal/ui"
	hook "github.com/robotn/gohook"
)

type Application struct {
	App fyne.App
	Api *api.ApiConnector

	window fyne.Window

	mainWindow  *ui.MainWindow
	loginWidnow *ui.LoginWindow
	groupWindow *ui.GroupWindow
	quickCreate *ui.QuickCreateWindow
}

func New(fyneApp fyne.App, api *api.ApiConnector) *Application {
	fyneApp.SetIcon(theme.GridIcon())
	app := Application{App: fyneApp, Api: api}

	app.window = fyneApp.NewWindow("Instant Notes")
	app.window.SetCloseIntercept(func() {
		app.window.Hide()
	})

	app.mainWindow = ui.NewMainWindow(app.window, app.Api)
	app.loginWidnow = ui.NewLoginWindow(app.App, app.Api)
	app.groupWindow = ui.NewGroupWindow(app.window, app.Api)
	app.quickCreate = ui.NewQuickCreateWindow(app.App, app.Api)

	if desk, ok := fyneApp.(desktop.App); ok {
		m := fyne.NewMenu("Instant Notes",
			fyne.NewMenuItem("Show", func() {
				app.window.Show()
			}),
			fyne.NewMenuItem("Create Note", func() {
				app.quickCreate.Window.Show()
			}))
		desk.SetSystemTrayMenu(m)
	}

	app.window.SetContent(app.loginWidnow.SignInW)

	//start hook
	go app.hook()

	// open main window...
	go func() {
		for {
			select {
			// open main window after logging in
			case <-app.loginWidnow.LogInChan:
				app.openMainWindow()
				go app.mainWindow.UpdateFolderContent(app.mainWindow.CurrentFolder.ID)
				// open main window after navigating back
			case <-app.groupWindow.LoadMainChan:
				app.openMainWindow()
				go app.mainWindow.UpdateFolderContent(app.mainWindow.CurrentFolder.ID)
			}

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

func (a *Application) hook() {
	hook.Register(hook.KeyDown, []string{"k", "shift", "ctrl"}, func(e hook.Event) {
		if a.Api.AuthToken != "" {
			a.quickCreate.Window.Show()
			a.quickCreate.Window.RequestFocus()
			a.quickCreate.Window.CenterOnScreen()
		}
	})
	s := hook.Start()
	defer hook.End()
	<-hook.Process(s)
}
