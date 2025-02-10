package main

import (
	"log"
	"github.com/josemontano1996/ai-chatbot-backend/api"
)

const (
	serverAddress = "0.0.0.0:8080"
)

func main() {

	server, err := api.NewServer()

	if err != nil {
		log.Fatal("cannot start server: ", err)
	}

	err = server.Start(serverAddress)

	if err != nil {
		log.Fatal("cannot start server: ", err)
	}
}
