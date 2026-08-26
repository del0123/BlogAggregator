package main

import (
	"context"
	"database/sql"
	"encoding/xml"
	"fmt"
	"gator/internal/database"
	"html"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/google/uuid"
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

// this variable is for scrapeFeeds()'s parsing of different time formats
// performance: outside of function: won't have to allocate memory for the string slice every time scrapeFeeds runs
var timeLayouts = []string{
	time.RFC1123,
	time.RFC1123Z,
	time.RFC3339,
	time.RFC822,
}

func fetchFeed(ctx context.Context, feedURL string) (*RSSFeed, error) {

	req, err := http.NewRequestWithContext(ctx, "GET", feedURL, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("User-Agent", "gator")

	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	rawData, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	var feed RSSFeed
	err = xml.Unmarshal(rawData, &feed)
	if err != nil {
		return nil, err
	}

	feed.Channel.Description = html.UnescapeString(feed.Channel.Description)
	feed.Channel.Title = html.UnescapeString(feed.Channel.Title)

	for i := range feed.Channel.Item {
		item := &feed.Channel.Item[i] //remember that when you create a new variable and want to change an original variable we use pointers, or else we are changing a copy...
		item.Title = html.UnescapeString(item.Title)
		item.Description = html.UnescapeString(item.Description)
	}

	return &feed, nil
}

func scrapeFeeds(s *state) error {
	ctx := context.Background()

	feedRecord, err := s.Db.GetNextFeedToFetch(ctx)
	if err != nil {
		return err
	}

	err = s.Db.MarkFeedFetched(ctx, feedRecord.ID)
	if err != nil {
		return err
	}

	feed, err := fetchFeed(ctx, feedRecord.Url)
	if err != nil {
		return err
	}

	for _, item := range feed.Channel.Item {
		ctx := context.Background()

		var pubTime sql.NullTime
		for _, layout := range timeLayouts {
			if t, err := time.Parse(layout, item.PubDate); err == nil {
				pubTime = sql.NullTime{Time: t, Valid: true}
				break
			}
		}

		description := sql.NullString{
			String: item.Description,
			Valid:  item.Description != "",
		}

		post, err := s.Db.CreatePost(ctx, database.CreatePostParams{
			ID:          uuid.New(),
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
			Title:       item.Title,
			Url:         item.Link,
			Description: description,
			PublishedAt: pubTime,
			FeedID:      feedRecord.ID,
		})
		if err != nil {
			if strings.Contains(err.Error(), "duplicate key") || strings.Contains(err.Error(), "unique constraint") {
				continue
			}
			fmt.Printf("Couldn't create post: %v\n", err)
		}
		fmt.Printf("Saved post %s with ID %s\n", post.Title, post.ID)
	}

	fmt.Printf("Feed %s fetched, %d posts processed\n", feedRecord.Name, len(feed.Channel.Item))

	return nil
}
