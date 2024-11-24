package main

import (
	"fmt"
	"log"

	"github.com/IBM/sarama"
)

func main() {
	// Kafka brokers and topic
	brokers := []string{"localhost:9092"}
	topic := "csv-topic"

	// Producer configuration
	config := sarama.NewConfig()
	config.Producer.RequiredAcks = sarama.WaitForAll        // Wait for all in-sync replicas to acknowledge
	config.Producer.Retry.Max = 5                           // Retry up to 5 times
	config.Producer.Return.Successes = true                 // Required for SyncProducer
	config.Producer.Partitioner = sarama.NewHashPartitioner // Ensure messages are routed to partitions based on the key

	// Create a new sync producer
	producer, err := sarama.NewSyncProducer(brokers, config)
	if err != nil {
		log.Fatalf("Error creating producer: %v", err)
	}
	defer func() {
		if err := producer.Close(); err != nil {
			log.Printf("Error closing producer: %v", err)
		} else {
			fmt.Println("Producer closed successfully")
		}
	}()

	// Get input from the user
	for {
		var consumerID, message string
		fmt.Println("Enter consumer ID (e.g., 'consumer-1'):")
		fmt.Scanln(&consumerID)
		fmt.Println("Enter message:")
		fmt.Scanln(&message)

		// Prepare the message
		msg := &sarama.ProducerMessage{
			Topic: topic,
			Key:   sarama.StringEncoder(consumerID), // Route messages based on the consumer ID
			Value: sarama.StringEncoder(message),
		}

		// Send the message
		partition, offset, err := producer.SendMessage(msg)
		if err != nil {
			log.Printf("Failed to send message: %v", err)
		} else {
			fmt.Printf("Message sent to partition %d with offset %d\n", partition, offset)
		}
	}
}
