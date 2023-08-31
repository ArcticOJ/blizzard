package main

import (
	"blizzard/blizzard/config"
	"blizzard/blizzard/db"
	"blizzard/blizzard/logger"
	"blizzard/migrations"
	"github.com/spf13/cobra"
	"github.com/uptrace/bun/dbfixture"
	"github.com/uptrace/bun/migrate"
	"os"
	"strings"
)

var migrator *migrate.Migrator

func _init() *cobra.Command {
	return &cobra.Command{
		Use:   "init",
		Short: "create migration tables",
		RunE: func(cmd *cobra.Command, args []string) error {
			return migrator.Init(cmd.Context())
		},
	}
}

func _migrate() *cobra.Command {
	return &cobra.Command{
		Use:   "migrate",
		Short: "migrate database",
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := migrator.Lock(cmd.Context()); err != nil {
				return err
			}
			defer migrator.Unlock(cmd.Context())
			group, err := migrator.Migrate(cmd.Context())
			if err != nil {
				return err
			}
			if group.IsZero() {
				logger.Logger.Info().Msg("there are no new migrations to run, database is up to date.")
				return nil
			}
			logger.Logger.Info().Msgf("migrated to %s", group)
			return nil
		},
	}
}

func rollback() *cobra.Command {
	return &cobra.Command{
		Use:   "rollback",
		Short: "rollback the last migration group",
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := migrator.Lock(cmd.Context()); err != nil {
				return err
			}
			defer migrator.Unlock(cmd.Context())
			group, err := migrator.Rollback(cmd.Context())
			if err != nil {
				return err
			}
			if group.IsZero() {
				logger.Logger.Info().Msg("there are no groups to roll back")
				return nil
			}
			logger.Logger.Info().Msgf("rolled back %s", group)
			return nil
		},
	}
}

func lock() *cobra.Command {
	return &cobra.Command{
		Use:   "lock",
		Short: "lock migrations",
		RunE: func(cmd *cobra.Command, args []string) error {
			return migrator.Lock(cmd.Context())
		},
	}
}

func unlock() *cobra.Command {
	return &cobra.Command{
		Use:   "unlock",
		Short: "unlock migrations",
		RunE: func(cmd *cobra.Command, args []string) error {
			return migrator.Unlock(cmd.Context())
		},
	}
}

func createGo() *cobra.Command {
	return &cobra.Command{
		Use:   "create_go",
		Short: "create Go migration",
		RunE: func(cmd *cobra.Command, args []string) error {
			name := strings.Join(args, "_")
			mf, err := migrator.CreateGoMigration(cmd.Context(), name)
			if err != nil {
				return err
			}
			logger.Logger.Info().Msgf("created migration %s (%s)", mf.Name, mf.Path)
			return nil
		},
	}
}

func createSQL() *cobra.Command {
	return &cobra.Command{
		Use:   "create_sql",
		Short: "create up and down SQL migrations",
		RunE: func(cmd *cobra.Command, args []string) error {
			name := strings.Join(args, "_")
			files, err := migrator.CreateSQLMigrations(cmd.Context(), name)
			if err != nil {
				return err
			}
			for _, mf := range files {
				logger.Logger.Info().Msgf("created migration %s (%s)", mf.Name, mf.Path)
			}
			return nil
		},
	}
}
func status() *cobra.Command {
	return &cobra.Command{
		Use:   "status",
		Short: "print migrations status",
		RunE: func(cmd *cobra.Command, args []string) error {
			ms, err := migrator.MigrationsWithStatus(cmd.Context())
			if err != nil {
				return err
			}
			logger.Logger.Info().Msgf(
				"migrations: %s\n"+
					"unapplied migrations: %s\n"+
					"last migration group: %s", ms, ms.Unapplied(), ms.LastGroup())
			return nil
		},
	}
}

func markApplied() *cobra.Command {
	return &cobra.Command{
		Use:   "mark_applied",
		Short: "mark migrations as applied without actually running them",
		RunE: func(cmd *cobra.Command, args []string) error {
			group, err := migrator.Migrate(cmd.Context(), migrate.WithNopMigration())
			if err != nil {
				return err
			}
			if group.IsZero() {
				logger.Logger.Info().Msg("there are no new migrations to mark as applied")
				return nil
			}
			logger.Logger.Info().Msgf("marked as applied %s", group)
			return nil
		},
	}
}

func reset() *cobra.Command {
	return &cobra.Command{
		Use:   "reset",
		Short: "recreate all tables and seed with example data",
		RunE: func(cmd *cobra.Command, args []string) error {
			if e := migrator.Reset(cmd.Context()); e != nil {
				return e
			}
			if _, e := migrator.Migrate(cmd.Context()); e != nil {
				return e
			}
			fixture := dbfixture.New(migrator.DB())
			return fixture.Load(cmd.Context(), os.DirFS("cmd/migrator"), "fixture.yml")
		},
	}
}

var cmds = []*cobra.Command{
	_init(),
	_migrate(),
	rollback(),
	lock(),
	unlock(),
	createGo(),
	createSQL(),
	status(),
	markApplied(),
	reset(),
}

func main() {
	config.Config.Debug = true
	migrator = migrate.NewMigrator(db.Database, migrations.Migrations)
	root := cobra.Command{
		Use:   "migrator",
		Short: "blizzard migration helper",
	}
	root.AddCommand(cmds...)
	_ = root.Execute()
}
