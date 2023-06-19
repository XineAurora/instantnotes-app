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
	CREATE_FOLDER_ROUTE      = "folder"
	GET_FOLDER_CONTENT_ROUTE = "folder/content/"
	GET_FOLDER_INFO_ROUTE    = "folder/"
	GET_PARENT_FOLDER_ROUTE  = "folder/parent/"

	// Groups operation routes
	GET_ALL_GROUPS_ROUTE    = "group/all"
	GET_GROUP_ROUTE         = "group/"
	GET_GROUP_MEMBERS_ROUTE = "group/members/"
	ADD_GROUP_MEMBER_ROUTE  = "group/member/"
	REMOVE_GROUP_MEMBER_ROUTE  = "group/member/"
)

type ApiConnector struct {
	AuthToken string
	client    http.Client
}
