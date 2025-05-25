package main

import (
	"github.com/uller91/goGator/internal/config"
	"errors"
	"fmt"
)

type state struct {
	config *config.Config
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

	err := s.config.SetUser(cmd.arguments[0])
	if err != nil {
		return err
	}

	fmt.Printf("%v is set!\n", cmd.arguments[0])
	return nil
}