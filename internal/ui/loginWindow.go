package ui

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	"github.com/XineAurora/instantnotes-app/internal/api"
)

type LoginWindow struct {
	SignInW          fyne.CanvasObject
	SignUpW          fyne.CanvasObject
	PasswrdRecoveryW fyne.CanvasObject

	LogInChan chan bool

	api *api.ApiConnector
}

func NewLoginWindow(app fyne.App, api *api.ApiConnector) *LoginWindow {
	w := LoginWindow{api: api}
	w.SignInW = w.initSignInUi()
	w.SignInW.Resize(fyne.NewSize(400, 200))
	w.SignUpW = w.initSignUpUi()
	w.SignUpW.Resize(fyne.NewSize(400, 200))
	w.PasswrdRecoveryW = w.initPasswrdRecoveryUi()
	w.PasswrdRecoveryW.Resize(fyne.NewSize(400, 200))

	w.LogInChan = make(chan bool)
	return &w
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
		w.LogInChan <- true

	}
	signUpButton := widget.NewButton("SignUp", func() {})
	passRecButton := widget.NewButton("Forgot Password", func() {})
	return container.NewVBox(form, container.NewHBox(signUpButton, passRecButton), errorLabel)
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
