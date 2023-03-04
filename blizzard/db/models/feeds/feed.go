package feeds

import (
	"backend/blizzard/db/models/shared"
	"github.com/uptrace/bun"
)

type Feed struct {
	bun.BaseModel `bun:"table:feeds"`
	Id            string        `bun:"type:uuid, unique, default:gen_random_uuid()" json:"id"`
	Title         string        `json:"title"`
	Timestamp     string        `json:"timestamp"`
	Author        shared.Author `json:"author"`
	Content       string        `json:"content"`
}
