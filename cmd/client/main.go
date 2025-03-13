package main

import (
	"fmt"
	"log"

	"github.com/bootdotdev/learn-pub-sub-starter/internal/gamelogic"
	"github.com/bootdotdev/learn-pub-sub-starter/internal/pubsub"
	"github.com/bootdotdev/learn-pub-sub-starter/internal/routing"

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

	username, err := gamelogic.ClientWelcome()
	if err != nil {
		log.Fatalf("Error accepting username: %v", err)
	}

	_, _, err = pubsub.DeclareAndBind(
		conn,
		routing.ExchangePerilDirect,
		routing.PauseKey+"."+username,
		routing.PauseKey,
		pubsub.Transient,
	)
	if err != nil {
		log.Fatalf("Error declaring and binding queue: %v", err)
	}

	gameState := gamelogic.NewGameState(username)

	for {
		input := gamelogic.GetInput()
		if len(input) > 0 {
			cmd := input[0]
			switch cmd {
			case "spawn":
				err := gameState.CommandSpawn(input)
				if err != nil {
					fmt.Println("Error: ", err)
				}
			case "move":
				mv, err := gameState.CommandMove(input)
				if err != nil {
					fmt.Println("Error: ", err)
				} else {
					fmt.Printf("%v\n", mv)
				}
			case "status":
				gameState.CommandStatus()
			case "help":
				gamelogic.PrintClientHelp()
			case "spam":
				fmt.Println("Spamming not allowed yet!")
			case "quit":
				gamelogic.PrintQuit()
				return
			default:
				fmt.Println("Unknown command")
			}
		}
	}
}
