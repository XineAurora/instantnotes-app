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

func (a *ApiConnector) GetGroup(id uint) (types.Group, error) {
	url := fmt.Sprintf("%s%s%v", os.Getenv("API_URL"), GET_GROUP_ROUTE, id)

	req, _ := http.NewRequest("GET", url, nil)
	a.SetAuthInfo(&req.Header)
	resp, err := a.client.Do(req)
	if err != nil {
		log.Fatal(err)
		return types.Group{}, err
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	var group types.Group

	err = json.Unmarshal(body, &group)
	if err != nil {
		log.Fatal(err)
		return types.Group{}, err
	}

	return group, nil
}

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

func (a *ApiConnector) GetGroupMembers(id uint) ([]types.User, error) {
	url := fmt.Sprintf("%s%s%v", os.Getenv("API_URL"), GET_GROUP_MEMBERS_ROUTE, id)

	req, _ := http.NewRequest("GET", url, nil)
	a.SetAuthInfo(&req.Header)
	resp, err := a.client.Do(req)
	if err != nil {
		log.Fatal(err)
		return []types.User{}, err
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	var members []types.User

	err = json.Unmarshal(body, &struct {
		Members *[]types.User `json:"members"`
	}{
		&members,
	})

	if err != nil {
		log.Fatal(err)
		return []types.User{}, err
	}

	return members, nil
}

func (a *ApiConnector) AddGroupMember(email string, groupId uint) error {
	url := fmt.Sprintf("%s%s%v", os.Getenv("API_URL"), ADD_GROUP_MEMBER_ROUTE, groupId)

	body, _ := json.Marshal(struct {
		Email string
	}{
		Email: email,
	})
	req, _ := http.NewRequest("POST", url, bytes.NewBuffer(body))
	a.SetAuthInfo(&req.Header)
	req.Header.Set("Content-Type", "application/json")

	resp, err := a.client.Do(req)
	if err != nil {
		log.Fatal(err)
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return errors.New(resp.Status)
	}
	return nil
}

func (a *ApiConnector) RemoveGroupMember(userId uint, groupId uint) error {
	url := fmt.Sprintf("%s%s%v", os.Getenv("API_URL"), REMOVE_GROUP_MEMBER_ROUTE, groupId)
	body, _ := json.Marshal(struct {
		UserID uint `json:"UserID"`
	}{
		UserID: userId,
	})

	req, _ := http.NewRequest("DELETE", url, bytes.NewBuffer(body))
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
