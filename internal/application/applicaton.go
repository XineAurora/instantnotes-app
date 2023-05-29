package application

import (
	"fyne.io/fyne/v2"
	"github.com/XineAurora/instantnotes-app/internal/api"
)

type Application struct {
	App fyne.App
	Api *api.ApiConnector
}
