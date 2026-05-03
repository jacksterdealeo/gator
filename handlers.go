package main

import (
	"context"
	"fmt"
	"gator/internal/database"
	"log"
	"strings"
	"time"

	"github.com/google/uuid"
)

func handlerLogin(s *state, cmd command) error {
	if len(cmd.args) < 1 {
		return fmt.Errorf("No args given.")
	}
	if _, err := s.db.GetUser(context.Background(), cmd.args[0]); err != nil {
		return err
	}
	if err := s.config.SetUser(cmd.args[0]); err != nil {
		return err
	}
	fmt.Printf("User %v has been set.\n", cmd.args[0])
	return nil
}

func handlerRegister(s *state, cmd command) error {
	if len(cmd.args) < 1 {
		return fmt.Errorf("No args given.")
	}
	if _, err := s.db.GetUser(context.Background(), cmd.args[0]); err != nil {
		// User does not exist.
	} else {
		return fmt.Errorf("User already exists.")
	}
	u, err := s.db.CreateUser(context.Background(),
		database.CreateUserParams{
			ID:        uuid.New(),
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
			Name:      cmd.args[0],
		},
	)
	if err != nil {
		return err
	}
	fmt.Printf("User %v has been created.\n", cmd.args[0])
	handlerLogin(s, cmd)
	log.Printf("DEBUG ~ User: %v\n", u)
	return nil
}

func handlerGetUsers(s *state, _ command) error {
	users, err := s.db.GetUsers(context.Background())
	if err != nil {
		return err
	}
	var userList strings.Builder
	for _, u := range users {
		fmt.Fprintf(&userList, "* %v", u.Name)
		if u.Name == s.config.CurrentUserName {
			fmt.Fprintf(&userList, " (current)")
		}
		fmt.Fprintf(&userList, "\n")
	}

	fmt.Print(userList.String())
	return nil
}

func handlerReset(s *state, _ command) error {
	if err := s.db.ResetUsers(context.Background()); err != nil {
		return err
	}
	return nil
}

func handlerAggregator(s *state, cmd command) error {
	a, err := fetchFeed(context.Background(), "https://www.wagslane.dev/index.xml")
	if err != nil {
		return err
	}
	fmt.Println(a)
	return nil
}
