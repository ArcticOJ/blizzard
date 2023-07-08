package user

import (
	"blizzard/blizzard/db/utils"
	"context"
	"github.com/google/uuid"
	"github.com/uptrace/bun"
)

type (
	User struct {
		ID            uuid.UUID         `bun:",pk,unique,type:uuid,default:gen_random_uuid()" json:"id"`
		DisplayName   string            `json:"displayName,omitempty"`
		Handle        string            `bun:",notnull,unique" json:"handle"`
		Email         string            `bun:",notnull,unique" json:"email,omitempty"`
		EmailVerified bool              `bun:",default:false" json:"emailVerified,omitempty"`
		Avatar        string            `bun:"-" json:"avatar"`
		Password      string            `bun:",notnull" json:"password,omitempty"`
		Organization  string            `json:"organization,omitempty"`
		RegisteredAt  *bun.NullTime     `bun:",nullzero,type:timestamptz,notnull,default:'now()'::timestamptz" json:"registeredAt,omitempty"`
		ApiKey        string            `json:"-"`
		Connections   []OAuthConnection `bun:",rel:has-many,join:id=user_id" json:"connections,omitempty"`
		Roles         []Role            `bun:",m2m:user_to_roles,join:User=Role" json:"roles,omitempty"`
		DeletedAt     *bun.NullTime     `bun:",soft_delete,nullzero" json:"deletedAt,omitempty"`
		Rating        uint16            `bun:",default:1000" json:"rating"`
	}

	UserToRole struct {
		RoleID uint16    `bun:",pk"`
		Role   *Role     `bun:"rel:belongs-to,join:role_id=id"`
		UserID uuid.UUID `bun:",pk,type:uuid"`
		User   *User     `bun:"rel:belongs-to,join:user_id=id"`
	}
)

func (UserToRole) BeforeCreateTable(_ context.Context, query *bun.CreateTableQuery) error {
	utils.Cascade(query, "user_id", "users", "id")
	utils.Cascade(query, "role_id", "roles", "id")
	return nil
}
