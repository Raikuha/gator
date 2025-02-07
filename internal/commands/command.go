package commands

import (
	"fmt"

	"github.com/Raikuha/gator/internal/config"
	"github.com/Raikuha/gator/internal/database"
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

func checkArgs(args []string, req int) error {
	if len(args) == 0 {
		return fmt.Errorf("no arguments given")
	}

	if len(args) < req {
		return fmt.Errorf("missing arguments")
	}

	return nil
}
