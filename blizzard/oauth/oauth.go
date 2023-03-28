package oauth

import (
	"blizzard/blizzard/config"
	"fmt"
	"github.com/ravener/discord-oauth2"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/github"
)

var EnabledProviders []string

type pConf struct {
	scopes   []string
	endpoint oauth2.Endpoint
}

var providerConf = map[string]pConf{
	"github": {
		scopes:   []string{"user:email"},
		endpoint: github.Endpoint,
	},
	"discord": {
		scopes:   []string{discord.ScopeIdentify, discord.ScopeEmail},
		endpoint: discord.Endpoint,
	},
}

var Conf = make(map[string]*oauth2.Config)

func Init() {
	for name, provider := range config.Config.OAuth {
		if c, ok := providerConf[name]; ok {
			EnabledProviders = append(EnabledProviders, name)
			Conf[name] = &oauth2.Config{
				// TODO: change this to config-based url
				RedirectURL:  fmt.Sprintf("https://localhost/api/oauth/%s/callback", name),
				Scopes:       c.scopes,
				Endpoint:     c.endpoint,
				ClientID:     provider.ClientId,
				ClientSecret: provider.ClientSecret,
			}
		}
	}
}
