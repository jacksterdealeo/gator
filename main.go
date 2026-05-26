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

	appCmds := commands{
		execute: make(map[string]func(*state, command) error),
	}
	appCmds.register("login", handlerLogin)
	appCmds.register("register", handlerRegister)
	appCmds.register("reset", handlerReset)
	appCmds.register("users", handlerGetUsers)
	appCmds.register("agg", handlerAggregator)
	appCmds.register("feeds", handlerFeeds)
	appCmds.register("addfeed", middlewareLoggedIn(handlerAddFeed))
	appCmds.register("follow", middlewareLoggedIn(handlerFollow))
	appCmds.register("following", middlewareLoggedIn(handlerFollowing))

	helpCmd := func(s *state, cmd command) error {
		if len(cmd.args) >= 1 {
			fmt.Printf(" ~ %v : %v\n", cmd.args[0], helpDocs[cmd.args[0]])
			return nil
		}
		fmt.Println("Here are all the avalible commands:")
		for command, _ := range appCmds.execute {
			fmt.Printf(" ~ %v: %v\n", command, helpDocs[command])
		}
		return nil
	}
	appCmds.register("help", helpCmd)

	if len(os.Args) < 2 {
		log.Fatalln("No command given.")
	}

	nowCommand := command{
		os.Args[1],
		os.Args[2:],
	}

	if err := appCmds.run(&appState, nowCommand); err != nil {
		log.Fatalln(err)
	}

	// fmt.Fprintf(os.Stderr, "Gator started.\nconfig: %v\n", appState.config)
}
