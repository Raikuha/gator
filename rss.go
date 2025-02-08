package main

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/Raikuha/gator/internal/database"
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
	if err := checkArgs(cmd.Args, 1); err != nil {
		return err
	}

	interval, err := time.ParseDuration(cmd.Args[0])
	if err != nil {
		return err
	}

	fmt.Printf("Collecting feeds every %v\n", interval)
	ticker := time.Tick(interval)

	for ; ; <-ticker {
		scrapeFeeds(s)
	}
}


func HandlerBrowse (s *state, cmd command, user database.User) error {
	limit := 2

	if len(cmd.Args) == 1 {
		if setLimit, err := strconv.Atoi(cmd.Args[0]); err == nil {
			limit = setLimit
		} else {
			return fmt.Errorf("invalid limit: %w", err)
		}
	}

	posts, err := s.db.GetPostsForUser(context.Background(), database.GetPostsForUserParams{
		UserID: user.ID,
		Limit: int32(limit),
	})
	if err != nil {
		return fmt.Errorf("couldn't get posts: %w", err)
	}

	fmt.Printf("Found %d posts for user %s\n", len(posts), user.Name)
	for _, post := range posts {
		fmt.Printf("%s from %s\n", post.PublishedAt.Time.Format("Mon Jan 2"), post.Feedname)
		fmt.Printf("-- %s --\n", post.Title)
		fmt.Printf("   %v\n", post.Description.String)
		fmt.Printf("Link: %s\n", post.Url)
		fmt.Println("======================================================")
	}
	return nil
}