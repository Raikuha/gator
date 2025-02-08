package main

import (
	"context"
	"database/sql"
	"encoding/xml"
	"fmt"
	"html"
	"io"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/Raikuha/gator/internal/database"
	"github.com/google/uuid"
)


func HandlerAddFeed(s *state, cmd command, user database.User) error {
	if err := checkArgs(cmd.Args, 2); err != nil {
		return err
	}

	ctx, name, url := context.Background(), cmd.Args[0], cmd.Args[1]

	feed, err := s.db.CreateFeed(ctx, database.CreateFeedParams{
		ID: uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Name: name,
		Url: url,
		UserID: user.ID,
	})
	if err != nil {
		return err
	}

	_, err = s.db.CreateFeedFollow(context.Background(), database.CreateFeedFollowParams{
		ID: uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		UserID: user.ID,
		FeedID: feed.ID,
	})
	if err != nil {
		return err
	}

	return nil
}


func HandlerFeeds(s *state, cmd command) error {
	feeds, err := s.db.GetFeeds(context.Background())
	if err != nil {
		return err
	}

	for _, feed := range feeds {
		fmt.Printf("Title: %s\n", feed.Name)
		fmt.Printf("URL: %s\n", feed.Url)
		fmt.Printf("Added by: %s\n\n", feed.User)
	}
	return nil
}


func scrapeFeeds(s *state) {

	ctx := context.Background()

	next, err := s.db.GetNextFeedToFetch(ctx)
	if err != nil {
		log.Printf("Couldn't fetch a feed: %v", err)
		return
	}

	next, err = s.db.MarkFeedFetched(ctx, next.ID)
	if err != nil {
		log.Printf("Failed to mark feed %s as fetched: %v", next.Name, err)
		return
	}

	rss, err := fetchFeed(ctx, next.Url)
	if err != nil {
		log.Printf("Couldn't fetch the feed %s: %v", next.Name, err)
		return
	}

	for _, item := range rss.Channel.Item {
		pubdate := sql.NullTime{}
		if date, err := time.Parse(time.RFC1123Z, item.PubDate); err == nil {
			pubdate = sql.NullTime{
				Time: date,
				Valid: true,
			}
		}

		_, err = s.db.CreatePost(ctx, database.CreatePostParams{
			ID: uuid.New(),
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
			FeedID: next.ID,
			Title: item.Title,
			Description: sql.NullString{
				String: item.Description,
				Valid: true,
			},
			Url: item.Link,
			PublishedAt: pubdate,
		})
		if err != nil {
			if strings.Contains(err.Error(), "duplicate") {
				continue
			}
			log.Printf("Couldn't create post: %v", err)
		}
	}

	log.Printf("Feed %s collected, %v posts found", next.Name, len(rss.Channel.Item))
}


func fetchFeed(ctx context.Context, feedURL string) (*RSSFeed, error) {
	req, err := http.NewRequestWithContext(ctx, "GET", feedURL, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("User-Agent", "gator")
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	data, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	var rss RSSFeed
	err = xml.Unmarshal(data, &rss)
	if err != nil {
		return nil, err
	}

	rss.Channel.Title = html.UnescapeString(rss.Channel.Title)
	rss.Channel.Description = html.UnescapeString(rss.Channel.Description)

	for i, field := range rss.Channel.Item {
		field.Title = html.UnescapeString(field.Title)
		field.Description = html.UnescapeString(field.Description)
		rss.Channel.Item[i] = field
	}

	return &rss, nil
}