package main

import (
	"context"
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

	username := cmd.args[0]
	_, err := s.database.GetUser(context.Background(), username)

	if err != nil {
		if err.Error() == "sql: no rows in result set" {
			logger.Error("user not found", "user", username)
			return fmt.Errorf("user not found: %s", username)
		} else {
			logger.Error(err)
			return err
		}
	}

	if err := s.config.SetUser(username, logger); err != nil {
		return err
	}

	logger.Info("username set!", "name", username)
	fmt.Printf("username set to: %s", username)

	return nil
}
