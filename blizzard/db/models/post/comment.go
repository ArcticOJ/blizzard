package post

import (
	"blizzard/blizzard/db/models/user"
	"github.com/google/uuid"
	"time"
)

type Comment struct {
	ID          uint32     `bun:",pk,autoincrement" json:"id"`
	CommentedAt *time.Time `bun:",nullzero,type:timestamptz,notnull,default:'now()'::timestamptz" json:"commentedAt,omitempty"`
	AuthorID    uuid.UUID  `bun:",type:uuid" json:"-"`
	Author      *user.User `bun:"rel:has-one,join:author_id=id" json:"author,omitempty"`
	PostID      uint32
	Post        *Post `bun:"rel:belongs-to,join:post_id=id"`
	Content     string
}
