package main

import (
	"fmt"
	"os"
	"time"

	"github.com/IBM/sarama"
)

func main() {
	// Kafka Broker address
	brokers := []string{"localhost:9092"}

	// Kafka topic to use
	topic := "csv-topic"

	// Set up producer
	producer, err := connectProducer(brokers)
	if err != nil {
		panic(err)
	}
	defer producer.Close()

	// Send a CSV file as a message
	err = sendCSV(producer, topic, "D:\\projects\\go\\hackathon\\ingestion\\Cleaned_Students_Performance.csv")
	if err != nil {
		panic(err)
	}

	fmt.Println("CSV file sent to Kafka!")

}

// connectProducer creates a new Kafka producer.
func connectProducer(brokers []string) (sarama.SyncProducer, error) {
	config := sarama.NewConfig()
	config.Producer.Return.Successes = true
	producer, err := sarama.NewSyncProducer(brokers, config)
	if err != nil {
		return nil, err
	}
	return producer, nil
}

// sendCSV reads a CSV file and sends each line as a message to Kafka.
func sendCSV(producer sarama.SyncProducer, topic, filepath string) error {
	// Open the CSV file
	file, err := os.Open(filepath)
	if err != nil {
		return err
	}
	defer file.Close()

	// Read the file line by line
	buf := make([]byte, 4096)
	for {
		n, err := file.Read(buf)
		if err != nil && err.Error() != "EOF" {
			return err
		}

		if n == 0 {
			break
		}

		// Send the data to Kafka
		message := &sarama.ProducerMessage{
			Topic: topic,
			Value: sarama.StringEncoder(buf[:n]),
		}

		// Produce the message
		_, _, err = producer.SendMessage(message)
		if err != nil {
			return err
		}
		fmt.Printf("Sent message: %s\n", string(buf[:n]))
		time.Sleep(1 * time.Second) // Sleep to simulate sending messages over time
	}

	return nil
}
