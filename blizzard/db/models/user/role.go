package user

import (
	"blizzard/blizzard/permission"
)

type Role struct {
	ID          uint16                `bun:",pk,autoincrement" json:"id"`
	Name        string                `bun:",unique,notnull" json:"name"`
	Permissions permission.Permission `bun:",default:0" json:"permissions"`
	Badge       string                `json:"badge"`
	Style       string                `bun:",default:'background-color:#90a8bb'" json:"style"`
	Priority    uint16                `bun:",notnull,default:1000" json:"priority"`
}
