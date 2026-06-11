package main

import (
	"bytes"
	"context"
	"encoding/xml"
	"fmt"
	"gator/internal/database"
	"html"
	"io"
	"net/http"
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

	fmt.Println("RSS Channel: ", rssFeed.Channel.Title)
	// Iterate over the items in the feed and print their titles to the console.
	for _, item := range rssFeed.Channel.Item {
		fmt.Println("RSS Item: ", item.Title)
	}
	// Mark it as fetched.
	if err := s.db.MarkFeedFetched(context.Background(), feed.ID); err != nil {
		return err
	}
	return nil
}
