package main

import (
	"context"
	"fmt"
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
