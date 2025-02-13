package main

import (
	"fmt"
	"os"

	"github.com/andybzn/gator/internal/config"
	"github.com/charmbracelet/log"
)

func main() {
	logger := log.NewWithOptions(os.Stderr, log.Options{
		ReportCaller: true,
		Level:        log.ErrorLevel,
	})
	if debugLevel := os.Getenv("LOG_LEVEL"); debugLevel == "DEBUG" {
		logger.SetLevel(log.DebugLevel)
	}

	dbConfig, err := config.Read(logger)
	if err != nil {
		logger.Fatal(err)
	}

	if err := dbConfig.SetUser("andy", logger); err != nil {
		logger.Fatal(err)
	}

	dbConfig, err = config.Read(logger)
	if err != nil {
		logger.Fatal(err)
	}
	fmt.Printf("%s\n", dbConfig.DbUrl)
	fmt.Printf("%s\n", dbConfig.CurrentUserName)
}
