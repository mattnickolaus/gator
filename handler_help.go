package main

import (
	"fmt"
)

func HandlerHelp(s *state, cmd command) error {
	if numArgs := len(cmd.Args); numArgs != 0 {
		return fmt.Errorf("help command accepts no arguments.")
	}

	helpMessage := `
  ____       _             
 / ___| __ _| |_ ___  _ __ 
| |  _ / _' | __/ _ \| '__|
| |_| | (_| | || (_) | |   
 \____|\__,_|\__\___/|_|   

A CLI RSS aggregator

Usage:
	gator <command> [arguments]

The commands are:

	login <username>
		Log in as a user.

	register <username>
		Register a new user.

	users
		List all registered users.

	addfeed <name> <url>
		Add a new RSS feed.

	feeds
		List all RSS feeds.

	follow <url>
		Follow a feed.

	following
		List all feeds you are following.

	unfollow <url>
		Unfollow a feed.

	browse [limit]
		Browse the latest posts from your followed feeds. Limit is optional and defaults to 2.

	agg <duration>
		Aggregates posts from feeds at a given frequency. Duration format is like '1s', '1m', '1h'.

	setup
		Setup the database schema.

	takedown
		Migrate down the database schema.

	reset
		Deletes all users from the database.

	help
		Shows this help message.
`
	fmt.Println(helpMessage)

	return nil
}
