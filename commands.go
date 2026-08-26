package main

import (
	"fmt"
	"gator/internal/config"
	"gator/internal/database"
)

type state struct {
	Db     *database.Queries
	Config *config.Config
}

type command struct {
	Name string
	Args []string
}

type commands struct {
	RegisteredCmds map[string]func(*state, command) error
}

func (c *commands) register(name string, f func(*state, command) error) {
	c.RegisteredCmds[name] = f
}

func (c *commands) run(s *state, cmd command) error {
	handler, exists := c.RegisteredCmds[cmd.Name]
	if !exists {
		return fmt.Errorf("command %q not found", cmd.Name)
	}

	return handler(s, cmd)
}
