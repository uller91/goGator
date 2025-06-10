package main

import (
	"github.com/uller91/goGator/internal/config"
	"github.com/uller91/goGator/internal/database"
	"errors"
	"fmt"
	"context"
	"github.com/google/uuid"
	"time"
	"os"
	"github.com/lib/pq"
	"strconv"
)

type state struct {
	config *config.Config
	database  *database.Queries
}

type command struct {
	name        string
	arguments	[]string
}


type commands struct {
	handlers	map[string]func(*state, command) error
}

func (c *commands) run(s *state, cmd command) error {
	hndl, exists := c.handlers[cmd.name] 
	if exists {
		err := hndl(s, cmd)
		if err != nil {
			return err
		}
	} else {
		return errors.New("No command with this name is registered")
	}
	
	return nil
}

func (c *commands) register(name string, f func(*state, command) error) {
	c.handlers[name] = f
}


func handlerLogin(s *state, cmd command) error {
	if len(cmd.arguments) != 1 {
		return errors.New("1 argument is expected")
	}

	_, err := s.database.GetUser(context.Background(), cmd.arguments[0])
	if err!= nil {
			fmt.Println("User with this name do not exist!")
			os.Exit(1)
	} 

	err = s.config.SetUser(cmd.arguments[0])
	if err != nil {
		return err
	}

	fmt.Printf("%v is set!\n", cmd.arguments[0])
	return nil
}

func handlerRegister(s *state, cmd command) error {
	if len(cmd.arguments) != 1 {
		return errors.New("1 argument is expected")
	}

	param := database.CreateUserParams{ID: uuid.New(), CreatedAt: time.Now(), UpdatedAt: time.Now(), Name: cmd.arguments[0]}
	user, err := s.database.CreateUser(context.Background(), param)
	if err!= nil {
		if pqError, ok := err.(*pq.Error); ok {
			if pqError.Code == "23505" {
				os.Exit(1)
			}
			return pqError
		} else {
			return err
		}
	}

	s.config.SetUser(cmd.arguments[0])
	fmt.Printf("the user %v was created at %v\n", user.Name, user.CreatedAt)

	return nil
}

func handlerReset(s *state, cmd command) error{
	if len(cmd.arguments) != 0 {
		return errors.New("0 arguments are expected")
	}

	err := s.database.Reset(context.Background())
	if err!= nil {
		if pqError, ok := err.(*pq.Error); ok {
			return pqError
		} else {
			return err
		}
	}
	return nil
}

func handlerUsers(s *state, cmd command) error{
	if len(cmd.arguments) != 0 {
		return errors.New("0 arguments are expected")
	}

	users, err := s.database.GetUsers(context.Background())
	if err!= nil {
			return err
	} 

	currentUser := s.config.UserName
	for _, user := range users {
		if user.Name == currentUser {
			fmt.Printf("* %v (current)\n", user.Name)
		} else {
			fmt.Printf("* %v\n", user.Name)
		}
	}

	return nil
}

func handlerAgg(s *state, cmd command) error{
	if len(cmd.arguments) != 1 {								//time_between_reqs arugement
		return errors.New("1 arguments are expected")
	}

	timeBetweenRequests, err := time.ParseDuration(cmd.arguments[0])
	if err!= nil {
		return err
	} 

	ticker := time.NewTicker(timeBetweenRequests)
	fmt.Printf("Collecting feeds every %v\n", timeBetweenRequests)
	for ; ; <-ticker.C {
		fmt.Println("\n*** Collecting feed... ***")
		err = scrapeFeeds(s, context.Background())
		if err!= nil {
			return err
		} 

	}

	return nil
}

func handlerAddFeed(s *state, cmd command, user database.User) error{
	if len(cmd.arguments) != 2 {
		return errors.New("2 arguments are expected")
	}

	param := database.CreateFeedParams{ID: uuid.New(), CreatedAt: time.Now(), UpdatedAt: time.Now(), Name: cmd.arguments[0], Url: cmd.arguments[1], UserID: user.ID}
	feed, err := s.database.CreateFeed(context.Background(), param)
	if err!= nil {
		if pqError, ok := err.(*pq.Error); ok {
			if pqError.Code == "23505" {
				os.Exit(1)
			}
			return pqError
		} else {
			return err
		}
	}

	fmt.Printf("the feed %v was created at %v. Url: %v\n", feed.Name, feed.CreatedAt, feed.Url)

	cmd.arguments = cmd.arguments[1:]
	err = handlerFollow(s, cmd, user)
	if err!= nil {
			return err
	} 

	return nil
}

