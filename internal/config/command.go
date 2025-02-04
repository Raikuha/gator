package config

import (
	"fmt"
)

type State struct {
	Cfg *Config
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
		fun(s, cmd)
		return nil
	}

	return fmt.Errorf("invalid command")
}


func HandlerLogin(s *State, cmd Command) error {
	if len(cmd.Args) == 0 {
		return fmt.Errorf("no arguments given")
	}

	s.Cfg.SetUser(cmd.Args[0])
	fmt.Println("New user has been set")
	return nil
}