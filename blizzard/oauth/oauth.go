package oauth

import (
	"blizzard/blizzard/config"
	"blizzard/blizzard/oauth/providers"
	"fmt"
	"golang.org/x/oauth2"
	"sort"
)

var EnabledProviders []string

// TODO: move to models

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

func Init() {
	for name, provider := range config.Config.OAuth {
		if c, ok := providerConf[name]; ok {
			EnabledProviders = append(EnabledProviders, name)
			Conf[name] = &oauth2.Config{
				// TODO: change this to config-based url
				RedirectURL:  fmt.Sprintf("https://localhost/oauth/%s/callback", name),
				Scopes:       c.Scopes,
				Endpoint:     c.Endpoint,
				ClientID:     provider.ClientId,
				ClientSecret: provider.ClientSecret,
			}
		}
	}
	sort.Slice(EnabledProviders, func(i, j int) bool {
		return EnabledProviders[i] < EnabledProviders[j]
	})
}
