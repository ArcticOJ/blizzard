package post

import (
	"github.com/ArcticOJ/blizzard/v0/db/models/user"
	"github.com/google/uuid"
	"time"
)

type Post struct {
	ID       uint32     `bun:",pk,autoincrement" json:"id"`
	Title    string     `bun:",notnull" json:"title"`
	PostedAt *time.Time `bun:",nullzero,notnull,default:'now()'" json:"postedAt,omitempty"`
	AuthorID uuid.UUID  `bun:",type:uuid" json:"authorID,omitempty"`
	Author   *user.User `bun:"rel:has-one,join:author_id=id" json:"-"`
	Comments []Comment  `bun:"rel:has-many,join:id=author_id" json:"comments,omitempty"`
}
