package commands

import (
	"context"
	"fmt"
	"time"

	"github.com/Raikuha/gator/internal/config"
	"github.com/Raikuha/gator/internal/database"
	"github.com/google/uuid"
)

type State struct {
	Cfg *config.Config
	Db *database.Queries
}

type Command struct {
	Name string
	Args []string
}

type Commands struct {
	List map[string]func(*State, Command) error
}

func (c *Commands) Register(name string, f func(*State, Command) error) {
	if _, ok := c.List[name]; !ok {
		c.List[name] = f
	}
}


func (c *Commands) Run(s *State, cmd Command) error {

	if fun, ok := c.List[cmd.Name]; ok {
		return fun(s, cmd)
	}

	return fmt.Errorf("invalid command")
}


func HandlerLogin(s *State, cmd Command) error {
	if len(cmd.Args) == 0 {
		return fmt.Errorf("no arguments given")
	}

	name := cmd.Args[0]
	_, err := s.Db.GetUser(context.Background(), name)
	if err != nil {
		return fmt.Errorf("user does not exist")
	}
	s.Cfg.SetUser(name)
	fmt.Println("New active user has been set")
	return nil
}

func HandlerRegister(s *State, cmd Command) error {
	if len(cmd.Args) == 0 {
		return fmt.Errorf("no arguments given")
	}

	user := database.CreateUserParams{
		ID:uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Name: cmd.Args[0]}

	_, err := s.Db.GetUser(context.Background(), user.Name)
	if err == nil {
	 	return fmt.Errorf("user already exists")
	}

	new_user, err := s.Db.CreateUser(context.Background(), user)
	if err != nil {
		return err
	}

	fmt.Printf("User %s created\n", new_user.Name)
	HandlerLogin(s, cmd)
	return nil
}