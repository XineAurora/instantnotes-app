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
)

func (a *ApiConnector) SignIn(email string, password string) error {
	url := fmt.Sprintf("%s%s", os.Getenv("API_URL"), SIGN_IN_ROUTE)
	// mashalize data
	jsonData, err := json.Marshal(struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}{
		email,
		password,
	})
	if err != nil {
		log.Fatal(err)
		return err
	}

	resp, err := http.Post(url, "application/json", bytes.NewReader(jsonData))
	if err != nil {
		//log.Fatal(err)
		return err
	}

	if resp.StatusCode == http.StatusOK {
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			//log.Fatal(err)
			return err
		}
		// demarshal data from body
		var data map[string]string

		err = json.Unmarshal(body, &data)
		if err != nil {
			//log.Fatal(err)
			return err
		}
		a.AuthToken = data["token"]
		return nil
	}
	return errors.New("wrong Email or Password")
}

func (a *ApiConnector) SetAuthInfo(h *http.Header) {
	h.Set("Authorization", a.AuthToken)
}
