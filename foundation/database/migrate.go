package database

import (
	"fmt"

	"github.com/pkg/errors"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq" // The database driver in use.
	"github.com/sirupsen/logrus"
)

func MigrateUp(db *sqlx.DB) error {
	m, err := getMigrate(db)
	if err != nil {
		return errors.WithMessage(err, "getting migration object failed")
	}

	err = m.Up()
	if err != nil {
		if err == migrate.ErrNoChange {
			logrus.Warn("no changes")
			return nil
		}
		logrus.WithError(err).Fatalf("migration failed")
		return err
	}
	logrus.Infof("migration successful")

	return nil
}

// Does th rollback. The DB will go to the previous state.
func MigrateDown(db *sqlx.DB) error {
	m, err := getMigrate(db)
	if err != nil {
		return errors.WithMessage(err, "getting migration object failed")
	}

	err = m.Down()
	if err != nil {
		if err == migrate.ErrNoChange {
			logrus.Info("no changes")
			return nil
		}

		logrus.WithError(err).
			WithField("stack_trace", fmt.Sprintf("%+v", err)).
			Fatal("rollback failed")
	}
	logrus.Info("rollback successful")
	return nil
}

func getMigrate(db *sqlx.DB) (*migrate.Migrate, error) {
	driver, err := postgres.WithInstance(db.DB, &postgres.Config{})
	if err != nil {
		return nil, errors.WithStack(err)
	}
	return migrate.NewWithDatabaseInstance(
		"file://migrations",
		"shark-auth", driver) //todo remove hardcoding
}
