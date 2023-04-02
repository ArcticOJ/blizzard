package feeds

import (
	"blizzard/blizzard/db/models/feeds"
	"blizzard/blizzard/models"
	"blizzard/blizzard/models/extra"
	"time"
)

func Index(ctx *extra.Context) models.Response {
	return ctx.Arr(feeds.Feed{
		Id:        "welcome",
		Title:     "Welcome to Arctic",
		Timestamp: time.Date(2021, time.September, 12, 9, 50, 24, 0, time.UTC).Format(time.RFC3339),
		Content:   "",
	})
}
