package user

import (
	"github.com/google/uuid"
	"time"
)

type (
	User struct {
		ID            uuid.UUID `bun:",pk,type:uuid,default:gen_random_uuid()" json:"id"`
		DisplayName   string    `json:"displayName"`
		Handle        string    `bun:",notnull,unique" json:"handle"`
		Email         string    `bun:",notnull,unique" json:"-"`
		EmailVerified bool      `bun:",default:false" json:"emailVerified,omitempty"`
		// TODO: cache avatar hashes
		Avatar        string            `bun:",scanonly" json:"avatar"`
		Password      string            `bun:",notnull" json:"-"`
		Organizations []Organization    `bun:",m2m:org_memberships,join:User=Organization" json:"organizations,omitempty"`
		RegisteredAt  *time.Time        `bun:",notnull,default:current_timestamp" json:"registeredAt,omitempty"`
		ApiKey        string            `json:"-"`
		Connections   []OAuthConnection `bun:",rel:has-many,join:id=user_id" json:"connections,omitempty"`
		Roles         []Role            `bun:",m2m:role_memberships,join:User=Role" json:"roles,omitempty"`
		BannedSince   *time.Time        `bun:",soft_delete" json:"deletedAt,omitempty"`
		Rating        uint16            `bun:",default:0" json:"rating"`
	}
)
