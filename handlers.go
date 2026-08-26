package main

import (
	"context"
	"fmt"
	"gator/internal/database"
	"strconv"
	"time"

	"github.com/google/uuid"
)

func handlerLogin(s *state, cmd command) error {
	if len(cmd.Args) != 1 {
		return fmt.Errorf("invalid login arguments. expected: <username>")
	}
	ctx := context.Background()

	user, err := s.Db.GetUser(ctx, cmd.Args[0])
	if err != nil {
		return err
	}

	err = s.Config.SetUser(user.Name)
	if err != nil {
		return err
	}

	fmt.Printf("Username %s login successful.\n", user.Name)

	return nil
}

func handlerRegister(s *state, cmd command) error {
	if len(cmd.Args) != 1 {
		return fmt.Errorf("invalid register arguments. expected: <username>")
	}

	ctx := context.Background()

	user, err := s.Db.CreateUser(ctx, database.CreateUserParams{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Name:      cmd.Args[0]})
	if err != nil {
		return err
	}

	err = s.Config.SetUser(cmd.Args[0])
	if err != nil {
		return err
	}

	fmt.Printf("Username %s\n registered successfully.\n", user.Name)

	return nil
}

func handlerReset(s *state, cmd command) error {
	if len(cmd.Args) >= 1 {
		return fmt.Errorf("too many arguments")
	}

	ctx := context.Background()

	err := s.Db.ResetUsers(ctx)
	if err != nil {
		return err
	}
	fmt.Println("All users have been reset successfully.")

	return nil
}

func handlerGetUsers(s *state, cmd command) error {
	if len(cmd.Args) >= 1 {
		return fmt.Errorf("too many arguments")
	}

	ctx := context.Background()

	userList, err := s.Db.GetUsers(ctx)
	if err != nil {
		return err
	}
	for _, user := range userList {
		if user.Name == s.Config.CurrentUserName {
			fmt.Printf("* %s (current)\n", user.Name)
		} else {
			fmt.Printf("* %s\n", user.Name)
		}
	}
	return nil

}

func handlerAggregate(s *state, cmd command) error {
	if len(cmd.Args) != 1 {
		return fmt.Errorf("usage: agg <time interval>")
	}

	time_between_reqs, err := time.ParseDuration(cmd.Args[0])
	if err != nil {
		return err
	}

	fmt.Printf("Collecting feeds every %s\n", time_between_reqs)

	ticker := time.NewTicker(time_between_reqs)
	defer ticker.Stop()

	for ; ; <-ticker.C {
		err := scrapeFeeds(s)
		if err != nil {
			fmt.Printf("Error scraping feed: %v\n", err)
		}
	}

}

func handlerAddFeed(s *state, cmd command, user database.User) error {
	if len(cmd.Args) != 2 {
		return fmt.Errorf("usage: addfeed <name> <url>")
	}

	ctx := context.Background()

	feed, err := s.Db.CreateFeed(ctx, database.CreateFeedParams{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Name:      cmd.Args[0],
		Url:       cmd.Args[1],
		UserID:    user.ID,
	})
	if err != nil {
		return err
	}

	feedFollow, err := s.Db.CreateFeedFollow(ctx, database.CreateFeedFollowParams{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		UserID:    user.ID,
		FeedID:    feed.ID, // Assuming URL is the feed ID for simplicity

	})
	if err != nil {
		return err
	}

	fmt.Printf("Feed created successfully: %s\n", feed.Name)
	fmt.Printf("User %s is now following feed '%s'\n", feedFollow.UserName, feedFollow.FeedName)

	return nil
}

func handlerGetFeeds(s *state, cmd command) error {
	if len(cmd.Args) != 0 {
		return fmt.Errorf("usage: feeds")
	}
	ctx := context.Background()

	feeds, err := s.Db.GetFeeds(ctx)
	if err != nil {
		return err
	}
	fmt.Printf("%+v\n", feeds)

	return nil
}

func handlerFollow(s *state, cmd command, user database.User) error {
	if len(cmd.Args) != 1 {
		return fmt.Errorf("usage: follow <url>")
	}
	ctx := context.Background()
	url := cmd.Args[0]

	feed, err := s.Db.GetFeedsByURL(ctx, url)
	if err != nil {
		return err
	}

	feedFollow, err := s.Db.CreateFeedFollow(ctx, database.CreateFeedFollowParams{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		UserID:    user.ID,
		FeedID:    feed.ID, // Assuming URL is the feed ID for simplicity

	})
	if err != nil {
		return err
	}

	fmt.Printf("User %s is now following feed '%s'\n", feedFollow.UserName, feedFollow.FeedName)

	return nil
}

func handlerFollowing(s *state, cmd command, user database.User) error {
	if len(cmd.Args) != 0 {
		return fmt.Errorf("usage: following")
	}
	ctx := context.Background()

	feedsByUser, err := s.Db.GetFeedFollowsForUser(ctx, user.ID)
	if err != nil {
		return err
	}
	for _, feedFollow := range feedsByUser {
		fmt.Printf("Following feed '%s'\n", feedFollow.FeedName)
	}

	return nil
}

func handlerUnfollow(s *state, cmd command, user database.User) error {
	if len(cmd.Args) != 1 {
		return fmt.Errorf("usage: unfollow <url>")
	}
	ctx := context.Background()

	feed, err := s.Db.GetFeedsByURL(ctx, cmd.Args[0])
	if err != nil {
		return err
	}

	err = s.Db.DeleteFeedFollow(ctx, database.DeleteFeedFollowParams{
		UserID: user.ID,
		FeedID: feed.ID,
	})
	if err != nil {
		return err
	}

	fmt.Printf("%s unfollowed feed: %s\n", user.Name, feed.Name)

	return nil
}

func handlerBrowse(s *state, cmd command, user database.User) error {
	ctx := context.Background()
	limit := 2

	if len(cmd.Args) == 1 {
		parsedLimit, err := strconv.Atoi(cmd.Args[0])
		if err != nil {
			return fmt.Errorf("usage: browse [limit]")
		}
		limit = parsedLimit
	}

	posts, err := s.Db.GetPostsForUser(ctx, database.GetPostsForUserParams{
		UserID: user.ID,
		Limit:  int32(limit),
	})
	if err != nil {
		return fmt.Errorf("couldn't fetch posts: %w", err)
	}

	fmt.Printf("Found %d posts for user %s:\n", len(posts), user.Name)

	for _, post := range posts {
		fmt.Printf("\n--- %s ---\n", post.Title)
		if post.PublishedAt.Valid {
			fmt.Printf("Published: %s\n", post.PublishedAt.Time)
		}
		fmt.Printf("Link: %s\n", post.Url)
		if post.Description.Valid {
			fmt.Printf("Description: %s\n", post.Description.String)
		}
	}

	return nil
}
