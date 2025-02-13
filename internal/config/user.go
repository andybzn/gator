package config

import (
	"fmt"
	"os"

	"github.com/charmbracelet/log"
)

func (cfg *Config) SetUser(name string, logger *log.Logger) error {
	if logger == nil {
		logger = log.New(os.Stderr)
	}

	if name == "" {
		logger.Error("username not provided")
		return fmt.Errorf("username not provided")
	}

	logger.Info("Setting username...")
	logger.Debug("username value", "name", name)
	cfg.CurrentUserName = name

	if cfg.CurrentUserName != name {
		logger.Error("Error setting username")
		return fmt.Errorf("Error setting username")
	}

	if err := Write(*cfg, logger); err != nil {
		logger.Error("Failed to write database config", "err", err)
	}

	return nil
}
