package user

import (
	"github.com/google/uuid"
	"time"
)

type (
	MinimalUser struct {
		ID           string      `json:"id,omitempty"`
		DisplayName  string      `json:"displayName,omitempty"`
		Handle       string      `json:"handle,omitempty"`
		Avatar       string      `json:"avatar,omitempty"`
		Organization string      `json:"organization,omitempty"`
		TopRole      interface{} `json:"topRole,omitempty"`
		Rating       uint16      `json:"rating"`
	}
	User struct {
		ID            uuid.UUID         `bun:",pk,unique,type:uuid,default:gen_random_uuid()" json:"id"`
		DisplayName   string            `json:"displayName"`
		Handle        string            `bun:",notnull,unique" json:"handle"`
		Email         string            `bun:",notnull,unique" json:"-"`
		EmailVerified bool              `bun:",default:false" json:"emailVerified,omitempty"`
		Avatar        string            `bun:"-" json:"avatar"`
		Password      string            `bun:",notnull" json:"-"`
		Organization  string            `json:"organization"`
		RegisteredAt  *time.Time        `bun:",nullzero,notnull,default:'now()'" json:"registeredAt,omitempty"`
		ApiKey        string            `json:"-"`
		Connections   []OAuthConnection `bun:"rel:has-many,join:id=user_id" json:"connections"`
		Roles         []Role            `bun:"m2m:user_to_roles,join:User=Role" json:"roles"`
		TopRole       *Role             `bun:"-" json:"topRole"`
		DeletedAt     *time.Time        `bun:",soft_delete,nullzero" json:"deletedAt,omitempty"`
		Rating        uint16            `bun:",default:0" json:"rating"`
	}

	UserToRole struct {
		RoleID uint16    `bun:",pk"`
		Role   *Role     `bun:"rel:belongs-to,join:role_id=id"`
		UserID uuid.UUID `bun:",pk,type:uuid"`
		User   *User     `bun:"rel:belongs-to,join:user_id=id"`
	}
)
