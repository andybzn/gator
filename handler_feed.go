package main

import (
	"context"
	"fmt"

	"github.com/charmbracelet/log"
)

func handlerAgg(s *state, _ command, logger *log.Logger) error {
	feedUrl := "https://wagslane.dev/index.xml"

	feedData, err := fetchFeed(context.Background(), feedUrl, logger)
	if err != nil {
		return err
	}

	fmt.Println(*feedData)

	return nil
}

func handlerAddFeed(s *state, cmd command, logger *log.Logger) error {
	if len(cmd.args) < 2 {
		logger.Errorf("No feed details provided. Usage %s <name> <url>", cmd.name)
		return fmt.Errorf("No feed details provided. Usage %s <name> <url>", cmd.name)
	}

	name := cmd.args[0]
	url := cmd.args[1]

	if err := addFeed(s, name, url, logger); err != nil {
		return err
	}

	return nil
}

func handlerFeeds(s *state, _ command, logger *log.Logger) error {
	feeds, err := s.database.GetFeeds(context.Background())
	if err != nil {
		logger.Error("could not retreive feeds", "err", err)
		return fmt.Errorf("could not retreive feeds: %v", err)
	}

	fmt.Printf("found %d feeds:\n", len(feeds))
	fmt.Printf("Feed Name, Feed URL, User\n")
	for _, feed := range feeds {
		fmt.Printf("%s , %s , %s\n", feed.Name, feed.Url, feed.Username)
	}

	return nil
}
