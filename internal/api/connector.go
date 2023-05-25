package api

import (
	"net/http"
)

const (
	// Auth operation routes
	SIGN_UP_ROUTE = "auth/signup"
	SIGN_IN_ROUTE = "auth/signin"

	// Notes operation routes
	CREATE_NOTE_ROUTE = "api/notes"
	GET_NOTE_ROUTE    = "api/notes/"
	UPDATE_NOTE_ROUTE = "api/notes/"
	DELETE_NOTE_ROUTE = "api/notes/"

	// Folder operation routes
	GET_FOLDER_CONTENT_ROUTE = "folder/content/"
	GET_FOLDER_INFO_ROUTE    = "folder/"
)

type ApiConnector struct {
	AuthToken string
	client    http.Client
}
