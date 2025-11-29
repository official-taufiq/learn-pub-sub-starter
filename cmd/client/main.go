package main

import (
	"fmt"
	"log"

	"github.com/official-taufiq/learn-pub-sub-starter/internal/gamelogic"
	"github.com/official-taufiq/learn-pub-sub-starter/internal/pubsub"
	"github.com/official-taufiq/learn-pub-sub-starter/internal/routing"

	amqp "github.com/rabbitmq/amqp091-go"
)

func main() {
	fmt.Println("Starting Peril client...")

	const connectionStr = "amqp://guest:guest@localhost:5672/"

	conn, err := amqp.Dial(connectionStr)
	if err != nil {
		log.Fatalf("Failed to connect to RabbitMQ:%v", err)
	}
	defer conn.Close()
	fmt.Println("Connected to RabbitMQ server")

	username, err := gamelogic.ClientWelcome()
	if err != nil {
		log.Fatalf("Couldn't get username: %v", err)
	}

	_, que, err := pubsub.DeclareAndBind(
		conn,
		routing.ExchangePerilDirect,
		routing.PauseKey+"."+username,
		routing.PauseKey,
		pubsub.SimpleQueueTransient,
	)
	if err != nil {
		log.Fatalf("could not subscribe to pause: %v", err)
	}

	fmt.Printf("Subscribed to pause messages on queue %s\n", que.Name)

	gameState := gamelogic.NewGameState(username)
	err = pubsub.SubscribeJSON(conn, routing.ExchangePerilDirect, routing.PauseKey+"."+gameState.GetUsername(), routing.PauseKey, pubsub.SimpleQueueTransient, handlerPause(gameState))
	if err != nil {
		log.Fatalf("could not subscribe to pause messages: %v", err)
	}
	for {
		words := gamelogic.GetInput()

		if len(words) == 0 {
			continue
		}
		switch words[0] {
		case "spawn":
			err = gameState.CommandSpawn(words)
			if err != nil {
				fmt.Println(err)
				continue
			}
		case "move":
			_, err = gameState.CommandMove(words)
			if err != nil {
				fmt.Println(err)
				continue
			}
			// fmt.Printf("Moving %d unit(s) to %s\n", len(armyMove.Units), armyMove.ToLocation)
		case "status":
			gameState.CommandStatus()
		case "help":
			gamelogic.PrintClientHelp()
		case "spam":
			fmt.Print("Spamming not allowed yet!\n")
		case "quit":
			gamelogic.PrintQuit()
			return
		default:
			log.Println("Unknown command. Type 'help' for possible commands.")
			continue
		}
	}
}
