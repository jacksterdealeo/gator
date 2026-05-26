package main

import (
	"context"
	"fmt"
	"gator/internal/database"
	"log"
	"os"
	"strings"
	"time"

	"github.com/google/uuid"
)

func handlerLogin(s *state, cmd command) error {
	if len(cmd.args) < 1 {
		return fmt.Errorf("No args given.")
	}
	if _, err := s.db.GetUser(context.Background(), cmd.args[0]); err != nil {
		return err
	}
	if err := s.config.SetUser(cmd.args[0]); err != nil {
		return err
	}
	fmt.Printf("User %v has been set.\n", cmd.args[0])
	return nil
}

func handlerRegister(s *state, cmd command) error {
	if len(cmd.args) < 1 {
		return fmt.Errorf("No args given.")
	}
	if _, err := s.db.GetUser(context.Background(), cmd.args[0]); err != nil {
		// User does not exist.
	} else {
		return fmt.Errorf("User already exists.")
	}
	u, err := s.db.CreateUser(context.Background(),
		database.CreateUserParams{
			ID:        uuid.New(),
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
			Name:      cmd.args[0],
		},
	)
	if err != nil {
		return err
	}
	fmt.Printf("User %v has been created.\n", cmd.args[0])
	handlerLogin(s, cmd)
	log.Printf("DEBUG ~ User: %v\n", u)
	return nil
}

func handlerGetUsers(s *state, _ command) error {
	users, err := s.db.GetUsers(context.Background())
	if err != nil {
		return err
	}
	var userList strings.Builder
	for _, u := range users {
		fmt.Fprintf(&userList, "* %v", u.Name)
		if u.Name == s.config.CurrentUserName {
			fmt.Fprintf(&userList, " (current)")
		}
		fmt.Fprintf(&userList, "\n")
	}

	fmt.Print(userList.String())
	return nil
}

func handlerReset(s *state, _ command) error {
	if err := s.db.ResetUsers(context.Background()); err != nil {
		return err
	}
	return nil
}

func handlerAggregator(s *state, cmd command) error {
	a, err := fetchFeed(context.Background(), "https://www.wagslane.dev/index.xml")
	if err != nil {
		return err
	}
	fmt.Println(a)
	return nil
}

func handlerAddFeed(s *state, cmd command, dbUser database.User) error {
	if len(cmd.args) < 2 {
		return fmt.Errorf("Not enough args given. (need name and URL)")
	}
	nameOfFeed := cmd.args[0]
	urlOfFeed := cmd.args[1]

	newFeed, err := s.db.CreateFeed(context.Background(),
		database.CreateFeedParams{
			ID:        uuid.New(),
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
			Name:      nameOfFeed,
			Url:       urlOfFeed,
			UserID:    dbUser.ID,
		})

	newFollow, err := s.db.CreateFeedFollow(context.Background(),
		database.CreateFeedFollowParams{
			ID:        uuid.New(),
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
			UserID:    dbUser.ID,
			FeedID:    newFeed.ID,
		})
	if err != nil {
		return err
	}
	fmt.Printf("Feed Name: %v\nUser Name: %v\n", newFollow[0].FeedName, newFollow[0].UserName)

	fmt.Println(newFeed)

	return nil
}

func handlerFeeds(s *state, cmd command) error {
	feeds, err := s.db.GetFeeds(context.Background())
	if err != nil {
		return err
	}

	var feedList strings.Builder
	for _, f := range feeds {
		usr, err := s.db.GetUserByID(context.Background(), f.UserID)
		if err != nil {
			return err
		}
		fmt.Fprintf(&feedList, "* %v\n  %v\n  %v\n",
			f.Name, f.Url, usr.Name)
	}

	fmt.Print(feedList.String())
	return nil
}

func handlerFollow(s *state, cmd command, dbUser database.User) error {
	if len(cmd.args) < 1 {
		return fmt.Errorf("No arguments given. (need URL)")
	}
	followUrl := cmd.args[0]

	feedByURL, err := s.db.GetFeedFromURL(context.Background(), followUrl)
	if err != nil {
		fmt.Fprintln(os.Stderr, "Feed likely wasn't registered with <addfeed>!")
		return err
	}

	newFollow, err := s.db.CreateFeedFollow(context.Background(),
		database.CreateFeedFollowParams{
			ID:        uuid.New(),
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
			UserID:    dbUser.ID,
			FeedID:    feedByURL.ID,
		})
	if err != nil {
		return err
	}
	fmt.Printf("Feed Name: %v\nUser Name: %v\n", newFollow[0].FeedName, newFollow[0].UserName)

	return nil
}

func handlerFollowing(s *state, cmd command, dbUser database.User) error {
	feedsFollowing, err := s.db.GetFeedFollowsForUser(context.Background(),
		dbUser.ID)
	if err != nil {
		return err
	}
	for _, f := range feedsFollowing {
		fmt.Println(f.FeedName)
	}
	return nil
}
