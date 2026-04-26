package main

import (
	"fmt"
	"gator/internal/config"
	"gator/internal/database"
)

type state struct {
	db     *database.Queries
	config *config.Config
}

type command struct {
	name string
	args []string
}

type commands struct {
	execute map[string]func(*state, command) error
}

func (c *commands) run(s *state, cmd command) error {
	if c, ok := c.execute[cmd.name]; ok {
		if err := c(s, cmd); err != nil {
			// log.Fatalln(err)
			return err
		}
	} else {
		return fmt.Errorf("Command %v not registered.", cmd.name)
	}
	return nil
}

func (c *commands) register(name string, f func(*state, command) error) {
	c.execute[name] = f
}
