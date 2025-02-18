package main

import (
	"context"
	"encoding/xml"
	"fmt"
	"html"
	"io"
	"net/http"
	"time"

	"github.com/andybzn/gator/internal/database"
	"github.com/charmbracelet/log"
	"github.com/google/uuid"
)

func fetchFeed(ctx context.Context, feedURL string, logger *log.Logger) (*RSSFeed, error) {
	client := http.Client{Timeout: 10 * time.Second}

	req, err := http.NewRequestWithContext(ctx, "GET", feedURL, nil)
	if err != nil {
		logger.Error("error forming request", "err", err)
		return nil, fmt.Errorf("error forming request: %v", err)
	}
	req.Header.Set("User-Agent", "gator")

	res, err := client.Do(req)
	if err != nil {
		logger.Error("error making request", "err", err)
		return nil, fmt.Errorf("error making request: %v", err)
	}
	defer res.Body.Close()

	data, err := io.ReadAll(res.Body)
	if err != nil {
		logger.Error("error reading response data", "err", err)
		return nil, fmt.Errorf("error reading response data: %v", err)
	}

	var feed RSSFeed
	if err := xml.Unmarshal(data, &feed); err != nil {
		logger.Error("error decoding response data", "err", err)
		return nil, fmt.Errorf("error decoding response data: %v", err)
	}

	feed.Channel.Title = html.UnescapeString(feed.Channel.Title)
	feed.Channel.Description = html.UnescapeString(feed.Channel.Description)
	for i, RSSItem := range feed.Channel.Item {
		feed.Channel.Item[i].Title = html.UnescapeString(RSSItem.Title)
		feed.Channel.Item[i].Description = html.UnescapeString(RSSItem.Description)
	}

	return &feed, nil
}

func addFeed(s *state, name, url string, logger *log.Logger) error {
	user, err := s.database.GetUser(context.Background(), s.config.CurrentUserName)
	if err != nil {
		logger.Error(err)
		return err
	}

	now := time.Now().UTC()
	feed, err := s.database.CreateFeed(context.Background(), database.CreateFeedParams{
		ID:        uuid.New(),
		CreatedAt: now,
		UpdatedAt: now,
		Name:      name,
		Url:       url,
		UserID:    user.ID,
	})
	if err != nil {
		logger.Error(err)
		return err
	}

	logger.Info("Feed created!", "feed", feed)
	fmt.Println("Feed created!")
	fmt.Printf("* ID:      %s\n", feed.ID)
	fmt.Printf("* Created: %s\n", feed.CreatedAt)
	fmt.Printf("* Updated: %s\n", feed.UpdatedAt)
	fmt.Printf("* Name:    %s\n", feed.Name)
	fmt.Printf("* Url:     %s\n", feed.Url)
	fmt.Printf("* UserId:  %s\n", feed.UserID)

	return nil
}
