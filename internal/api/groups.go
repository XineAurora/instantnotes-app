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

func (a *ApiConnector) GetAllGroups() ([]types.Group, error) {
	url := fmt.Sprintf("%s%s", os.Getenv("API_URL"), GET_ALL_GROUPS_ROUTE)
	req, _ := http.NewRequest("GET", url, nil)
	a.SetAuthInfo(&req.Header)

	resp, err := a.client.Do(req)
	if err != nil {
		log.Fatal(err)
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		log.Printf("error occured during request %v, server returned %v", url, resp.Status)
		return nil, errors.New("error " + resp.Status)
	}

	// read body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}
	// demarshal data from body
	var groups []types.Group

	err = json.Unmarshal(body, &groups)

	if err != nil {
		log.Fatal(err)
	}

	return groups, nil
}
