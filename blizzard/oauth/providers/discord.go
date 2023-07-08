package providers

import (
	"blizzard/blizzard/logger/debug"
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
	Scopes:   []string{discord.ScopeIdentify},
	Endpoint: discord.Endpoint,
}

func GetDiscordUser(client *http.Client) *UserInfo {
	info := readBody[discordUserInfo](client.Get("https://discord.com/api/v10/users/@me"))
	var avatarUrl string
	debug.Dump(info)
	if info == nil {
		return nil
	}
	inf := &UserInfo{
		Id:       info.Id,
		Username: info.Username,
	}
	// user has migrated to the new username system
	if len(info.Tag) == 4 {
		inf.Username = fmt.Sprintf("%s#%s", info.Username, info.Tag)
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
	inf.Avatar = avatarUrl
	return inf
}
