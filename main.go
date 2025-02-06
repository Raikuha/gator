package main

import (
	"database/sql"
	"fmt"
	"os"

	"github.com/Raikuha/gator/internal/config"
	"github.com/Raikuha/gator/internal/database"
	"github.com/Raikuha/gator/internal/commands"
	_ "github.com/lib/pq"
)

func main () {
	var state commands.State
	state.Cfg = config.Read()

	dbURL := state.Cfg.DB_url 

	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	state.Db = database.New(db)

	cmds := commands.Commands{List:make(map[string]func(*commands.State, commands.Command) error, 3)}
	cmds.Register("login", commands.HandlerLogin)
	cmds.Register("register", commands.HandlerRegister)


	args := os.Args

	if len(args) < 3 {
		fmt.Println("Missing arguments")
		os.Exit(1)
	}

	cmd := commands.Command{
		Name: args[1],
		Args: args[2:],
	}

	err = cmds.Run(&state, cmd)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}