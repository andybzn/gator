package main

import (
	"context"
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
	commands.register("users", handlerUsers, logger)
	commands.register("agg", handlerAgg, logger)
	commands.register("addfeed", middlewareLoggedIn(handlerAddFeed), logger)
	commands.register("feeds", handlerFeeds, logger)
	commands.register("follow", middlewareLoggedIn(handlerFollow), logger)
	commands.register("following", middlewareLoggedIn(handlerFollowing), logger)
	commands.register("unfollow", middlewareLoggedIn(handlerUnfollow), logger)

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

func middlewareLoggedIn(handler func(s *state, cmd command, user database.User, logger *log.Logger) error) func(*state, command, *log.Logger) error {
	return func(s *state, cmd command, logger *log.Logger) error {
		user, err := s.database.GetUser(context.Background(), s.config.CurrentUserName)
		if err != nil {
			return err
		}

		return handler(s, cmd, user, logger)
	}
}
