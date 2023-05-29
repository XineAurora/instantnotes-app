package main

import (
	"log"

	"fyne.io/fyne/v2/app"
	"github.com/XineAurora/instantnotes-app/internal/api"
	"github.com/XineAurora/instantnotes-app/internal/application"
	"github.com/XineAurora/instantnotes-app/internal/ui"
	"github.com/joho/godotenv"
)

func init() {
	if err := godotenv.Load(); err != nil {
		log.Fatal(err)
	}
}

func main() {
	appl := application.Application{App: app.New(), Api: &api.ApiConnector{}}

	mainWindow := ui.NewMainWindow(&appl)
	mainWindow.Window.Show()

	loginWindow := ui.NewLoginWindow(&appl)
	loginWindow.SignInW.Show()

	quickCreate := ui.NewQuickCreateWindow(&appl)
	quickCreate.Window.Show()

	groupWindow := ui.NewGroupWindow(&appl)
	groupWindow.Window.Show()

	// mw := ui.NewMainWindow(&appl)
	// mw.Window.ShowAndRun()
	appl.App.Run()
}
