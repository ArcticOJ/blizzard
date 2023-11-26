package routes

import (
	"github.com/ArcticOJ/blizzard/v0/routes/apex"
	"github.com/ArcticOJ/blizzard/v0/routes/auth"
	"github.com/ArcticOJ/blizzard/v0/routes/contests"
	"github.com/ArcticOJ/blizzard/v0/routes/oauth"
	"github.com/ArcticOJ/blizzard/v0/routes/posts"
	"github.com/ArcticOJ/blizzard/v0/routes/problems"
	"github.com/ArcticOJ/blizzard/v0/routes/submissions"
	"github.com/ArcticOJ/blizzard/v0/routes/user"
	"github.com/ArcticOJ/blizzard/v0/routes/users"
	"github.com/ArcticOJ/blizzard/v0/server/http"
)

var Map = map[string]http.RouteMap{
	"/problems":    problems.Map,
	"/posts":       posts.Map,
	"/contests":    contests.Map,
	"/auth":        auth.Map,
	"/oauth":       oauth.Map,
	"/users":       users.Map,
	"/user":        user.Map,
	"/submissions": submissions.Map,
	"/":            apex.Map,
}
