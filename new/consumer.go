package main

import (
	"fmt"

	"github.com/IBM/sarama"
)

func main() {

	brokers := []string{"localhost:"}

	// Kafka topic to use
	topic := "csv-topic"

	// Set up consumer
	consumer, err := connectConsumer(brokers)
	if err != nil {
		panic(err)
	}
	defer consumer.Close()

	// Start consumer to consume messages
	startConsumer(consumer, topic)
}

// connectConsumer creates a new Kafka consumer.
func connectConsumer(brokers []string) (sarama.Consumer, error) {
	config := sarama.NewConfig()
	config.Consumer.Return.Errors = true
	consumer, err := sarama.NewConsumer(brokers, config)
	if err != nil {
		return nil, err
	}
	return consumer, nil
}

// startConsumer consumes messages from Kafka and prints them.
func startConsumer(consumer sarama.Consumer, topic string) {
	partitionConsumer, err := consumer.ConsumePartition(topic, 0, sarama.OffsetOldest)
	if err != nil {
		panic(err)
	}
	defer partitionConsumer.Close()

	fmt.Println("Consumer started, waiting for messages...")

	// Consume messages
	for msg := range partitionConsumer.Messages() {
		fmt.Printf("Received message: %s\n", string(msg.Value))
	}
}
