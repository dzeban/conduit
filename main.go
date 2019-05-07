package main

import (
	"log"

	"github.com/koding/multiconfig"
)

func main() {
	var config Config
	multiconfig.New().MustLoad(&config)

	server, err := NewServer(config)
	if err != nil {
		log.Fatal("failed to start server: ", err)
	}

	log.Printf("Start listening on %d\n", config.Server.Port)
	server.Run()
}
