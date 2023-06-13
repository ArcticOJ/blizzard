package users

import (
	"github.com/google/uuid"
	"github.com/uptrace/bun"
	"time"
)

type User struct {
	bun.BaseModel
	UUID         uuid.UUID         `bun:",pk,unique,type:uuid,default:gen_random_uuid()" json:"uuid"`
	DisplayName  string            `json:"displayName,omitempty"`
	Handle       string            `bun:",notnull,unique" json:"handle"`
	Email        string            `bun:",notnull,unique" json:"email"`
	Password     string            `bun:",notnull" json:"password,omitempty"`
	Organization string            `json:"organization,omitempty"`
	RegisteredAt time.Time         `bun:",nullzero,type:timestamptz,notnull,default:'now()'::timestamptz" json:"registeredAt,omitempty"`
	ApiKey       string            `json:"-"`
	Avatar       string            `bun:"-" json:"avatar"`
	Connections  []OAuthConnection `bun:",rel:has-many,join:uuid=user_id" json:"-"`
	//DeletedAt    time.Time         `bun:",soft_delete,nullzero"`
	//LastLogin     time.Time `bun:"lastLogin,nullzero,type:timestamptz" json:"lastLogin"`
}
