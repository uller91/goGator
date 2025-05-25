package main

import (
	"fmt"
	"github.com/uller91/goGator/internal/config"
	"os"
)


func main() {
	//fmt.Print("Hello world!\n")

	//var cfg config.Config
	//cfg, err := config.Read()

	var cfg config.Config
	var st state
	var err error

	cfg, err = config.Read()
	if err != nil {
		fmt.Println(err)
	}
	st.config = &cfg

	var cmds commands
	handlers := make(map[string]func(*state, command) error)
	cmds.handlers = handlers

	cmds.register("login", handlerLogin)

	args := os.Args[:]
	if len(args) < 2 {
		err := fmt.Errorf("Not enough arguments!")
		fmt.Println(err.Error())
		os.Exit(1)
	}

	commandName := args[1]
	commandArgs := []string{}
	if len(args) > 2 {
		commandArgs = args[2:]
	}
	cmd := command{name: commandName, arguments: commandArgs}

	err = cmds.run(&st, cmd)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}


	/* 
	err = st.config.SetUser("Ilya")
	if err != nil {
		fmt.Println(err)
	}
	
	*st.config, err = config.Read()
	fmt.Println("New .json!")
	fmt.Println(st.config)
	*/

}

