package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/IBM/sarama"
)

// ConsumerGroupHandler implements sarama.ConsumerGroupHandler
type ConsumerGroupHandler struct {
	consumerID string
}

// Setup is run before consuming
func (ConsumerGroupHandler) Setup(sarama.ConsumerGroupSession) error {
	fmt.Println("Consumer group setup completed")
	return nil
}

// Cleanup is run after consuming
func (ConsumerGroupHandler) Cleanup(sarama.ConsumerGroupSession) error {
	fmt.Println("Consumer group cleanup completed")
	return nil
}

// ConsumeClaim is run to consume messages
func (h ConsumerGroupHandler) ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	for message := range claim.Messages() {
		// Filter messages based on consumerID (key)
		if string(message.Key) == h.consumerID {
			// Print message for the consumer with matching ID
			fmt.Printf("Consumer ID %s received message: %s from topic: %s\n", h.consumerID, string(message.Value), message.Topic)
			session.MarkMessage(message, "")
		}
	}
	return nil
}

func main() {
	// Kafka brokers and topic configuration
	brokers := []string{"localhost:9092"}
	topics := []string{"csv-topic"} // Kafka topic name

	// Get the consumer ID to listen to specific messages
	var consumerID string
	fmt.Println("Enter your consumer ID to receive messages (e.g., 'consumer-1'):")
	fmt.Scanln(&consumerID)

	// Create Kafka consumer group configuration
	config := sarama.NewConfig()
	config.Consumer.Group.Rebalance.Strategy = sarama.BalanceStrategyRoundRobin
	config.Consumer.Offsets.Initial = sarama.OffsetOldest
	config.Consumer.Return.Errors = true

	// Create Kafka consumer group
	consumerGroup, err := sarama.NewConsumerGroup(brokers, "my-consumer-group", config)
	if err != nil {
		log.Fatalf("Failed to create consumer group: %v", err)
	}
	defer func() {
		if err := consumerGroup.Close(); err != nil {
			log.Printf("Error closing consumer group: %v", err)
		} else {
			fmt.Println("Consumer group closed successfully")
		}
	}()

	// Handle graceful shutdown on SIGINT (Ctrl+C)
	sigchan := make(chan os.Signal, 1)
	signal.Notify(sigchan, syscall.SIGINT, syscall.SIGTERM)

	// Create the consumer handler with the consumerID to filter messages
	handler := ConsumerGroupHandler{consumerID: consumerID}

	// Start consuming messages in a separate goroutine
	go func() {
		for {
			ctx := context.Background()
			if err := consumerGroup.Consume(ctx, topics, handler); err != nil {
				log.Printf("Error consuming messages: %v", err)
				break
			}
		}
	}()

	<-sigchan // Wait for termination signal
	fmt.Println("Received termination signal, shutting down consumer...")
}
