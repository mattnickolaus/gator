package main

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/mattnickolaus/gator/internal/database"
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

func HandlerAddFeed(s *state, cmd command) error {
	if numArgs := len(cmd.Args); numArgs != 2 {
		return fmt.Errorf("addfeed command requires two arguments\n\tname: The name of the feed\n\turl: The url of the feed\n")
	}
	name := cmd.Args[0]
	url := cmd.Args[1]

	currentUserName := s.cfg.CurrentUserName
	currentUser, err := s.db.GetUser(context.Background(), currentUserName)
	if err != nil {
		return fmt.Errorf("Error finding user ID: %v\n", err)
	}

	feed, err := fetchFeed(context.Background(), url)
	if err != nil {
		return err
	}
	feed.unescapeFields()

	feedParam := database.CreateFeedParams{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Name:      name,
		Url:       url,
		UserID:    currentUser.ID,
	}
	createdFeed, err := s.db.CreateFeed(context.Background(), feedParam)
	if err != nil {
		return fmt.Errorf("Error creating feed in DB: %v\n", err)
	}

	fmt.Printf("created new feed record with the following data:\n%+v\n", createdFeed)

	return nil
}
