package main

import (
	"fmt"

	"github.com/mattnickolaus/gator/sql"
)

func HandlerSetUp(s *state, cmd command) error {
	if numArgs := len(cmd.Args); numArgs != 0 {
		return fmt.Errorf("setup command accepts 0 arguments\n")
	}

	sql.UpMigrations(s.dbConn)

	fmt.Printf("Postgres gator schema setup successfully!\n")

	return nil
}

func HandlerTakeDown(s *state, cmd command) error {
	if numArgs := len(cmd.Args); numArgs != 0 {
		return fmt.Errorf("takedown command accepts 0 arguments\n")
	}

	sql.DownMigrations(s.dbConn)

	fmt.Printf("Postgres gator schema migrated down/removed successfully!\n")

	return nil
}
