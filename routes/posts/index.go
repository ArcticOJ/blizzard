package posts

import (
	"blizzard/db"
	"blizzard/db/models/post"
	"blizzard/server/http"
)

func Index(ctx *http.Context) http.Response {
	var posts []post.Post
	if db.Database.NewSelect().Model(&posts).Scan(ctx.Request().Context()) != nil {
		return ctx.InternalServerError("Could not fetch posts.")
	}
	return ctx.Respond(posts)
}
