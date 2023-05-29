package ui

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	"github.com/XineAurora/instantnotes-app/internal/api"
	"github.com/XineAurora/instantnotes-app/internal/application"
)

type LoginWindow struct {
	SignInW          fyne.Window
	SignUpW          fyne.Window
	PasswrdRecoveryW fyne.Window
	api              *api.ApiConnector
}

func NewLoginWindow(app *application.Application) LoginWindow {
	w := LoginWindow{api: app.Api}
	w.SignInW = app.App.NewWindow("SignIn")
	w.SignInW.SetContent(w.initSignInUi())
	w.SignInW.Resize(fyne.NewSize(400, 200))
	w.SignUpW = app.App.NewWindow("SignUp")
	w.SignUpW.SetContent(w.initSignUpUi())
	w.SignUpW.Resize(fyne.NewSize(400, 200))
	w.PasswrdRecoveryW = app.App.NewWindow("Password Recovery")
	w.PasswrdRecoveryW.SetContent(w.initPasswrdRecoveryUi())
	w.PasswrdRecoveryW.Resize(fyne.NewSize(400, 200))
	return w
}

func (w *LoginWindow) initSignInUi() fyne.CanvasObject {
	emailEntry := widget.NewEntry()
	passwrdEntry := widget.NewPasswordEntry()
	form := widget.NewForm(
		widget.NewFormItem("email: ", emailEntry),
		widget.NewFormItem("password: ", passwrdEntry),
	)

	errorLabel := widget.NewLabel("")
	form.OnSubmit = func() {
		err := w.api.SignIn(emailEntry.Text, passwrdEntry.Text)
		if err != nil {
			errorLabel.SetText(err.Error())
			passwrdEntry.SetText("")
			return
		}
		// close this window, open main window
		w.SignInW.Hide()
		emailEntry.SetText("")
		passwrdEntry.SetText("")
	}
	return container.NewVBox(form, errorLabel)
}

func (w *LoginWindow) initSignUpUi() fyne.CanvasObject {
	nameEntry := widget.NewEntry()
	emailEntry := widget.NewEntry()
	passwrdEntry := widget.NewPasswordEntry()
	passwrdEntry2 := widget.NewPasswordEntry()
	form := widget.NewForm(
		widget.NewFormItem("name: ", nameEntry),
		widget.NewFormItem("email: ", emailEntry),
		widget.NewFormItem("password: ", passwrdEntry),
		widget.NewFormItem("repeat password: ", passwrdEntry2),
	)

	errorLabel := widget.NewLabel("")
	form.OnSubmit = func() {

	}
	signInButton := widget.NewButton("SignIn", func() {})
	return container.NewVBox(form, errorLabel, signInButton)
}

func (w *LoginWindow) initPasswrdRecoveryUi() fyne.CanvasObject {
	passwrdEntry := widget.NewPasswordEntry()
	passwrdEntry2 := widget.NewPasswordEntry()
	form := widget.NewForm(
		widget.NewFormItem("password: ", passwrdEntry),
		widget.NewFormItem("repeat password: ", passwrdEntry2),
	)

	errorLabel := widget.NewLabel("")
	form.OnSubmit = func() {

	}
	return container.NewVBox(form, errorLabel)
}
