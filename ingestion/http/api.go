package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"

	"github.com/IBM/sarama"
)

// User represents the user data fetched from the API
type User struct {
	ID    int    `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
}

// FetchUsers fetches user data from JSONPlaceholder API
func FetchUsers() ([]User, error) {
	url := "https://jsonplaceholder.typicode.com/users"
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("failed to fetch users: %s", string(body))
	}

	var users []User
	err = json.NewDecoder(resp.Body).Decode(&users)
	if err != nil {
		return nil, err
	}

	return users, nil
}

func main() {
	// Kafka Configuration
	brokers := []string{"localhost:9092"} // Replace with your Kafka broker addresses
	topic := "csv-topic"                  // Kafka topic name

	// Create a new Kafka producer with configuration
	config := sarama.NewConfig()
	config.Producer.Return.Successes = true // Ensure the producer waits for success confirmation
	config.Producer.Partitioner = sarama.NewHashPartitioner

	producer, err := sarama.NewSyncProducer(brokers, config)
	if err != nil {
		log.Fatalf("Failed to start Kafka producer: %v", err)
	}
	defer producer.Close()

	// Get the consumer ID from user input
	fmt.Print("Enter consumer ID (e.g., 'consumer-1') to send data to: ")
	var consumerID string
	fmt.Scanln(&consumerID)

	// Fetch users from the JSONPlaceholder API
	users, err := FetchUsers()
	if err != nil {
		log.Fatalf("Error fetching users: %v", err)
	}

	// Send each user as a message to Kafka
	for _, user := range users {
		// Create a JSON string with user data
		userData, err := json.Marshal(user)
		if err != nil {
			log.Printf("Failed to marshal user data: %v", err)
			continue
		}

		// Send message to Kafka with consumerID as key
		message := &sarama.ProducerMessage{
			Topic: topic,
			Key:   sarama.StringEncoder(consumerID), // Route messages based on consumerID
			Value: sarama.StringEncoder(string(userData)),
		}

		_, _, err = producer.SendMessage(message)
		if err != nil {
			log.Printf("Failed to send message for user ID %d: %v", user.ID, err)
		} else {
			log.Printf("Sent user data to Kafka: %s", string(userData))
		}
	}

	log.Println("All user data has been sent to Kafka.")
}
