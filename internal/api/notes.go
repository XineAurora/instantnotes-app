package api

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"

	"github.com/XineAurora/instantnotes-app/internal/types"
)

func (a *ApiConnector) GetNote(id uint) (types.Note, error) {
	url := fmt.Sprintf("%s%s%v", os.Getenv("API_URL"), GET_NOTE_ROUTE, id)

	req, _ := http.NewRequest("GET", url, nil)
	a.SetAuthInfo(&req.Header)
	resp, err := a.client.Do(req)
	if err != nil {
		log.Fatal(err)
		return types.Note{}, err
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	var note types.Note

	err = json.Unmarshal(body, &note)
	if err != nil {
		log.Fatal(err)
		return types.Note{}, err
	}

	return note, nil
}

func (a *ApiConnector) UpdateNote(id uint, title string, content string, folderId uint) error {
	url := fmt.Sprintf("%s%s%v", os.Getenv("API_URL"), UPDATE_NOTE_ROUTE, id)

	body, _ := json.Marshal(struct {
		Title    string
		Content  string
		FolderID uint
	}{
		Title:    title,
		Content:  content,
		FolderID: folderId,
	})
	req, _ := http.NewRequest("PUT", url, bytes.NewBuffer(body))
	a.SetAuthInfo(&req.Header)
	req.Header.Set("Content-Type", "application/json")

	resp, err := a.client.Do(req)
	if err != nil {
		log.Fatal(err)
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		log.Printf("error occured during request %v, server returned %v", url, resp.Status)
		return errors.New("error " + resp.Status)
	}

	return nil
}
