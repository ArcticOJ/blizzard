package shared

import (
	"github.com/google/uuid"
	"github.com/uptrace/bun"
	"time"
)

type User struct {
	bun.BaseModel `bun:"table:User"`
	ID            uuid.UUID `bun:"id,pk,unique,type:uuid,default:gen_random_uuid()" json:"id"`
	DisplayName   string    `bun:"displayName" json:"displayName,omitempty"`
	Handle        string    `bun:"handle,notnull,unique" json:"handle"`
	Email         string    `bun:"email,notnull,unique" json:"email"`
	Password      string    `bun:"password,notnull" json:"password,omitempty"`
	Organization  string    `bun:"organization" json:"organization,omitempty"`
	RegisteredAt  time.Time `bun:"registeredAt,nullzero,type:timestamptz,notnull,default:'now()'::timestamptz" json:"registeredAt,omitempty"`
	ApiKey        string    `bun:"apiKey" json:"apiKey,omitempty"`
	//LastLogin     time.Time `bun:"lastLogin,nullzero,type:timestamptz" json:"lastLogin"`
}
