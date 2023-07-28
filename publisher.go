package main

import (
	"context"
	"encoding/json"
	amqp "github.com/rabbitmq/amqp091-go"
	"log"
)

var (
	rabbitMQURL  = "amqp://guest:guest@localhost:5672/" // Replace with your RabbitMQ URL if needed
	exchangeName = "notification_exchange"
	exchangeType = "direct"
	routingKey   = "notification_key"
)

func main() {
	// Connect to RabbitMQ
	conn, err := amqp.Dial(rabbitMQURL)
	if err != nil {
		log.Fatalf("Failed to connect to RabbitMQ: %v", err)
	}
	defer conn.Close()

	// Create a channel
	ch, err := conn.Channel()
	if err != nil {
		log.Fatalf("Failed to open a channel: %v", err)
	}
	defer ch.Close()

	// Declare the exchange
	err = ch.ExchangeDeclare(exchangeName, exchangeType, true, false, false, false, nil)
	if err != nil {
		log.Fatalf("Failed to declare the exchange: %v", err)
	}

	// Send a sample message
	messages := map[string]any{
		"verify": map[string]any{
			"email":   "verify@example.com",
			"message": "Verify your email address",
		},
		"welcome": map[string]any{
			"email":   "welcome@example.com",
			"message": "Welcome to our platform",
		},
	}

	for _, message := range messages {
		//panic(message.(map[string]any)["email"])
		err = publishMessage(ch, message)
		if err != nil {
			log.Fatalf("Failed to publish a message: %v", err)
		}

		log.Printf("A message was sent: %v", message)
	}

}

func publishMessage(ch *amqp.Channel, body any) error {
	//marshal message
	message, err := json.Marshal(body)
	if err != nil {
		return err
	}

	return ch.PublishWithContext(
		context.TODO(),
		exchangeName,
		routingKey,
		false,
		false,
		amqp.Publishing{
			ContentType: "application/json",
			Body:        message,
		},
	)
}
