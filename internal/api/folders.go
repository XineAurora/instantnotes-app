package api

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"

	"github.com/XineAurora/instantnotes-app/internal/types"
)

func (a *ApiConnector) GetFolderContent(id uint) ([]types.Folder, []types.Note, error) {
	// assemble url string
	url := fmt.Sprintf("%s%s%v", os.Getenv("API_URL"), GET_FOLDER_CONTENT_ROUTE, id)
	// request data

	//create request
	req, _ := http.NewRequest("GET", url, nil)

	//set auth info in header
	a.SetAuthInfo(&req.Header)

	resp, err := a.client.Do(req)
	if err != nil {
		log.Fatal("Failed to get folder content", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		log.Printf("error occured during request %v, server returned %v", url, resp.Status)
		return nil, nil, errors.New("error " + resp.Status)
	}

	// read body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}
	// demarshal data from body
	var notes []types.Note
	var folders []types.Folder

	err = json.Unmarshal(body, &struct {
		Folders *[]types.Folder `json:"folders"`
		Notes   *[]types.Note   `json:"notes"`
	}{
		&folders,
		&notes,
	})

	if err != nil {
		log.Fatal(err)
	}

	return folders, notes, nil
}

func (a *ApiConnector) GetFolderInfo(id uint) (types.Folder, error) {
	if id == 0 {
		return types.Folder{ID: 0, Name: ""}, nil
	}
	url := fmt.Sprintf("%s%s%v", os.Getenv("API_URL"), GET_FOLDER_INFO_ROUTE, id)
	req, _ := http.NewRequest("GET", url, nil)
	a.SetAuthInfo(&req.Header)
	resp, err := a.client.Do(req)
	if err != nil {
		log.Fatal(err)
		return types.Folder{}, err
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	var folder types.Folder

	err = json.Unmarshal(body, &folder)
	if err != nil {
		log.Fatal(err)
		return types.Folder{}, err
	}

	return folder, nil
}
