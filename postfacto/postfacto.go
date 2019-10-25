package postfacto

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
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

func (c *RetroClient) GetUnfinishedActionItems() ([]ActionItem, error) {
	req, err := c.authorizedRequest("GET", fmt.Sprintf("%s/retros/%s", c.ApiHost, c.ID),
		new(bytes.Buffer))
	if err != nil {
		return nil, err
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("processing action items, HTTP status code: %d", resp.StatusCode)
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	fmt.Println(string(body))
	retro, err := UnmarshalRetro(body)
	if err != nil {
		return nil, err
	}
	actionItems := []ActionItem{}
	for _, item := range retro.Retro.ActionItems {
		if !item.Done {
			actionItems = append(actionItems, item)
		}
	}
	return actionItems, nil
}

func (c *RetroClient) Add(i RetroItem) error {
	b := new(bytes.Buffer)
	_ = json.NewEncoder(b).Encode(i)
	req, err := c.authorizedRequest("POST", fmt.Sprintf("%s/retros/%s/items", c.ApiHost, c.ID), b)
	if err != nil {
		return err
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

func (c *RetroClient) authorizedRequest(method string, url string, b *bytes.Buffer) (*http.Request, error) {
	req, err := http.NewRequest(method, url, b)
	if err != nil {
		return nil, err
	}

	req.Header.Add("Accept", "application/json")
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Origin", c.AppHost)
	req.Header.Add("Referer", fmt.Sprintf("%s/retros/%s", c.AppHost, c.ID))
	req.Header.Add("Sec-Fetch-Mode", "cors")

	if c.Password != "" {
		token, err := c.fetchToken()
		if err != nil {
			return nil, fmt.Errorf("unable to fetch the token: %s", err)
		}
		req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", token))
	}
	return req, err
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
