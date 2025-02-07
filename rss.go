package main

import (
	"context"
	"encoding/xml"
	"fmt"
	"html"
	"io"
	"net/http"
	"time"

	"github.com/Raikuha/gator/internal/database"
	"github.com/google/uuid"
)

type RSSFeed struct {
	Channel struct {
		Title string `xml:"title"`
		Link string `xml:"link"`
		Description string `xml:"description"`
		Item []RSSItem `xml:"item"`
	} `xml:"channel"`
}

type RSSItem struct {
	Title string `xml:"title"`
	Link string `xml:"link"`
	Description string `xml:"description"`
	PubDate string `xml:"item"`
}


func HandlerAgg(s *state, cmd command) error {
	rss, err := fetchFeed(context.Background(), "https://www.wagslane.dev/index.xml")
	if err != nil {
		return err
	}

	fmt.Println(*rss)

	return nil
}


func HandlerAddFeed(s *state, cmd command, user database.User) error {
	if err := checkArgs(cmd.Args, 2); err != nil {
		return err
	}

	ctx, name, url := context.Background(), cmd.Args[0], cmd.Args[1]

	user, err := s.db.GetUser(ctx, s.cfg.Current_user_name)
	if err != nil {
		return err
	}

	feed_data := database.CreateFeedParams{
		ID: uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Name: name,
		Url: url,
		UserID: user.ID,

	}

	feed, err := s.db.CreateFeed(ctx, feed_data)
	if err != nil {
		return err
	}

	follow_data := database.CreateFeedFollowParams{
		ID: uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		UserID: user.ID,
		FeedID: feed.ID,
	}

	_, err = s.db.CreateFeedFollow(context.Background(), follow_data)
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
		fmt.Printf("%s\n%s\n%s\n\n", feed.Name, feed.Url, feed.User)
	}
	return nil
}

func HandlerFollowFeed (s *state, cmd command, user database.User) error {
	if err := checkArgs(cmd.Args, 1); err != nil {
		return err
	}

	user, err:= s.db.GetUser(context.Background(), s.cfg.Current_user_name)
	if err != nil {
		return err
	}

	ctx, url := context.Background(), cmd.Args[0]

	feed_id, err := s.db.GetFeed(ctx, url)
	if err != nil {
		return err
	}

	follow_data := database.CreateFeedFollowParams{
		ID: uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		UserID: user.ID,
		FeedID: feed_id,
	}

	follow, err := s.db.CreateFeedFollow(context.Background(), follow_data)
	if err != nil {
		return err
	}

	fmt.Printf("%s is now following %s\nCreated by %s",
		s.cfg.Current_user_name, follow.Feed, follow.Creator)
	return nil
}

func HandlerUnfollow(s *state, cmd command, user database.User) error {
	if err := checkArgs(cmd.Args, 1); err != nil {
		return err
	}

	url := cmd.Args[0]

	feed_id, err := s.db.GetFeed(context.Background(), url)
	if err != nil {
		return err
	}

	unfollow := database.UnfollowParams{
		UserID: user.ID,
		FeedID: feed_id,
	}

	s.db.Unfollow(context.Background(), unfollow)
	return nil
}

func HandlerFollowing(s *state, cmd command, user database.User) error {
	ctx := context.Background()
	user, err := s.db.GetUser(ctx, s.cfg.Current_user_name)
	if err != nil {
		return err
	}

	follows, err := s.db.GetFeedFollowsForUser(ctx, user.ID)
	if err != nil {
		return err
	}

	fmt.Println("Currently following:")
	for _, feed := range follows {
		fmt.Printf(" * %s\n", feed.Title)
	}
	return nil
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

	data, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	var rss RSSFeed
	if err = xml.Unmarshal(data, &rss); err != nil {
		return nil, err
	}

	rss.Channel.Title = html.UnescapeString(rss.Channel.Title)
	rss.Channel.Description = html.UnescapeString(rss.Channel.Description)

	for _, field := range rss.Channel.Item {
		fmt.Println("TITLE", field.Description)
		field.Title = html.UnescapeString(field.Title)
		fmt.Println("DECODED", field.Description)
		field.Description = html.UnescapeString(field.Description)
	}

	return &rss, nil
}
