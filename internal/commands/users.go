package commands

import (
	"context"
	"fmt"
	"time"

	"github.com/Raikuha/gator/internal/database"
	"github.com/google/uuid"
)

func HandlerLogin(s *State, cmd Command) error {
	if err := checkArgs(cmd.Args, 1); err != nil {
		return err
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
	if err := checkArgs(cmd.Args, 1); err != nil {
		return err
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

func HandlerUsers(s *State, cmd Command) error {
	userlist, err := s.Db.GetUsers(context.Background())
	if err != nil {
		return err
	}

	if len(userlist) == 0 {
		fmt.Println("No users registered yet")
	}

	for _, name := range userlist {
		if name == s.Cfg.Current_user_name {
			name += " (current)"
		}

		fmt.Printf("* %s\n", name)
	}

	return nil
}

func HandlerReset(s *State, cmd Command) error {
	err := s.Db.Reset(context.Background())
	if err != nil {
		return err
	}

	fmt.Println("Users table purged correctly")
	return nil
}