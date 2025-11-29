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
	fmt.Println("Starting Peril server...")
	connectionStr := "amqp://guest:guest@localhost:5672/"
	conn, err := amqp.Dial(connectionStr)
	if err != nil {
		panic(err)
	}
	defer conn.Close()

	newChan, err := conn.Channel()
	if err != nil {
		log.Fatalf("Failed to open a channel:%v", err)
	}
	_, que, err := pubsub.DeclareAndBind(
		conn,
		routing.ExchangePerilTopic,
		routing.GameLogSlug,
		routing.GameLogSlug+"."+"*",
		pubsub.SimpleQueueDurable,
	)
	if err != nil {
		log.Fatalf("could not subscribe to pause: %v", err)
	}
	fmt.Printf("Queue %v declared and bound!\n", que.Name)

	gamelogic.PrintServerHelp()

	for {
		command := gamelogic.GetInput()[0]

		if len(command) == 0 {
			continue
		}

		switch command {
		case "pause":
			err = pubsub.PublishJSON(newChan, routing.ExchangePerilDirect, routing.PauseKey, routing.PlayingState{
				IsPaused: true,
			})
			if err != nil {
				log.Printf("could not publish time:%v", err)
			}
			fmt.Println("Pause message sent")

		case "resume":
			err = pubsub.PublishJSON(newChan, routing.ExchangePerilDirect, routing.PauseKey, routing.PlayingState{
				IsPaused: false,
			})
			if err != nil {
				log.Printf("could not publish time:%v", err)
			}
			fmt.Println("Resume message sent")

		case "quit":
			log.Printf("Server quitting...")
			return
		case "help":
			gamelogic.PrintServerHelp()
		default:
			log.Println("Unknown command. Type 'help' for possible commands.")

		}
	}
}
