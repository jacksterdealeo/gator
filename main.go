package main

import (
	"database/sql"
	"fmt"
	"gator/internal/config"
	"gator/internal/database"
	"log"
	"os"

	_ "github.com/lib/pq"
)

func main() {
	conf, err := config.Read()
	if err != nil {
		log.Fatalln(err)
	}
	db, err := sql.Open("postgres", conf.DbURL)
	if err != nil {
		log.Fatalln(err)
	}
	dbQueries := database.New(db)
	appState := state{dbQueries, conf}

	appCommands := commands{
		execute: make(map[string]func(*state, command) error),
	}

	appCommands.register("login", handlerLogin)
	appCommands.register("register", handlerRegister)
	appCommands.register("reset", handlerReset)
	appCommands.register("users", handlerGetUsers)
	appCommands.register("agg", handlerAggregator)

	if len(os.Args) < 2 {
		log.Fatalln("No command given.")
	}

	nowCommand := command{
		os.Args[1],
		os.Args[2:],
	}

	if err := appCommands.run(&appState, nowCommand); err != nil {
		log.Fatalln(err)
	}

	fmt.Printf("Gator started.\nconfig: %v\n", appState.config)
}
