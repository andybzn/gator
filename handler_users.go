package main

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/charmbracelet/log"
	"github.com/google/uuid"

	"github.com/andybzn/gator/internal/database"
)

func handlerRegister(s *state, cmd command, logger *log.Logger) error {
	if len(cmd.args) == 0 {
		logger.Errorf("No username provided. Usage %s <name>", cmd.name)
		return fmt.Errorf("No username provided. Usage %s <name>", cmd.name)
	}

	t := time.Now().UTC()
	user, err := s.database.CreateUser(context.Background(), database.CreateUserParams{
		ID:        uuid.New(),
		CreatedAt: t,
		UpdatedAt: t,
		Name:      cmd.args[0],
	})

	logger.Debug(user)
	logger.Debug(err)

	if err != nil {
		logger.Error(err)
		return err
	}

	if err := s.config.SetUser(cmd.args[0], logger); err != nil {
		return err
	}

	return nil
}

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
