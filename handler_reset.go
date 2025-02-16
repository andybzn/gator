package main

import (
	"context"

	"github.com/charmbracelet/log"
)

func handlerReset(s *state, _ command, logger *log.Logger) error {
	logger.Info("Database: clearing `users` table...")

	if err := s.database.DeleteUsers(context.Background()); err != nil {
		logger.Error(err)
		return err
	}

	logger.Info("Database: `users` table cleared!")
	return nil
}
