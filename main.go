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

	appCommands.register(
		"login", func(appState *state, cmd command) error {
			if len(cmd.args) < 1 {
				log.Fatalln("Not enough arguments for login.")
			}
			return handlerLogin(appState, cmd)
		},
	)

	appCommands.register(
		"register", func(appState *state, cmd command) error {
			if len(cmd.args) < 1 {
				log.Fatalln("Not enough arguments for login.")
			}
			return handlerRegister(appState, cmd)
		},
	)

	if len(os.Args) < 2 {
		log.Fatalln("No command given.")
	}

	nowCommand := command{
		os.Args[1],
		os.Args[2:],
	}

	appCommands.run(&appState, nowCommand)

	fmt.Println("Gator started.\n", appState.config)
}
