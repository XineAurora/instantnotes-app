package main

import (
	"log"

	"fyne.io/fyne/v2/app"
	"github.com/XineAurora/instantnotes-app/internal/api"
	"github.com/XineAurora/instantnotes-app/internal/application"
	"github.com/joho/godotenv"
)

func init() {
	if err := godotenv.Load(); err != nil {
		log.Fatal(err)
	}
}

func main() {
	appl := application.New(app.New(), &api.ApiConnector{})
	appl.Run()
	
}
