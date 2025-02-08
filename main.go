package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"os"

	"github.com/Raikuha/gator/internal/config"
	"github.com/Raikuha/gator/internal/database"
	_ "github.com/lib/pq"
)

type state struct {
	cfg *config.Config
	db *database.Queries
}

func main () {
	var s state
	s.cfg = config.Read()

	db, err := sql.Open("postgres", s.cfg.DB_url)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	s.db = database.New(db)

	cmds := commands{
		registered: make(map[string]func(*state, command) error, 9),
	}

	register_commands(cmds)

	args := os.Args

	if len(args) < 2 {
		log.Fatal("usage <command> [args...]")
	}

	cmd := command{
		Name: args[1],
		Args: args[2:],
	}

	err = cmds.Run(&s, cmd)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func register_commands(cmds commands) {
	// User commands
	cmds.register("users", HandlerUsers)
	cmds.register("register", HandlerRegister)
	cmds.register("login", HandlerLogin)
	cmds.register("reset", HandlerReset)

	// Feeds commands
	cmds.register("feeds", HandlerFeeds)
	cmds.register("addfeed", middlewareLoggedIn(HandlerAddFeed))

	// Follow commands
	cmds.register("follow", middlewareLoggedIn(HandlerFollowFeed))
	cmds.register("unfollow", middlewareLoggedIn(HandlerUnfollow))
	cmds.register("following", middlewareLoggedIn(HandlerFollowing))

	// Posts commands
	cmds.register("agg", HandlerAgg)
	cmds.register("browse", middlewareLoggedIn(HandlerBrowse))
}

func middlewareLoggedIn(handler func(s *state, cmd command, user database.User) error) func(*state, command) error {

	return func(s *state, cmd command) error {
		user, err := s.db.GetUser(context.Background(), s.cfg.Current_user_name)
		if err != nil {
			return err
		}

		return handler(s, cmd, user)
	}
}