func handlerFeeds(s *state, cmd command) error{
	if len(cmd.arguments) != 0 {
		return errors.New("0 arguments are expected")
	}

	feeds, err := s.database.GetFeeds(context.Background())
	if err!= nil {
			return err
	} 

	for _, feed := range feeds {
		userName, err := s.database.GetUserName(context.Background(), feed.UserID)
		if err!= nil {
			return err
		} 
		fmt.Printf("* Feed  %v (%v). Created by: %v\n", feed.Name, feed.Url, userName)
	}

	return nil
}

func handlerFollow(s *state, cmd command, user database.User) error{
	if len(cmd.arguments) != 1 {
		return errors.New("1 argument is expected")
	}
	
	feed, err := s.database.GetFeedUrl(context.Background(), cmd.arguments[0])
	if err!= nil {
			return err
	} 

	param := database.CreateFeedFollowParams{ID: uuid.New(), CreatedAt: time.Now(), UpdatedAt: time.Now(), UserID: user.ID, FeedID: feed.ID}
	_, err = s.database.CreateFeedFollow(context.Background(), param)
	if err!= nil {
		if pqError, ok := err.(*pq.Error); ok {
			return pqError
		} else {
			return err
		}
	}
	
	fmt.Printf("the user %v just followed the feed %v\n", user.Name, feed.Name)

	return nil
}

func handlerFollowing(s *state, cmd command, user database.User) error{
	if len(cmd.arguments) != 0 {
		return errors.New("0 arguments are expected")
	}

	follows, err := s.database.GetFeedFollowsForUser(context.Background(), user.ID)
	if err!= nil {
			return err
	} 

	currentUser := s.config.UserName
	fmt.Printf("Follows of the user  %v:\n", currentUser)
	for _, follow := range follows {
		fmt.Printf("* Feed  %v\n", follow.FeedName)
	}

	return nil
}


func middlewareLoggedIn(handler func(s *state, cmd command, user database.User) error) func(*state, command) error{
	return func(s *state, cmd command) error {
		user, err := s.database.GetUser(context.Background(), s.config.UserName)
		if err!= nil {
			return err
		}
		return handler(s, cmd, user)
	}
}


func handlerUnfollow(s *state, cmd command, user database.User) error{
	if len(cmd.arguments) != 1 {
		return errors.New("1 argument is expected")
	}
	
	feed, err := s.database.GetFeedUrl(context.Background(), cmd.arguments[0])
	if err!= nil {
			return err
	} 

	param := database.DeleteFeedFollowParams{UserID: user.ID, FeedID: feed.ID}
	err = s.database.DeleteFeedFollow(context.Background(), param)
	if err!= nil {
		if pqError, ok := err.(*pq.Error); ok {
			return pqError
		} else {
			return err
		}
	}
	
	fmt.Printf("the user %v just unfollowed the feed %v\n", user.Name, feed.Name)

	return nil
}

func handlerAggTest(s *state, cmd command) error{
	if len(cmd.arguments) != 0 {
		return errors.New("0 arguments are expected")
	}

	err := scrapeFeeds(s, context.Background())
	if err!= nil {
			return err
	} 

	return nil
}

func handlerBrowse(s *state, cmd command, user database.User) error{
	if len(cmd.arguments) > 1 {
		return errors.New("0 or 1 argument are expected")
	}

	limit := 2
	if len(cmd.arguments) == 1 {
		if i, err := strconv.Atoi(cmd.arguments[0]); err == nil {
			limit = i
		} else {
			return err
		}
	}

	param := database.GetPostForUserParams{UserID: user.ID, Limit: int32(limit)}
	posts, err := s.database.GetPostForUser(context.Background(), param)
	if err!= nil {
		if pqError, ok := err.(*pq.Error); ok {
			return pqError
		} else {
			return err
		}
	}

	fmt.Printf("posts found for user %v:\n", user.Name)
	for _, post := range posts {
		fmt.Printf("%v from %v\n", post.PublishedAt.Time.Format("Mon Jan 2"), post.FeedName)
		fmt.Printf("* %v\n", post.Title)
		fmt.Printf("	%v\n", post.Description.String)
		fmt.Printf("Link: %v\n", post.Url)
		fmt.Println("")
	}

	return nil
}