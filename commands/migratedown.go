package commands

import (
	"github.com/jmoiron/sqlx"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"

	"shark-auth/foundation/database"
)

func newMigrateDownCommand(db *sqlx.DB) *cobra.Command {
	return &cobra.Command{
		Use:   "migratedown",
		Short: "perform database rollback",
		Run: func(cmd *cobra.Command, args []string) {
			err := database.MigrateDown(db)
			if err != nil {
				logrus.WithError(err).Errorf("Error during DB rollback: %+v", err)
			}
		},
	}
}
