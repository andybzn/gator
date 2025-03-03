package main

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/andybzn/gator/internal/database"
	"github.com/charmbracelet/log"
	"github.com/google/uuid"
)

func handlerAgg(s *state, cmd command, logger *log.Logger) error {
	if len(cmd.args) < 1 {
		logger.Errorf("No time duration provided. Usage %s <duration(s|m|h)>", cmd.name)
		return fmt.Errorf("No time duration provided. Usage %s <duration(s|m|h)>", cmd.name)
	}

	timeBetweenRequests, err := time.ParseDuration(cmd.args[0])
	if err != nil {
		return err
	}

	fmt.Printf("Collecting feeds every %v", timeBetweenRequests)

	ticker := time.NewTicker(timeBetweenRequests)
	for ; ; <-ticker.C {
		scrapeFeeds(s, logger)
	}
}

func handlerAddFeed(s *state, cmd command, user database.User, logger *log.Logger) error {
	if len(cmd.args) < 2 {
		logger.Errorf("No feed details provided. Usage %s <name> <url>", cmd.name)
		return fmt.Errorf("No feed details provided. Usage %s <name> <url>", cmd.name)
	}

	name := cmd.args[0]
	url := cmd.args[1]

	if err := addFeed(s, name, url, logger); err != nil {
		return err
	}

	cmd.args = cmd.args[1:]

	if err := handlerFollow(s, cmd, user, logger); err != nil {
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

func handlerFollow(s *state, cmd command, user database.User, logger *log.Logger) error {
	if len(cmd.args) == 0 {
		logger.Errorf("No feed url provided. Usage %s <url>", cmd.name)
		return fmt.Errorf("No feed url provided. Usage %s <url>", cmd.name)
	}
	feedUrl := cmd.args[0]

	ctx := context.Background()

	feed, err := s.database.GetFeedByUrl(ctx, feedUrl)
	if err != nil {
		logger.Error("Feed does not exist or could not be found", "feedUrl", feedUrl, "err", err)
		return fmt.Errorf("Feed %s does not exist or could not be found: %v", feedUrl, err)
	}

	now := time.Now().UTC()
	followed, err := s.database.CreateFeedFollow(ctx, database.CreateFeedFollowParams{
		ID:        uuid.New(),
		CreatedAt: now,
		UpdatedAt: now,
		UserID:    user.ID,
		FeedID:    feed.ID,
	})
	if err != nil {
		logger.Error("Feed could not be followed", "err", err)
		return fmt.Errorf("Feed could not be followed: %v", err)
	}

	logger.Info("Feed followed!", "followedFeed", followed)
	fmt.Println("Feed followed!")
	fmt.Printf("* ID:       %s\n", followed.ID)
	fmt.Printf("* Created:  %s\n", followed.CreatedAt)
	fmt.Printf("* Updated:  %s\n", followed.UpdatedAt)
	fmt.Printf("* UserId:   %s\n", followed.UserID)
	fmt.Printf("* Username: %s\n", followed.Username)
	fmt.Printf("* FeedId:   %s\n", followed.FeedID)
	fmt.Printf("* Feedname: %s\n", followed.Feedname)

	return nil
}

func handlerFollowing(s *state, _ command, user database.User, logger *log.Logger) error {
	ctx := context.Background()

	feeds, err := s.database.GetFeedFollowsForUser(ctx, user.ID)
	if err != nil {
		logger.Error("Error retrieving feed follows", "err", err)
		return fmt.Errorf("Error retrieving feed follows: %v", err)
	}

	if len(feeds) == 0 {
		return nil
	}

	for _, feed := range feeds {
		fmt.Println(feed.Name)
	}

	return nil
}

func handlerUnfollow(s *state, cmd command, user database.User, logger *log.Logger) error {
	if len(cmd.args) == 0 {
		logger.Errorf("No feed url provided. Usage %s <url>", cmd.name)
		return fmt.Errorf("No feed url provided. Usage %s <url>", cmd.name)
	}

	ctx := context.Background()
	feed, err := s.database.GetFeedByUrl(ctx, cmd.args[0])
	if err != nil {
		logger.Error("Could not fetch feed data", "err", err)
		return fmt.Errorf("Could not fetch feed data: %v", err)
	}

	if err := s.database.RemoveFeedFollow(ctx, database.RemoveFeedFollowParams{
		FeedID: feed.ID,
		UserID: user.ID,
	}); err != nil {
		logger.Error("Could not unfollow feed", "err", err)
		return fmt.Errorf("Could not unfollow feed: %v", err)
	}

	return nil
}

func handlerBrowse(s *state, cmd command, user database.User, logger *log.Logger) error {
	limit := 2
	if len(cmd.args) > 0 {
		if parsedLimit, err := strconv.Atoi(cmd.args[0]); err == nil {
			limit = parsedLimit
		} else {
			logger.Error(err)
			return fmt.Errorf("Unable to parse limit: %v", err)
		}
	}

	posts, err := s.database.GetPostsForUser(context.Background(), database.GetPostsForUserParams{
		UserID: user.ID,
		Limit:  int32(limit),
	})
	if err != nil {
		logger.Error(err)
		return fmt.Errorf("Unable to get posts: %v", err)
	}

	if len(posts) == 0 {
		return nil
	}

	for _, post := range posts {
		fmt.Printf("%s (%s)\n", post.FeedName, post.PublishedAt.Format("Mon 2 Jan"))
		fmt.Printf("~~~ %s ~~~\n", post.Title)
		fmt.Printf("    %s\n", post.Description)
		fmt.Printf("Link: %s\n", post.Url)
		fmt.Println("")
	}

	return nil
}
