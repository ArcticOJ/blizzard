package user

import (
	"context"
	"github.com/ArcticOJ/blizzard/v0/db/utils"
	"github.com/google/uuid"
	"github.com/uptrace/bun"
)

type OAuthConnection struct {
	Provider      string    `bun:",pk,notnull" json:"provider"`
	Username      string    `bun:",notnull" json:"username"`
	UserID        uuid.UUID `bun:",pk,type:uuid,notnull" json:"-"`
	ShowInProfile bool      `bun:",default:true" json:"-"`
}

func (OAuthConnection) BeforeCreateTable(_ context.Context, query *bun.CreateTableQuery) error {
	return utils.Cascade(query, "user_id", "users", "id")
}
