package main

import (
	"fmt"
	"os"

	"github.com/mattnickolaus/gator/internal/config"
)

func main() {
	c, err := config.Read()
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}

	s := config.InitState(&c)
	cmds := config.InitCommands()
	cmds.Register("login", config.HandlerLogins)

	args := os.Args
	if len(args) < 2 {
		fmt.Printf("Error: Not enough arguments were provided\n")
		os.Exit(1)
	}
	// args[0] is the program name so ignore
	cmdName := args[1]
	var cmdArgs []string
	if len(args) > 2 {
		cmdArgs = args[2:]
	}
	cmd := config.Command{
		Name: cmdName,
		Args: cmdArgs,
	}

	err = cmds.Run(s, cmd)
	if err != nil {
		fmt.Printf("%v\n", err)
		os.Exit(1)
	}

	c, err = config.Read()
	if err != nil {
		fmt.Printf("%v\n", err)
		os.Exit(1)
	}
	fmt.Printf("Config:\n %+v\n", c)
}
