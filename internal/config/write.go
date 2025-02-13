package config

import (
	"encoding/json"
	"os"

	"github.com/charmbracelet/log"
)

func Write(cfg Config, logger *log.Logger) error {
	if logger == nil {
		logger = log.New(os.Stderr)
	}

	configFilePath, err := getConfigFilepath(logger)
	if err != nil {
		logger.Error(err)
		return err
	}
	logger.Debug("config file", "path", configFilePath)

	marshalled, err := json.MarshalIndent(cfg, "", "  ")
	if err != nil {
		logger.Error("Unable to marshal config", "path", configFilePath, "err", err)
		return err
	}

	if err := os.WriteFile(configFilePath, marshalled, 0700); err != nil {
		logger.Error("Error writing config file", "err", err)
		return err
	}

	logger.Info("Database Config written to file!", "path", configFilePath)
	return nil
}
