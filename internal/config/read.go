package config

import (
	"encoding/json"
	"os"

	"github.com/charmbracelet/log"
)

func Read(logger *log.Logger) (Config, error) {
	if logger == nil {
		logger = log.New(os.Stderr)
	}
	logger.Info("Reading config file...")

	configFilePath, err := getConfigFilepath(logger)
	if err != nil {
		logger.Error(err)
		return Config{}, err
	}
	logger.Debug("config file", "path", configFilePath)

	configData, err := os.ReadFile(configFilePath)
	if err != nil {
		logger.Error("Error reading config file", "err", err)
		return Config{}, err
	}

	var config Config
	if err := json.Unmarshal(configData, &config); err != nil {
		logger.Error("Unable to parse config", "path", configFilePath, "err", err)
		return Config{}, err
	}

	logger.Info("Finished reading config file")
	logger.Debug(config)

	return config, nil
}
