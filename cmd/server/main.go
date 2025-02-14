package main

import (
	"fmt"
	"internal/pubsub"
	"internal/routing"
	"log"
	"os"
	"os/signal"

	amqp "github.com/rabbitmq/amqp091-go"
)

func main() {
	connString := "amqp://guest:guest@localhost:5672/"
	conn, err := amqp.Dial(connString)
	if err != nil {
		log.Fatalf("Error dialing amqp server: %v", err)
	}
	defer conn.Close()
	fmt.Println("Connection successful.")

	ch, err := conn.Channel()
	if err != nil {
		log.Fatalf("Error creating connection channel: %v", err)
	}
	//don't forget about "peril_direct exchange needed in RabbitMQ web config"
	err = pubsub.PublishJSON(ch, routing.ExchangePerilDirect, routing.PauseKey, routing.PlayingState{IsPaused: true})
	if err != nil {
		log.Fatalf("Error publishing message: %v", err)
	}

	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, os.Interrupt)
	exit := <-signalChan
	log.Fatal("Got signal: ", exit)

}
