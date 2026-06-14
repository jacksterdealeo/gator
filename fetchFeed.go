package main

import (
	"bytes"
	"context"
	"database/sql"
	"encoding/xml"
	"html"
	"io"
	"log"
	"net/http"
	"time"

	"github.com/jacksterdealeo/gator/internal/database"

	"github.com/google/uuid"
	"github.com/lib/pq"
)

type RSSFeed struct {
	Channel struct {
		Title       string    `xml:"title"`
		Link        string    `xml:"link"`
		Description string    `xml:"description"`
		Item        []RSSItem `xml:"item"`
	} `xml:"channel"`
}

type RSSItem struct {
	Title       string `xml:"title"`
	Link        string `xml:"link"`
	Description string `xml:"description"`
	PubDate     string `xml:"pubDate"`
}

func fetchFeed(ctx context.Context, feedURL string) (*RSSFeed, error) {
	var client http.Client
	req, err := http.NewRequestWithContext(
		ctx,
		"GET",
		feedURL,
		bytes.NewBuffer([]byte{}),
	)
	if err != nil {
		return nil, err
	}
	req.Header.Set("User-Agent", "gator")

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var feed RSSFeed
	err = xml.Unmarshal(body, &feed)
	if err != nil {
		return nil, err
	}

	feed.Channel.Title = html.UnescapeString(feed.Channel.Title)
	feed.Channel.Description = html.UnescapeString(feed.Channel.Description)

	for i := 0; i < len(feed.Channel.Item); i++ {
		feed.Channel.Item[i].Title = html.UnescapeString(feed.Channel.Item[i].Title)
		feed.Channel.Item[i].Description = html.UnescapeString(feed.Channel.Item[i].Description)
	}

	return &feed, nil
}

/*
This function only scrapes one feed at a time.
This should be ran in a loop.
*/
func scrapeFeeds(s *state, dbUser database.User) error {
	// Get the next feed to fetch from the DB.
	feed, err := s.db.GetNextFeedToFetch(context.Background(), dbUser.ID)
	if err != nil {
		return err
	}
	// Fetch the feed using the URL (we already wrote this function)
	rssFeed, err := fetchFeed(context.Background(), feed.Url)
	if err != nil {
		return err
	}

	for _, post := range rssFeed.Channel.Item {
		if len(post.Title) == 0 {
			continue
		}
		postToCreate := database.CreatePostParams{
			ID:          uuid.New(),
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
			Title:       post.Title,
			Url:         post.Link,
			Description: sql.NullString{String: post.Description, Valid: true},
			PublishedAt: sql.NullTime{Time: time.Now(), Valid: false},
			FeedID:      feed.ID,
		}
		if len(rssFeed.Channel.Description) == 0 {
			postToCreate.Description.Valid = false
		}
		var timeFormats = []string{
			// time.RFC1123Z seems to be most common?
			time.RFC1123Z,
			time.RFC1123,
			time.RFC822Z,
			time.RFC822,
			time.Layout,
			time.ANSIC,
			time.UnixDate,
			time.RubyDate,
			time.RFC850,
			time.RFC3339,
			time.RFC3339Nano,
		}
		for _, timeFormat := range timeFormats {
			t, err := time.Parse(timeFormat, post.PubDate)
			if err != nil {
				continue
			}
			postToCreate.PublishedAt = sql.NullTime{Time: t, Valid: true}
			break
		}
		if !postToCreate.PublishedAt.Valid {
			log.Println("ERR: Time could not be parsed")
		}

		if err := s.db.CreatePost(context.Background(), postToCreate); err != nil {
			pqErr := err.(*pq.Error)
			// 23505 is 'unique_violation' error.
			if pqErr.Code != "23505" {
				log.Println("ERR: pq err ", pqErr)
			}
		}
	}

	// fmt.Println("RSS Channel: ", rssFeed.Channel.Title)
	// // Iterate over the items in the feed and print their titles to the console.
	// for _, item := range rssFeed.Channel.Item {
	// 	fmt.Println("RSS Item: ", item.Title)
	// }

	// Mark it as fetched.
	if err := s.db.MarkFeedFetched(context.Background(), feed.ID); err != nil {
		return err
	}
	return nil
}
