package main

import (
	"fmt"
)

type command struct {
	Name string
	Args []string
}

type commands struct {
	handlers map[string]func(*state, command) error
}

func (cmds *commands) Run(s *state, cmd command) error {
	handlerFunc, ok := cmds.handlers[cmd.Name]
	if !ok {
		return fmt.Errorf("Error: command '%v' does not exist in map", cmd.Name)
	}

	err := handlerFunc(s, cmd)
	if err != nil {
		return err
	}

	return nil
}

func (cmds *commands) Register(name string, f func(*state, command) error) {
	cmds.handlers[name] = f
}
