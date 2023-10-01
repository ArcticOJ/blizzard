package routes

import (
	"blizzard/models/extra"
	"blizzard/routes/auth"
	"blizzard/routes/contests"
	"blizzard/routes/oauth"
	"blizzard/routes/posts"
	"blizzard/routes/problems"
	"blizzard/routes/root"
	"blizzard/routes/submissions"
	"blizzard/routes/user"
	"blizzard/routes/users"
)

var Map = map[string]extra.RouteMap{
	"/problems":    problems.Map,
	"/posts":       posts.Map,
	"/contests":    contests.Map,
	"/auth":        auth.Map,
	"/oauth":       oauth.Map,
	"/users":       users.Map,
	"/user":        user.Map,
	"/submissions": submissions.Map,
	"/":            root.Map,
}
