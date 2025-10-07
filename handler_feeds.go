package main

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/mattnickolaus/gator/internal/database"
)

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

func HandlerFeeds(s *state, cmd command) error {
	if numArgs := len(cmd.Args); numArgs != 0 {
		return fmt.Errorf("The feeds command accepts no arguments\n")
	}

	// GetFeeds joins off of users to get the userName
	allFeeds, err := s.db.GetFeeds(context.Background())
	if err != nil {
		return fmt.Errorf("Error getting all feeds: %v", err)
	}

	fmt.Printf("All Feeds:\n")
	fmt.Printf("\tFeed Name\tUrl\tUsername\tCreated At\n")
	for i, feed := range allFeeds {
		formattedCreated := feed.CreatedAt.Format("2006-01-02 15:04:05")

		fmt.Printf("%d\t%s\t%v\t%v\t%v\n", i, feed.Name, feed.Url, feed.UserName, formattedCreated)
	}

	return nil
}
