package main

import (
	"database/sql"
	"os"

	"github.com/andybzn/gator/internal/config"
	"github.com/andybzn/gator/internal/database"
	"github.com/charmbracelet/log"

	_ "github.com/lib/pq"
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

	db, err := sql.Open("postgres", dbConfig.DbUrl)
	if err != nil {
		logger.Fatal(err)
	}
	dbQueries := database.New(db)

	appState := state{
		database: dbQueries,
		config:   &dbConfig,
	}

	commands := commands{
		commands: make(map[string]func(*state, command, *log.Logger) error),
	}
	commands.register("login", handlerLogin, logger)
	commands.register("register", handlerRegister, logger)
	commands.register("reset", handlerReset, logger)

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
