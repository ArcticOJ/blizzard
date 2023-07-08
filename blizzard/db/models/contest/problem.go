package contest

import "blizzard/blizzard/db/models/user"

type Problem struct {
	ID       string   `bun:",pk" json:"id"`
	Tags     []string `bun:",array,notnull" json:"tags"`
	AuthorID string
	Author   user.User `bun:"rel:has-one,join:author_id=id"`
	Content  string    `json:"content"`
}
