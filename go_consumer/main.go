package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/segmentio/kafka-go"
)

func main() {
	// Read environment variables
	kafkaBrokers := strings.Split(os.Getenv("KAFKA_BROKERS"), ",")
	kafkaTopic := os.Getenv("KAFKA_TOPIC")

	// Initialize Kafka reader
	config := kafka.ReaderConfig{
		Brokers: kafkaBrokers,
		Topic:   kafkaTopic,
		GroupID: "test-group",
	}

	reader := kafka.NewReader(config)
	defer reader.Close()

	fmt.Printf("Listening to topic: %s\n", kafkaTopic)

	// Handle OS signals to gracefully shutdown
	sigchan := make(chan os.Signal, 1)
	signal.Notify(sigchan, syscall.SIGINT, syscall.SIGTERM)

	// Start reading messages
	for {
		select {
		case <-sigchan:
			fmt.Println("Shutting down...")
			return
		default:
			msg, err := reader.FetchMessage(context.Background())
			if err != nil {
				fmt.Printf("Error fetching message: %v\n", err)
				continue
			}
			fmt.Printf("Message from topic %s: %s\n", msg.Topic, msg.Value)
			err = reader.CommitMessages(context.Background(), msg)
			if err != nil {
				fmt.Printf("Error committing message: %v\n", err)
			}
		}
	}
}
