package providers

import (
	"blizzard/blizzard/utils"
	"fmt"
	"github.com/ravener/discord-oauth2"
	"net/http"
	"strings"
)

type (
	discordUserInfo struct {
		Id       string `json:"id"`
		Username string `json:"username"`
		Tag      string `json:"discriminator"`
		Avatar   string `json:"avatar"`
	}
)

var DiscordProviderConfig = ProviderConfig{
	Scopes:   []string{discord.ScopeIdentify, discord.ScopeEmail},
	Endpoint: discord.Endpoint,
}

func GetDiscordUser(client *http.Client) *UserInfo {
	info := readBody[discordUserInfo](client.Get("https://discord.com/api/v10/users/@me"))
	var avatarUrl string
	if info == nil {
		return nil
	}
	if len(info.Avatar) > 0 {
		ext := "png"
		if strings.HasPrefix(info.Avatar, "a_") {
			ext = "gif"
		}
		avatarUrl = fmt.Sprintf("https://cdn.discordapp.com/avatars/%s/%s.%s", info.Id, info.Avatar, ext)
	} else {
		avatarUrl = fmt.Sprintf("https://cdn.discordapp.com/embed/avatars/%d.png", utils.ParseInt(info.Tag)%5)
	}
	return &UserInfo{
		Id:       info.Id,
		Username: fmt.Sprintf("%s#%s", info.Username, info.Tag),
		Avatar:   avatarUrl,
	}
}
