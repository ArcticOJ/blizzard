package providers

import (
	"fmt"
	"github.com/ArcticOJ/blizzard/v0/logger/debug"
	"github.com/ArcticOJ/blizzard/v0/utils"
	"github.com/ravener/discord-oauth2"
	"net/http"
	"strings"
)

type (
	discordUserInfo struct {
		ID       string `json:"id"`
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
		ID:       info.ID,
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
		avatarUrl = fmt.Sprintf("https://cdn.discordapp.com/avatars/%s/%s.%s", info.ID, info.Avatar, ext)
	} else {
		avatarUrl = fmt.Sprintf("https://cdn.discordapp.com/embed/avatars/%d.png", utils.ParseInt(info.Tag)%5)
	}
	inf.Avatar = avatarUrl
	return inf
}
