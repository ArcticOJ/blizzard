package posts

import (
	"github.com/ArcticOJ/blizzard/v0/db"
	"github.com/ArcticOJ/blizzard/v0/db/schema/post"
	"github.com/ArcticOJ/blizzard/v0/server/http"
)

// GetPosts GET /
func GetPosts(ctx *http.Context) http.Response {
	var posts []post.Post
	if db.Database.NewSelect().Model(&posts).Scan(ctx.Request().Context()) != nil {
		return ctx.InternalServerError("Could not fetch posts.")
	}
	return ctx.Respond(posts)
}
