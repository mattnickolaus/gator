package main

import (
	"database/sql"
	"log"
	"os"

	_ "github.com/lib/pq"
	"github.com/mattnickolaus/gator/internal/config"
	"github.com/mattnickolaus/gator/internal/database"
)

type state struct {
	db  *database.Queries
	cfg *config.Config
}

func main() {
	c, err := config.Read()
	if err != nil {
		log.Fatalf("Error reading config: %v", err)
	}

	db, err := sql.Open("postgres", c.DbURL)
	dbQueries := database.New(db)

	programState := &state{
		db:  dbQueries,
		cfg: &c,
	}

	cmds := commands{
		handlers: make(map[string]func(*state, command) error),
	}
	cmds.Register("login", HandlerLogins)
	cmds.Register("register", HanlderRegister)
	cmds.Register("reset", HandlerReset)
	cmds.Register("users", HandlerUsers)
	cmds.Register("agg", HandlerAgg)
	cmds.Register("addfeed", middlewareLoggedIn(HandlerAddFeed))
	cmds.Register("feeds", HandlerFeeds)
	cmds.Register("follow", middlewareLoggedIn(HandlerFeedFollow))
	cmds.Register("following", middlewareLoggedIn(HandlerFollowing))
	cmds.Register("unfollow", middlewareLoggedIn(HandlerUnfollow))
	cmds.Register("browse", middlewareLoggedIn(HandlerBrowse))

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
		log.Fatalf("%v\n", err)
	}
}
