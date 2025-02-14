package main

import (
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

	appState := state{
		config: &dbConfig,
	}

	commands := commands{
		commands: make(map[string]func(*state, command, *log.Logger) error),
	}
	commands.register("login", handlerLogin, logger)

	args := os.Args
	if len(args) < 2 {
		logger.Fatal("Usage: cli <command> [arguments...]")
	}

	command := command{
		name: args[1],
		args: args[2:],
	}

	if err := commands.run(&appState, command, logger); err != nil {
		log.Fatal(err)
	}
}
