package main

import (
	"context"
	"fmt"
	"gator/internal/database"
	"log"
	"time"

	"github.com/google/uuid"
)

func handlerLogin(s *state, cmd command) error {
	if len(cmd.args) < 1 {
		return fmt.Errorf("No args given.")
	}
	if _, err := s.db.GetUser(context.Background(), cmd.args[0]); err != nil {
		// log.Fatalf("User \"%v\" does not exist.\n", cmd.args[0])
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
		// log.Fatalln("User already exists.")
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

func handlerReset(s *state, _ command) error {
	if err := s.db.ResetUsers(context.Background()); err != nil {
		return err
	}
	return nil
}
