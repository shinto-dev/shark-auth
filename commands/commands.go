package commands

import (
	"github.com/jmoiron/sqlx"
	"github.com/spf13/cobra"
)

func NewCLI(db *sqlx.DB) *cobra.Command {
	cli := &cobra.Command{
		Use:   "shark-auth",
		Short: "shark-auth is a decentralized authentication service",
	}

	cli.AddCommand(newStartServerCommand(db))
	cli.AddCommand(newMigrateUpCommand(db))
	cli.AddCommand(newMigrateDownCommand(db))

	return cli
}
