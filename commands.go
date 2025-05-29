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
	fmt.Printf("the user %v was created at %v, updated at %v with id %v\n", user.Name, user.CreatedAt, user.UpdatedAt, user.ID)

	return nil
}