package main

import (
	"context"
	"fmt"
)

func HandlerAgg(s *state, cmd command) error {
	if numArgs := len(cmd.Args); numArgs > 0 {
		return fmt.Errorf("agg command accepts no arguments")
	}

	feed, err := fetchFeed(context.Background(), "https://www.wagslane.dev/index.xml")
	if err != nil {
		return err
	}
	feed.unescapeFields()

	fmt.Printf("Feed found with the following data:\n%+v", feed)

	return nil
}
