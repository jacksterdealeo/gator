package main

import (
	"fmt"
	"gator/internal/config"
)

type state struct {
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
	if err := s.config.SetUser(cmd.args[0]); err != nil {
		return err
	}
	fmt.Printf("User %v has been set.\n", cmd.args[0])
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
