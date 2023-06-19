package application

import "fyne.io/fyne/v2"

func (a *Application) openSignInWindow() {
	a.window.SetContent(a.loginWidnow.SignInW)
	a.window.Resize(fyne.NewSize(400, 200))
}

func (a *Application) openMainWindow() {
	a.window.SetContent(a.mainWindow.Window)
	a.window.Resize(fyne.NewSize(800, 600))
}

func (a *Application) openGroupWindow() {
	a.window.SetContent(a.groupWindow.Window)
	a.window.Resize(fyne.NewSize(800, 600))
}