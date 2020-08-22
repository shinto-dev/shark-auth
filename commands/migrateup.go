package commands

import (
	"github.com/jmoiron/sqlx"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"

	"shark-auth/foundation/database"
)

func newMigrateUpCommand(db *sqlx.DB) *cobra.Command {
	return &cobra.Command{
		Use:   "migrateup",
		Short: "Perform database migration",
		Long:  "",
		Run: func(cmd *cobra.Command, args []string) {
			err := database.MigrateUp(db)
			if err != nil {
				logrus.WithError(err).Errorf("Error during DB migrations: %+v", err)
			}
		},
	}
}
