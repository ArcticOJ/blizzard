package routes

import (
	"blizzard/blizzard/models/extra"
	"blizzard/blizzard/routes/auth"
	"blizzard/blizzard/routes/contests"
	"blizzard/blizzard/routes/feeds"
	"blizzard/blizzard/routes/oauth"
	"blizzard/blizzard/routes/problems"
	"blizzard/blizzard/routes/root"
	"blizzard/blizzard/routes/user"
)

var Map = map[string]extra.RouteMap{
	"/problems": problems.Map,
	"/feeds":    feeds.Map,
	"/contests": contests.Map,
	"/auth":     auth.Map,
	"/oauth":    oauth.Map,
	"/user":     user.Map,
	"/":         root.Map,
}
