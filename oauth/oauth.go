package oauth

import (
	"fmt"
	"github.com/ArcticOJ/blizzard/v0/config"
	"github.com/ArcticOJ/blizzard/v0/oauth/providers"
	"golang.org/x/oauth2"
	"slices"
)

var EnabledProviders []string

var providerConf = map[string]providers.ProviderConfig{
	"github":  providers.GitHubProviderConfig,
	"discord": providers.DiscordProviderConfig,
}

var Conf = make(map[string]*oauth2.Config)

var UserInfoHandler = map[string]providers.UserInfoHandler{
	"github":  providers.GetGitHubUser,
	"discord": providers.GetDiscordUser,
}

var AllowedActions = []string{"link", "register", "login"}

func init() {
	for name, provider := range config.Config.Blizzard.OAuth {
		if c, ok := providerConf[name]; ok {
			EnabledProviders = append(EnabledProviders, name)
			Conf[name] = &oauth2.Config{
				// TODO: change this to config-based url
				RedirectURL:  fmt.Sprintf("https://localhost/oauth/%s/callback", name),
				Scopes:       c.Scopes,
				Endpoint:     c.Endpoint,
				ClientID:     provider.ClientID,
				ClientSecret: provider.ClientSecret,
			}
		}
	}
	slices.Sort(EnabledProviders)
}
