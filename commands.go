package main

import (
	"fmt"

	"github.com/charmbracelet/log"
)

func (c *commands) register(name string, f func(*state, command, *log.Logger) error, logger *log.Logger) {
	logger.Info("Registering new handler function", "name", name, "func", f)
	c.commands[name] = f
}

func (c *commands) run(s *state, cmd command, logger *log.Logger) error {
	command, exists := c.commands[cmd.name]
	if !exists {
		logger.Error("Command not found", "cmd", cmd)
		return fmt.Errorf("Command not found: %s", cmd.name)
	}

	if err := command(s, cmd, logger); err != nil {
		logger.Error(err)
		return err
	}

	return nil
}
