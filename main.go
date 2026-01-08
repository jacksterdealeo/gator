package main

import (
	"fmt"
	"gator/internal/config"
	"log"
)

func main() {
	conf, err := config.Read()
	if err != nil {
		log.Fatalln(err)
	}
	if err := conf.SetUser("jack"); err != nil {
		log.Fatalln(err)
	}

	conf, err = config.Read()
	if err != nil {
		log.Fatalln(err)
	}

	fmt.Println("Gator started.\n", *conf)
}
