package providers

import (
	"encoding/json"
	"golang.org/x/oauth2"
	"net/http"
)

type (
	ProviderConfig struct {
		Scopes   []string
		Endpoint oauth2.Endpoint
	}

	UserInfo struct {
		Username string `json:"username"`
		Id       string `json:"id"`
		Avatar   string `json:"avatar"`
	}

	UserInfoHandler = func(client *http.Client) *UserInfo
)

func readBody[T any](response *http.Response, err error) *T {
	if err != nil || response.StatusCode != http.StatusOK {
		return nil
	}
	d := json.NewDecoder(response.Body)
	var body T
	if !d.More() || d.Decode(&body) != nil {
		return nil
	}
	return &body
}
