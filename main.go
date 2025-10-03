package main

import (
	"fmt"
	"log"
	"os"

	"github.com/mattnickolaus/gator/internal/config"
)

type state struct {
	cfg *config.Config
}

func main() {
	c, err := config.Read()
	if err != nil {
		log.Fatalf("Error reading config: %v", err)
	}

	programState := &state{
		cfg: &c,
	}
	cmds := commands{
		handlers: make(map[string]func(*state, command) error),
	}
	cmds.Register("login", HandlerLogins)

	args := os.Args
	if len(args) < 2 {
		log.Fatalf("Error: Not enough arguments were provided\nUseage: cli <command> [args...]")
	}
	// args[0] is the program name so ignore
	cmdName := args[1]
	cmdArgs := args[2:]

	cmd := command{
		Name: cmdName,
		Args: cmdArgs,
	}

	err = cmds.Run(programState, cmd)
	if err != nil {
		fmt.Printf("%v\n", err)
		os.Exit(1)
		log.Fatalf("%v\n", err)
	}

	c, err = config.Read()
	if err != nil {
		log.Fatalf("%v\n", err)
	}
	fmt.Printf("Config:\n %+v\n", c)
}
