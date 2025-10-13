package main

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/mattnickolaus/gator/internal/database"
)

func scrapeFeeds(s *state) error {
	nextFeed, err := s.db.GetNextFeedToFetch(context.Background())
	if err != nil {
		return err
	}

	var nullTimeLastFetched sql.NullTime
	nullTimeLastFetched.Time = time.Now()
	nullTimeLastFetched.Valid = true

	markParams := database.MarkFeedFetchedParams{
		UpdatedAt:     time.Now(),
		LastFetchedAt: nullTimeLastFetched,
		ID:            nextFeed.ID,
	}

	err = s.db.MarkFeedFetched(context.Background(), markParams)
	if err != nil {
		return err
	}

	feedData, err := fetchFeed(context.Background(), nextFeed.Url)
	if err != nil {
		return err
	}
	feedData.unescapeFields()

	fmt.Printf("%v Items:\n", feedData.Channel.Title)
	for i, items := range feedData.Channel.Item {
		fmt.Printf("\t%d. %v\n", i+1, items.Title)
	}

	return nil
}

func HandlerAgg(s *state, cmd command) error {
	if numArgs := len(cmd.Args); numArgs != 1 {
		return fmt.Errorf("agg command accepts 1 argument <time_between_requests>:\n\tThe designated frequency gator will pull updates from your followed feeds.\n\tMust be in one of the following formats: '1s'= 1 second, '1m'= 1 minute, '1h'= 1 hour,\n\tand accepts the more complex '1h10m10s' = 1 hour 10 mintes 10 second\n")
	}
	timeBetweenRequests := cmd.Args[0]
	tickTime, err := time.ParseDuration(timeBetweenRequests)
	if err != nil {
		return fmt.Errorf("agg command accepts 1 argument <time_between_requests>:\n\tThe designated frequency gator will pull updates from your followed feeds.\n\tMust be in one of the following formats: '1s'= 1 second, '1m'= 1 minute, '1h'= 1 hour,\n\tand accepts the more complex '1h10m10s' = 1 hour 10 mintes 10 second\n")
	}

	fmt.Printf("Collecting feeds every %v\n", tickTime)
	ticker := time.NewTicker(tickTime)
	for ; ; <-ticker.C {
		scrapeFeeds(s)
	}

	return nil
}
