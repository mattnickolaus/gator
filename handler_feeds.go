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

	followParam := database.CreateFeedFollowParams{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		UserID:    currentUser.ID,
		FeedID:    createdFeed.ID,
	}

	followedRowData, err := s.db.CreateFeedFollow(context.Background(), followParam)
	if err != nil {
		return fmt.Errorf("Error writing the feed follow: %v\n", err)
	}
	fmt.Printf("created new feed record with the following data:\n%+v\n", createdFeed)
	fmt.Printf("You (%v) are now following this feed\n", followedRowData.UserName)

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

func HandlerFeedFollow(s *state, cmd command) error {
	if numArgs := len(cmd.Args); numArgs != 1 {
		return fmt.Errorf("The follow command accepts a single argument:\n\tUrl: the url to the chosen feed to follow")
	}
	feedsUrl := cmd.Args[0]

	chosenFeed, err := s.db.GetFeedByUrl(context.Background(), feedsUrl)
	if err != nil {
		return fmt.Errorf("Error: Unable to retrieve feed with that url\n%v\n", err)
	}

	currentUserName := s.cfg.CurrentUserName
	currentUser, err := s.db.GetUser(context.Background(), currentUserName)
	if err != nil {
		return fmt.Errorf("Error finding user ID: %v\n", err)
	}

	followParam := database.CreateFeedFollowParams{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		UserID:    currentUser.ID,
		FeedID:    chosenFeed.ID,
	}

	followedRowData, err := s.db.CreateFeedFollow(context.Background(), followParam)
	if err != nil {
		return fmt.Errorf("Error writing the feed follow: %v\n", err)
	}

	formattedCreated := followedRowData.CreatedAt.Format("2006-01-02 15:04:05")
	fmt.Printf("Follow Record Created:\nUser: %v\tFeed: %v\t%v\n", followedRowData.UserName, followedRowData.FeedName, formattedCreated)

	return nil
}

func HandlerFollowing(s *state, cmd command) error {
	if numArgs := len(cmd.Args); numArgs != 0 {
		return fmt.Errorf("The following command accepts no arguments")
	}

	currentUserName := s.cfg.CurrentUserName
	currentUser, err := s.db.GetUser(context.Background(), currentUserName)
	if err != nil {
		return fmt.Errorf("Error finding user ID: %v\n", err)
	}

	following, err := s.db.GetFeedFollowsForUser(context.Background(), currentUser.ID)
	if err != nil {
		return fmt.Errorf("Error retrieving followed feeds: %v\n", err)
	}

	fmt.Printf("You (%v) are following these feeds:\n", currentUser.Name)
	for _, feed := range following {
		fmt.Printf("\t* %v\n", feed.FeedName)
	}

	return nil
}
