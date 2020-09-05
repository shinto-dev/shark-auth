package main

import (
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"

	"shark-auth/commands"
	"shark-auth/foundation/config"
	"shark-auth/foundation/database"
	"shark-auth/foundation/redisclient"
)

func main() {
	err := run()
	if err != nil {
		logrus.Panic(err)
	}
}

func run() error {
	db, err := database.Open(database.Config{
		User:       config.GetString("database.user"),
		Password:   config.GetStringOrDefault("database.password", ""),
		Host:       config.GetString("database.host"),
		Name:       config.GetStringOrDefault("database.name", "shark-auth"),
		DisableTLS: true,
	})
	if err != nil {
		return errors.Wrap(err, "error opening DB connection")
	}

	redisClient, err := redisclient.NewRedisClient(redisclient.Config{
		Host: config.GetStringOrDefault("redis.host", "localhost"),
		Port: config.GetInt("redis.port"),
	})
	if err != nil {
		return errors.Wrap(err, "error opening redis connection")
	}

	cli := commands.NewCLI(db, redisClient)
	//cli.Version = fmt.Sprintf("%s (Commit: %s)", version, commit)
	err = cli.Execute()
	if err != nil {
		return errors.Wrap(err, "error initializing the command")
	}
	return nil
}
