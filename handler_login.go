package main

import (
	"fmt"
	"os"

	"github.com/charmbracelet/log"
)

func handlerLogin(s *state, cmd command, logger *log.Logger) error {
	if logger == nil {
		logger = log.New(os.Stderr)
	}

	if len(cmd.args) == 0 {
		logger.Errorf("No username provided. Usage: %s <name>", cmd.name)
		return fmt.Errorf("No username provided. Usage: %s <name>", cmd.name)
	}

	if err := s.config.SetUser(cmd.args[0], logger); err != nil {
		return err
	}

	logger.Info("username set!", "name", cmd.args[0])
	fmt.Printf("username set to: %s", cmd.args[0])

	return nil
}
