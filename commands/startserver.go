package commands

import (
	"github.com/go-redis/redis/v7"
	"github.com/jmoiron/sqlx"
	"github.com/spf13/cobra"

	"shark-auth/apis"
)

func newStartServerCommand(db *sqlx.DB, redisClient *redis.Client) *cobra.Command {
	return &cobra.Command{
		Use:     "startserver",
		Short:   "Start HTTP API server",
		Aliases: []string{"startapp", "runserver"},
		Run: func(_ *cobra.Command, _ []string) {
			router := apis.API(db, redisClient)
			err := router.Run(":8080")
			if err != nil {
				panic("starting server failed")
			}
		},
	}
}
