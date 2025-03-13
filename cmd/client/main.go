package main

import (
	"fmt"
	"log"

	"github.com/bootdotdev/learn-pub-sub-starter/internal/gamelogic"

	amqp "github.com/rabbitmq/amqp091-go"
)

func main() {
	fmt.Println("Starting Peril client...")
	connString := "amqp://guest:guest@localhost:5672/"
	conn, err := amqp.Dial(connString)
	if err != nil {
		log.Fatalf("Error dialing amqp server: %v", err)
	}
	defer conn.Close()
	fmt.Println("Connection successful.")

	_, err = gamelogic.ClientWelcome()
	if err != nil {
		log.Fatalf("Error accepting username: %v", err)
	}
}
