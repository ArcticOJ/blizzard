package oauth

import (
	"github.com/ArcticOJ/blizzard/v0/oauth"
	"github.com/ArcticOJ/blizzard/v0/server/http"
)

// GetProviders GET /
func GetProviders(ctx *http.Context) http.Response {
	return ctx.Respond(oauth.EnabledProviders)
}
