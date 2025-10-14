package main

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/mattnickolaus/gator/internal/database"
)

func HandlerAddFeed(s *state, cmd command, user database.User) error {
	if numArgs := len(cmd.Args); numArgs != 2 {
		return fmt.Errorf("addfeed command requires two arguments\n\tname: The name of the feed\n\turl: The url of the feed\n")
	}
	name := cmd.Args[0]
	url := cmd.Args[1]

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
		UserID:    user.ID,
	}
	createdFeed, err := s.db.CreateFeed(context.Background(), feedParam)
	if err != nil {
		return fmt.Errorf("Error creating feed in DB: %v\n", err)
	}

	followParam := database.CreateFeedFollowParams{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		UserID:    user.ID,
		FeedID:    createdFeed.ID,
	}

	followedRowData, err := s.db.CreateFeedFollow(context.Background(), followParam)
	if err != nil {
		return fmt.Errorf("Error writing the feed follow: %v\n", err)
	}
	fmt.Printf("created new feed record with the following data:\n")
	fmt.Printf("\tFeed Name\tUrl\tCreated At\n")
	fmt.Printf("\t%s\t%v\t%v\n", createdFeed.Name, createdFeed.Url, createdFeed.CreatedAt.Format("2006-01-02 15:04:05"))

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

func HandlerFeedFollow(s *state, cmd command, user database.User) error {
	if numArgs := len(cmd.Args); numArgs != 1 {
		return fmt.Errorf("The follow command accepts a single argument:\n\tUrl: the url to the chosen feed to follow")
	}
	feedsUrl := cmd.Args[0]

	chosenFeed, err := s.db.GetFeedByUrl(context.Background(), feedsUrl)
	if err != nil {
		return fmt.Errorf("Error: Unable to retrieve feed with that url\n%v\n", err)
	}

	followParam := database.CreateFeedFollowParams{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		UserID:    user.ID,
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

func HandlerUnfollow(s *state, cmd command, user database.User) error {
	if numArgs := len(cmd.Args); numArgs != 1 {
		return fmt.Errorf("The unfollow command accepts a single argument:\n\tUrl: the url to the chosen feed to unfollow")
	}
	feedsUrl := cmd.Args[0]

	chosenFeed, err := s.db.GetFeedByUrl(context.Background(), feedsUrl)
	if err != nil {
		return fmt.Errorf("Error: Unable to retrieve feed with that url\n%v\n", err)
	}

	deleteParam := database.DeleteFeedFollowByIDsParams{
		UserID: user.ID,
		FeedID: chosenFeed.ID,
	}

	err = s.db.DeleteFeedFollowByIDs(context.Background(), deleteParam)
	if err != nil {
		return fmt.Errorf("Error deleting recrod from the database: %v", err)
	}
	fmt.Printf("Successfully unfollowed %v", chosenFeed.Name)

	return nil
}

func HandlerFollowing(s *state, cmd command, user database.User) error {
	if numArgs := len(cmd.Args); numArgs != 0 {
		return fmt.Errorf("The following command accepts no arguments")
	}

	following, err := s.db.GetFeedFollowsForUser(context.Background(), user.ID)
	if err != nil {
		return fmt.Errorf("Error retrieving followed feeds: %v\n", err)
	}

	fmt.Printf("You (%v) are following these feeds:\n", user.Name)
	for _, feed := range following {
		fmt.Printf("\t* %v - %v\n", feed.FeedName, feed.FeedUrl)
	}

	return nil
}
