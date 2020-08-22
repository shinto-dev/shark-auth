package main

import (
	"shark-auth/commands"
	"shark-auth/foundation/database"
)

func main() {
	db, err := database.Open(database.Config{
		User:       "",
		Password:   "",
		Host:       "localhost",
		Name:       "shark-auth",
		DisableTLS: true,
	})
	if err != nil {
		panic("error opening DB connection")
	}

	cli := commands.NewCLI(db)
	//cli.Version = fmt.Sprintf("%s (Commit: %s)", version, commit)
	err = cli.Execute()
	if err != nil {
		panic("error initializing the command")
	}
}
