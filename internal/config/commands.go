package config

import (
	"fmt"
)

type state struct {
	cfg *Config
}

type Command struct {
	Name string
	Args []string
}

type commands struct {
	handlers map[string]func(*state, Command) error
}

func InitState(cfg *Config) *state {
	return &state{
		cfg: cfg,
	}
}

func InitCommands() commands {
	hndlr := make(map[string]func(*state, Command) error)
	return commands{
		handlers: hndlr,
	}
}

func (cmds *commands) Run(s *state, cmd Command) error {
	handlerFunc, ok := cmds.handlers[cmd.Name]
	if !ok {
		return fmt.Errorf("Error: Command '%v' does not exist in map", cmd.Name)
	}

	err := handlerFunc(s, cmd)
	if err != nil {
		return err
	}

	return nil
}

func (cmds *commands) Register(name string, f func(*state, Command) error) {
	cmds.handlers[name] = f
}

func HandlerLogins(s *state, cmd Command) error {
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
