package main

import (
	"log"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"github.com/XineAurora/instantnotes-app/internal/ui"
	"github.com/joho/godotenv"
)

func init() {
	if err := godotenv.Load(); err != nil {
		log.Fatal(err)
	}
}

func main() {
	a := app.New()
	w := a.NewWindow("Hello World")

	w.Resize(fyne.NewSize(600, 400))
	var mw ui.MainWindow
	w.SetContent(mw.InitUi())
	w.ShowAndRun()
}
