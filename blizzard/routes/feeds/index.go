package feeds

import (
	"backend/blizzard/db/models/feeds"
	"backend/blizzard/db/models/shared"
	"backend/blizzard/models"
	"time"
)

func Index(ctx *models.Context) models.Response {
	author := shared.Author{
		Id:       "130139",
		Username: "AlphaNecron",
	}
	return ctx.Arr(feeds.Feed{
		Id:        "welcome",
		Title:     "Welcome to Arctic",
		Timestamp: time.Date(2021, time.September, 12, 9, 50, 24, 0, time.UTC).Format(time.RFC3339),
		Author:    author,
		Content:   "",
	})
}
