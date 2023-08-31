package routes

import (
	"blizzard/blizzard/models/extra"
	"blizzard/blizzard/routes/auth"
	"blizzard/blizzard/routes/contests"
	"blizzard/blizzard/routes/oauth"
	"blizzard/blizzard/routes/posts"
	"blizzard/blizzard/routes/problems"
	"blizzard/blizzard/routes/root"
	"blizzard/blizzard/routes/submissions"
	"blizzard/blizzard/routes/user"
	"blizzard/blizzard/routes/users"
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
