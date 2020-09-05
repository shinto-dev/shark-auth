package commands

import (
	"github.com/go-redis/redis/v7"
	"github.com/jmoiron/sqlx"
	"github.com/spf13/cobra"
)

func NewCLI(db *sqlx.DB, redisClient *redis.Client) *cobra.Command {
	cli := &cobra.Command{
		Use:   "shark-auth",
		Short: "shark-auth is a decentralized authentication service",
	}

	cli.AddCommand(newStartServerCommand(db, redisClient))
	cli.AddCommand(newMigrateUpCommand(db))
	cli.AddCommand(newMigrateDownCommand(db))

	return cli
}
