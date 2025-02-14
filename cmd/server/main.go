package main

import (
	"fmt"
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
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, os.Interrupt)
	exit := <-signalChan
	log.Fatal("Got signal: ", exit)

}
