package config

import (
	"fmt"
	"os"

	"github.com/charmbracelet/log"
)

func getConfigFilepath(logger *log.Logger) (string, error) {
	if logger == nil {
		logger = log.New(os.Stderr)
	}

	homeDir, err := os.UserHomeDir()
	if err != nil {
		logger.Error("unable to access $HOME", "err", err)
		return "", err
	}

	return fmt.Sprintf("%s/%s", homeDir, configFileName), nil
}
