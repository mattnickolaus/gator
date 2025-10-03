package main

import (
	"fmt"
)

func HandlerLogins(s *state, cmd command) error {
	if numArgs := len(cmd.Args); numArgs == 0 || numArgs > 1 {
		return fmt.Errorf("Error: the login expects a single argument, the username")
	}

	err := s.cfg.SetUser(cmd.Args[0])
	if err != nil {
		return err
	}

	fmt.Printf("The user %s has been set", s.cfg.CurrentUserName)
	return nil
}
