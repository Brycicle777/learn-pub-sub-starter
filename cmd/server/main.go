package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"

	"github.com/bootdotdev/learn-pub-sub-starter/internal/gamelogic"
	"github.com/bootdotdev/learn-pub-sub-starter/internal/pubsub"
	"github.com/bootdotdev/learn-pub-sub-starter/internal/routing"

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
	//and peril_topic now too

	_, _, err = pubsub.DeclareAndBind(
		conn,
		routing.ExchangePerilTopic,
		routing.GameLogSlug,
		routing.GameLogSlug+".",
		pubsub.Durable,
	)
	if err != nil {
		log.Fatalf("Error declaring and binding queue: %v", err)
	}

	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, os.Interrupt)
	for {
		select {
		case exit := <-signalChan:
			log.Fatal("Got signal: ", exit)
		default:
			input := gamelogic.GetInput()
			if len(input) > 0 {
				firstWord := input[0]
				switch firstWord {
				case "pause":
					fmt.Println("Sending pause message")
					err = pubsub.PublishJSON(ch, routing.ExchangePerilDirect, routing.PauseKey, routing.PlayingState{IsPaused: true})
					if err != nil {
						log.Fatalf("Error publishing message: %v", err)
					}
				case "resume":
					fmt.Println("Sending resume message")
					err = pubsub.PublishJSON(ch, routing.ExchangePerilDirect, routing.PauseKey, routing.PlayingState{IsPaused: false})
					if err != nil {
						log.Fatalf("Error publishing message: %v", err)
					}
				case "quit":
					fmt.Println("Exiting..")
					return
				default:
					fmt.Println("Unknown command")
				}

			}
		}
	}

}
