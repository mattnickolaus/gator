package main

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/mattnickolaus/gator/internal/database"
)

func middlewareLoggedIn(handler func(s *state, cmd command, user database.User) error) func(*state, command) error {
	return func(s *state, cmd command) error {
		currentUser, err := s.db.GetUser(context.Background(), s.cfg.CurrentUserName)
		if err != nil {
			return fmt.Errorf("Error finding user ID: %v\n", err)
		}
		return handler(s, cmd, currentUser)
	}
}

func HandlerLogins(s *state, cmd command) error {
	if numArgs := len(cmd.Args); numArgs == 0 || numArgs > 1 {
		return fmt.Errorf("Error: the login expects a single argument, the username")
	}
	userName := cmd.Args[0]

	response, err := s.db.GetUser(context.Background(), userName)
	if err != nil {
		return fmt.Errorf("Error: User login with name '%s' was not found\n%v", userName, err)
	}

	err = s.cfg.SetUser(response.Name)
	if err != nil {
		return err
	}

	fmt.Printf("The user %s has been set\n", s.cfg.CurrentUserName)
	return nil
}

func HanlderRegister(s *state, cmd command) error {
	if numArgs := len(cmd.Args); numArgs == 0 || numArgs > 1 {
		return fmt.Errorf("Error: the register expects a single argument, the username")
	}
	newUserName := cmd.Args[0]

	// Should return err = sql: no rows in result set
	response, err := s.db.GetUser(context.Background(), newUserName)
	if err != sql.ErrNoRows {
		return fmt.Errorf("Error: User '%s' already exists\n", response.Name)
	}

	regParam := database.CreateUserParams{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Name:      newUserName,
	}

	user, err := s.db.CreateUser(context.Background(), regParam)
	if err != nil {
		return fmt.Errorf("Error: %v", err)
	}

	err = s.cfg.SetUser(user.Name)
	if err != nil {
		return fmt.Errorf("Error %v", err)
	}

	fmt.Printf("New user with name '%v' was created\n", user.Name)

	return nil
}

func HandlerReset(s *state, cmd command) error {
	if numArgs := len(cmd.Args); numArgs > 0 {
		return fmt.Errorf("reset command accepts no arguments")
	}

	err := s.db.DeleteAllUsers(context.Background())
	if err != nil {
		return fmt.Errorf("Error: %v", err)
	}

	fmt.Printf("Reset successful: All users deleted from the database\n")

	return nil
}

func HandlerUsers(s *state, cmd command) error {
	if numArgs := len(cmd.Args); numArgs > 0 {
		return fmt.Errorf("users command accepts NO arguments")
	}

	users, err := s.db.GetUsers(context.Background())
	if err != nil {
		return fmt.Errorf("Error: %v", err)
	}

	if len(users) == 0 {
		return fmt.Errorf("Error: There are currently no registered users\n\tRun gator register <username> to setup a user")
	}

	for _, u := range users {
		userLine := fmt.Sprintf("* %s", u)
		if u == s.cfg.CurrentUserName {
			userLine += " (current)"
		}
		fmt.Println(userLine)
	}

	return nil

}
