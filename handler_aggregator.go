package main

import (
	"context"
	"database/sql"
	"fmt"
	"strconv"
	"time"

	"github.com/google/uuid"
	"github.com/lib/pq"
	"github.com/mattnickolaus/gator/internal/database"
)

var rssDateLayouts = []string{
	time.RFC1123Z,                     // Mon, 02 Jan 2006 15:04:05 -0700
	time.RFC1123,                      // Mon, 02 Jan 2006 15:04:05 MST
	time.RFC3339,                      // 2006-01-02T15:04:05Z07:00
	time.RFC3339Nano,                  // 2006-01-02T15:04:05.999999999Z07:00
	"Mon, 02 Jan 2006 15:04:05 +0000", // Common RFC 822 variant
	"Mon, 02 Jan 2006 15:04:05 GMT",   // The example format from your query
	"2006-01-02T15:04:05-0700",        // ISO 8601 without colon in timezone
	"2006-01-02 15:04:05",             // Common non-standard format without timezone
}

func ParseFlexibleTime(dateTimeStr string) (time.Time, error) {
	for _, layout := range rssDateLayouts {
		t, err := time.Parse(layout, dateTimeStr)
		if err == nil {
			return t, nil
		}
	}
	return time.Time{}, fmt.Errorf("unrecognized date/time format: %s", dateTimeStr)
}

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

	fmt.Printf("Currently Aggregating from: %v\n", feedData.Channel.Title)

	for _, items := range feedData.Channel.Item {
		pubDate, _ := ParseFlexibleTime(items.PubDate)

		postCreateParam := database.CreatePostParams{
			ID:          uuid.New(),
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
			Title:       items.Title,
			Url:         items.Link,
			Description: sql.NullString{String: items.Description, Valid: true},
			PublishedAt: pubDate,
			FeedID:      nextFeed.ID,
		}
		_, err = s.db.CreatePost(context.Background(), postCreateParam)
		// We don't do anything with the returned Post struct
		if err != nil {
			if pqErr, ok := err.(*pq.Error); ok && pqErr.Code == "23505" {
				// do nothing
			} else {
				return fmt.Errorf("Error inserting post: %v", err)
			}
		}
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
		err := scrapeFeeds(s)
		if err != nil {
			return err
		}
	}

	return nil
}

func HandlerBrowse(s *state, cmd command, user database.User) error {
	if numArgs := len(cmd.Args); numArgs > 1 {
		return fmt.Errorf("browse command optionally accepts 1 argument:\n\tlimit: the number of posts that will be displayed (defaults to 2) from your followed feeds\n\tThese posts are sorted by most recent published date.")
	}
	var limit int32 = 2
	if len(cmd.Args) == 1 {
		lim64, err := strconv.ParseInt(cmd.Args[0], 10, 32)
		if err != nil {
			return fmt.Errorf("Error parsing the string parameter to int: %v", err)
		}
		limit = int32(lim64)
	}

	postQueryParam := database.GetPostsParams{
		UserID: user.ID,
		Limit:  limit,
	}

	posts, err := s.db.GetPosts(context.Background(), postQueryParam)
	if err != nil {
		return fmt.Errorf("Error querying posts: %v", err)
	}

	fmt.Printf("Top %v Most recent Posts:\n\n", limit)

	for _, post := range posts {
		fmt.Printf("%v\n", post.Title)
		fmt.Printf("\tUrl: %v\n", post.Url)
		fmt.Printf("\tPublished: %v\n", post.PublishedAt.Format("Jan 02 2006"))
		fmt.Printf("\tDescription: %v\n\n", post.Description.String)
	}

	return nil
}
