package user

import (
	"github.com/google/uuid"
	"time"
)

type (
	User struct {
		ID                 uuid.UUID         `bun:",pk,unique,type:uuid,default:gen_random_uuid()" json:"id,omitempty"`
		DisplayName        string            `json:"displayName,omitempty"`
		Handle             string            `bun:",notnull,unique" json:"handle,omitempty"`
		Email              string            `bun:",notnull,unique" json:"email,omitempty"`
		EmailVerified      bool              `bun:",default:false" json:"emailVerified,omitempty"`
		Avatar             string            `bun:"-" json:"avatar,omitempty"`
		Password           string            `bun:",notnull" json:"password,omitempty"`
		Organization       string            `json:"organization,omitempty"`
		RegisteredAt       *time.Time        `bun:",nullzero,type:timestamptz,notnull,default:'now()'::timestamptz" json:"registeredAt,omitempty"`
		ApiKey             string            `json:"-"`
		Connections        []OAuthConnection `bun:"rel:has-many,join:id=user_id" json:"connections,omitempty"`
		Roles              []Role            `bun:"m2m:user_to_roles,join:User=Role" json:"roles,omitempty"`
		ProblemsSolved     uint16            `bun:",scanonly" json:"problemsSolved,omitempty"`
		DeletedAt          *time.Time        `bun:",soft_delete,nullzero" json:"deletedAt,omitempty"`
		Rating             uint16            `bun:",default:0" json:"rating"`
		LastUsernameChange *time.Time        `bun:",nullzero,type:timestamptz" json:"lastUsernameChange,omitempty"`
	}

	UserToRole struct {
		RoleID uint16    `bun:",pk"`
		Role   *Role     `bun:"rel:belongs-to,join:role_id=id"`
		UserID uuid.UUID `bun:",pk,type:uuid"`
		User   *User     `bun:"rel:belongs-to,join:user_id=id"`
	}
)
