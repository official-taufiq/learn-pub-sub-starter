package pubsub

import (
	"encoding/json"
	"fmt"

	amqp "github.com/rabbitmq/amqp091-go"
)

func SubscribeJSON[T any](conn *amqp.Connection, exchange, queueName, key string, queueType SimpleQueueType, handler func(T)) error {

	ch, que, err := DeclareAndBind(conn, exchange, queueName, key, queueType)
	if err != nil {
		return fmt.Errorf("failed to declare and bind queue: %v", err)
	}
	amqpDel, err := ch.Consume(que.Name, "", false, false, false, false, nil)
	if err != nil {
		return fmt.Errorf("could not consume messages: %v", err)
	}

	go func() {
		defer ch.Close()
		for del := range amqpDel {
			var msg T
			err := json.Unmarshal(del.Body, &msg)
			if err != nil {
				fmt.Printf("Failed to unmarshal JSON: %v", err)
				continue
			}
			handler(msg)
			if err := del.Ack(false); err != nil {
				fmt.Printf("Failed to ack message: %v", err)
			}
		}
	}()

	return nil
}
