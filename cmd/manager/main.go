package main

import (
	"blizzard/blizzard/config"
	"blizzard/blizzard/core"
	"blizzard/blizzard/db"
	"blizzard/blizzard/db/models/user"
	"blizzard/blizzard/logger"
	"blizzard/blizzard/logger/debug"
	"blizzard/blizzard/permission"
	"context"
	"fmt"
	"github.com/spf13/cobra"
	"github.com/uptrace/bun"
	"strconv"
	"strings"
)

func init() {
	logger.Init()
}

func init() {
	config.Load()
	// enable debug mode to trace queries
	config.Config.Debug = true
}

func init() {
	db.Init()
}

func createUser() *cobra.Command {
	c := &cobra.Command{
		Use:   "create_user <handle> <password> <email> <role_id>",
		Short: "create a new user",
		Args:  cobra.ExactArgs(4),
		RunE: func(cmd *cobra.Command, args []string) error {
			r, e := core.HashConfig.HashEncoded([]byte(args[1]))
			if e != nil {
				return e
			}
			roleId, e := strconv.ParseUint(args[3], 10, 16)
			if e != nil {
				return fmt.Errorf("invalid role_id, role_id must be a valid 16-bit integer")
			}
			if e != nil {
				return e
			}
			var roles []user.Role
			e = db.Database.NewSelect().Model(&roles).Where("id = ?", roleId).Column("id").Scan(cmd.Context())
			if e != nil {
				return e
			}
			if len(roles) != 1 {
				return fmt.Errorf("no roles found with that id")
			}
			u := &user.User{
				Handle:   strings.ToLower(args[0]),
				Email:    strings.ToLower(args[2]),
				Password: string(r),
			}
			if name, e := cmd.Flags().GetString("name"); e == nil && name != "" {
				u.DisplayName = name
			}
			if org, e := cmd.Flags().GetString("organization"); e == nil && org != "" {
				u.Organization = org
			}
			if e := db.Database.RunInTx(cmd.Context(), nil, func(ctx context.Context, tx bun.Tx) error {
				_, err := tx.NewInsert().Model(u).Returning("id").Exec(ctx)
				if err != nil {
					return err
				}
				_, err = tx.NewInsert().Model(&user.UserToRole{
					RoleID: roles[0].ID,
					UserID: u.ID,
				}).Exec(ctx)
				if err != nil {
					return err
				}
				return nil
			}); e != nil {
				return e
			}
			return nil
		},
	}
	c.Flags().String("name", "", "display name for created user")
	c.Flags().String("org", "", "affiliated organization of user")
	return c
}

func createRole() *cobra.Command {
	c := &cobra.Command{
		Use:  "create_role <name> <priority> [permissions...]",
		Args: cobra.MinimumNArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			perms := args[1:]
			var permBits permission.Permission = 0
			priority, e := strconv.Atoi(args[1])
			if e != nil {
				return e
			}
			for _, p := range perms {
				permBits |= permission.StringToPermission(strings.ToLower(p))
			}
			r := &user.Role{
				Name:        args[0],
				Priority:    uint16(priority),
				Permissions: permBits,
			}
			if badge, e := cmd.Flags().GetString("badge"); e == nil && badge != "" {
				r.Badge = badge
			}
			if css, e := cmd.Flags().GetString("style"); e == nil && css != "" {
				r.Style = css
			}
			_, err := db.Database.NewInsert().Model(r).Exec(cmd.Context())
			if err != nil {
				return err
			}
			debug.Dump()
			return nil
		},
	}
	c.Flags().String("badge", "", "path to badge image")
	c.Flags().String("style", "", "css style of role")
	return c
}

var cmds = []*cobra.Command{
	createUser(),
	createRole(),
}

func main() {
	root := &cobra.Command{
		Use:   "manager",
		Short: "blizzard database manager",
	}
	root.AddCommand(cmds...)
	_ = root.Execute()
}
