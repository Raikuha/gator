package main

import (
	"github.com/Raikuha/gator/internal/config"
	"fmt"
)

func main () {
	cfg := config.Read()
	cfg.SetUser("raikuha")

	cfg = config.Read()
	fmt.Println(cfg.DB_url)
	fmt.Println(cfg.Current_user_name)
}