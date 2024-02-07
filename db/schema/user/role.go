package user

import (
	"github.com/ArcticOJ/blizzard/v0/permission"
	"github.com/google/uuid"
)

type (
	Role struct {
		ID          uint16                `bun:",pk,autoincrement" json:"id,omitempty"`
		Name        string                `bun:",unique,notnull" json:"name,omitempty"`
		Permissions permission.Permission `bun:",default:0" json:"permissions,omitempty"`
		Icon        string                `json:"icon"`
		Color       string                `json:"color"`
		Priority    uint16                `bun:",notnull,default:1000,unique" json:"priority,omitempty"`
		Members     []User                `bun:",m2m:role_memberships,join:Role=User" json:"members,omitempty"`
	}

	RoleMembership struct {
		RoleID uint16    `bun:",pk"`
		Role   *Role     `bun:",rel:belongs-to,join:role_id=id"`
		UserID uuid.UUID `bun:",pk,type:uuid"`
		User   *User     `bun:",rel:belongs-to,join:user_id=id"`
	}
)
