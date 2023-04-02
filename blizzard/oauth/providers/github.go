package providers

import (
	"encoding/json"
	"golang.org/x/oauth2/github"
	"net/http"
)

type (
	githubUserInfo struct {
		Username  string      `json:"login"`
		Id        json.Number `json:"id"`
		AvatarUrl string      `json:"avatar_url"`
	}
)

var GitHubProviderConfig = ProviderConfig{
	Scopes:   []string{"user:email", "read:user"},
	Endpoint: github.Endpoint,
}

func GetGitHubUser(client *http.Client) *UserInfo {
	info := readBody[githubUserInfo](client.Get("https://api.github.com/user"))
	if info == nil {
		return nil
	}
	return &UserInfo{
		Id:       info.Id.String(),
		Username: info.Username,
		Avatar:   info.AvatarUrl,
	}
}
