package user

import (
	"context"
	"github.com/ArcticOJ/blizzard/v0/db/utils"
	"github.com/google/uuid"
	"github.com/uptrace/bun"
)

type OAuthConnection struct {
	ID            string    `bun:",pk" json:"-"`
	Provider      string    `bun:",pk,unique:provider,notnull" json:"provider"`
	Username      string    `bun:",notnull" json:"username"`
	UserID        uuid.UUID `bun:",type:uuid,unique:provider" json:"-"`
	ShowInProfile bool      `bun:",default:true" json:"-"`
}

func (OAuthConnection) BeforeCreateTable(_ context.Context, query *bun.CreateTableQuery) error {
	return utils.Cascade(query, "user_id", "users", "id")
}