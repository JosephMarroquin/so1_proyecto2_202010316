package main

import (
	"context"
	"fmt"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/segmentio/kafka-go"
)

func sendMessage(c *fiber.Ctx) error {
	// Parse request body
	type Message struct {
		Topic   string `json:"topic"`
		Message string `json:"message"`
	}
	var reqBody Message
	if err := c.BodyParser(&reqBody); err != nil {
		return err
	}

	// Load Kafka brokers from environment variable
	kafkaBrokers := os.Getenv("KAFKA_BROKERS")
	if kafkaBrokers == "" {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "KAFKA_BROKERS environment variable is not set",
		})
	}

	// Create a new Kafka writer
	w := kafka.NewWriter(kafka.WriterConfig{
		Brokers: []string{kafkaBrokers},
		Topic:   reqBody.Topic,
	})

	defer w.Close()

	// Produce message to Kafka topic
	err := w.WriteMessages(context.Background(), kafka.Message{
		Value: []byte(reqBody.Message),
	})

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": fmt.Sprintf("Error sending message: %s", err),
		})
	}

	return c.JSON(fiber.Map{
		"message": "Message sent successfully",
	})
}

func main() {
	// Create a new Fiber instance
	app := fiber.New()

	// Define route to send messages
	app.Post("/send-message", sendMessage)

	// Start the server
	port := os.Getenv("PORT")
	if port == "" {
		port = "3000" // Default port if PORT environment variable is not set
	}
	app.Listen(":" + port)
}
