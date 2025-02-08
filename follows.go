package main

import (
	"context"
	"fmt"
	"time"

	"github.com/Raikuha/gator/internal/database"
	"github.com/google/uuid"
)

func HandlerFollowFeed (s *state, cmd command, user database.User) error {
	if err := checkArgs(cmd.Args, 1); err != nil {
		return err
	}

	ctx, url := context.Background(), cmd.Args[0]

	feed_id, err := s.db.GetFeed(ctx, url)
	if err != nil {
		return err
	}

	follow, err := s.db.CreateFeedFollow(ctx, database.CreateFeedFollowParams{
		ID: uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		UserID: user.ID,
		FeedID: feed_id,
	})
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

	s.db.Unfollow(context.Background(), database.UnfollowParams{
		UserID: user.ID,
		FeedID: feed_id,
	})
	return nil
}


func HandlerFollowing(s *state, cmd command, user database.User) error {

	follows, err := s.db.GetFeedFollowsForUser(context.Background(), user.ID)
	if err != nil {
		return err
	}

	fmt.Println("Currently following:")
	for _, feed := range follows {
		fmt.Printf(" * %s\n", feed.Title)
	}
	return nil
}