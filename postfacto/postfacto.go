package postfacto

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httputil"
)

type RetroClient struct {
	ApiHost  string
	AppHost  string
	ID       string
	Password string
}

type Category string

const (
	CategoryHappy Category = "happy"
	CategoryMeh   Category = "meh"
	CategorySad   Category = "sad"
)

type RetroItem struct {
	Description string   `json:"description"`
	Category    Category `json:"category"`
}

type passwordPayload struct {
	Retro retroPassword `json:"retro"`
}

type retroPassword struct {
	Password string `json:"password"`
}

type tokenReply struct {
	Token string `json:"token"`
}

func (c *RetroClient) Add(i RetroItem) error {
	b := new(bytes.Buffer)
	_ = json.NewEncoder(b).Encode(i)
	req, err := http.NewRequest("POST", fmt.Sprintf("%s/retros/%s/items", c.ApiHost, c.ID), b)
	if err != nil {
		return err
	}

	req.Header.Add("Accept", "application/json")
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Origin", c.AppHost)
	req.Header.Add("Referer", fmt.Sprintf("%s/retros/%s", c.AppHost, c.ID))
	req.Header.Add("Sec-Fetch-Mode", "cors")

	if c.Password != "" {
		token, err := c.fetchToken()
		if err != nil {
			return fmt.Errorf("unable to fetch the token: %s", err)
		}
		req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", token))
	}

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}

	if res.StatusCode != http.StatusCreated {
		b, _ := httputil.DumpResponse(res, true)
		return fmt.Errorf("unexpected response code (%d) - %s", res.StatusCode, string(b))
	}

	return nil
}

func (c *RetroClient) fetchToken() (string, error) {
	payload := passwordPayload{Retro: retroPassword{Password: c.Password}}
	b := new(bytes.Buffer)
	_ = json.NewEncoder(b).Encode(payload)
	loginURL := fmt.Sprintf("%s/retros/%s/login", c.ApiHost, c.ID)
	req, err := http.NewRequest("PUT", loginURL, b)
	if err != nil {
		return "", fmt.Errorf("unable to create a request - %s", err)
	}

	req.Header.Add("Accept", "application/json")
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Origin", c.AppHost)
	req.Header.Add("Referer", fmt.Sprintf("%s/retros/%s/login", c.AppHost, c.ID))
	req.Header.Add("Sec-Fetch-Mode", "cors")

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("unable to perform login request: %s", err)
	}

	defer res.Body.Close()

	reply := tokenReply{}
	err = json.NewDecoder(res.Body).Decode(&reply)
	if err != nil {
		return "", fmt.Errorf("unable to decode the reply: %s\nStatus code: %d", err, res.StatusCode)
	}

	return reply.Token, err
}
