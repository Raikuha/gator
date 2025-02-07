package main

import (
	"fmt"
)

type command struct {
	Name string
	Args []string
}

type commands struct {
	registered map[string]func(*state, command) error
}

func (c *commands) register(name string, f func(*state, command) error) {
	if _, ok := c.registered[name]; !ok {
		c.registered[name] = f
	}
}

func (c *commands) Run(s *state, cmd command) error {

	if fun, ok := c.registered[cmd.Name]; ok {
		return fun(s, cmd)
	}

	return fmt.Errorf("invalid command")
}

func checkArgs(args []string, req int) error {
	if len(args) == 0 {
		return fmt.Errorf("no arguments given")
	}

	if len(args) < req {
		return fmt.Errorf("missing arguments")
	}

	return nil
}
