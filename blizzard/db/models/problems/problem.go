package problems

import (
	"backend/blizzard/db/models/shared"
)

type Problem struct {
	Id      string        `json:"id"`
	Title   string        `json:"title"`
	Contest string        `json:"contest"`
	Tags    []string      `bun:"tags" json:"tags"`
	Author  shared.Author `json:"author"`
	Content string        `json:"content"`
}
