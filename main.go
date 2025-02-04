package main

import (
	"github.com/Raikuha/gator/internal/config"
	"fmt"
	"os"
)

func main () {
	var state config.State
	state.Cfg = config.Read()

	commands := config.Commands{List:make(map[string]func(*config.State, config.Command) error, 3)}
	commands.Register("login", config.HandlerLogin)

	args := os.Args

	if len(args) < 3 {
		fmt.Println("Missing arguments")
		os.Exit(1)
	}

	cmd := config.Command{
		Name: args[1],
		Args: args[2:],
	}

	err := commands.Run(&state, cmd)
	if err != nil {
		fmt.Println(err)
	}
}