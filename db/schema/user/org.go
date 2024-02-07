package user

import (
	"github.com/google/uuid"
	"time"
)

type (
	Organization struct {
		ID          string `bun:",pk" json:"id"`
		Name        string `json:"name"`
		Description string `json:"description,omitempty"`
		Members     []User `bun:",m2m:org_memberships,join:Organization=User" json:"members,omitempty"`
	}

	OrgRole uint8

	OrgMembership struct {
		OrgID        string        `bun:",pk"`
		Organization *Organization `bun:",rel:belongs-to,join:org_id=id"`
		Role         OrgRole       `bun:",default:0"`
		JoinedAt     time.Time     `bun:",notnull,default:current_timestamp" json:"joinedAt"`
		UserID       uuid.UUID     `bun:",pk,type:uuid"`
		User         *User         `bun:",rel:belongs-to,join:user_id=id"`
	}
)

const (
	OrgMember OrgRole = iota
	OrgAdmin
	OrgOwner
)
