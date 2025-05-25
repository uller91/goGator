package main

import (
	"fmt"
	"github.com/uller91/goGator/internal/config"
)


func main() {
	//fmt.Print("Hello world!\n")

	var cfg config.Config
	cfg, err := config.Read()
	if err != nil {
		fmt.Println(err)
	}

	cfg.SetUser("Ilya")
	
	cfg, err = config.Read()
	fmt.Println("New .json!")
	fmt.Println(cfg)
}

