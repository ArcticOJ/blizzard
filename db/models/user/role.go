package user

import (
	"blizzard/permission"
)

type Role struct {
	ID          uint16                `bun:",pk,autoincrement" json:"id,omitempty"`
	Name        string                `bun:",unique,notnull" json:"name,omitempty"`
	Permissions permission.Permission `bun:",default:0" json:"permissions,omitempty"`
	Icon        string                `json:"icon"`
	Color       string                `json:"color"`
	Priority    uint16                `bun:",notnull,default:1000" json:"priority,omitempty"`
	Members     []User                `bun:"m2m:user_to_roles,join:Role=User" json:"members,omitempty"`
}
