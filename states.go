package main

import (
	"context"
	"fmt"
	"gator/internal/config"
	"gator/internal/database"
	"log"
	"time"

	"github.com/google/uuid"
)

type state struct {
	db     *database.Queries
	config *config.Config
}

type command struct {
	name string
	args []string
}

func handlerLogin(s *state, cmd command) error {
	if len(cmd.args) < 1 {
		return fmt.Errorf("No args given.")
	}
	if _, err := s.db.GetUser(context.Background(), cmd.args[0]); err != nil {
		log.Fatalf("User \"%v\" does not exist.\n", cmd.args[0])
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
		log.Fatalln("User already exists.")
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

type commands struct {
	execute map[string]func(*state, command) error
}

// not confident
func (c *commands) run(s *state, cmd command) error {
	if c, ok := c.execute[cmd.name]; ok {
		c(s, cmd)
	} else {
		return fmt.Errorf("Command %v not registered.", cmd.name)
	}
	return nil
}

func (c *commands) register(name string, f func(*state, command) error) {
	c.execute[name] = f
}
