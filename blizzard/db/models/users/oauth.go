package users

import (
	"blizzard/blizzard/db/utils"
	"context"
	"github.com/google/uuid"
	"github.com/uptrace/bun"
)

type OAuthConnection struct {
	ID       string    `bun:",pk"`
	Provider string    `bun:",notnull"`
	Username string    `bun:",notnull"`
	UserID   uuid.UUID `bun:",type:uuid"`
	//	User     *User     `bun:"rel:belongs-to,join:user_id=uuid"`
}

func (OAuthConnection) BeforeCreateTable(_ context.Context, query *bun.CreateTableQuery) error {
	return utils.Cascade(query, "user_id", "users", "uuid")
}
