package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"regexp"
	"strings"

	"github.com/IBM/sarama"
	"github.com/gocolly/colly"
	"github.com/jdkato/prose/v2" // NLP library for basic text processing
)

func main() {
	// Kafka Configuration
	brokers := []string{"localhost:9092"} // Replace with your Kafka broker addresses
	topic := "csv-topic"                  // Kafka topic name

	// Create a new Kafka producer
	producer, err := sarama.NewSyncProducer(brokers, nil)
	if err != nil {
		log.Fatalf("Failed to start Kafka producer: %v", err)
	}
	defer producer.Close()

	// Get user input (URL, data selector, and consumer ID)
	reader := bufio.NewReader(os.Stdin)
	fmt.Print("Enter consumer ID (e.g., 'consumer-1') to send data to: ")
	consumerID, _ := reader.ReadString('\n')
	consumerID = strings.TrimSpace(consumerID)

	fmt.Print("Enter URL to scrape (leave blank for default): ")
	url, _ := reader.ReadString('\n')
	url = strings.TrimSpace(url)
	if url == "" {
		url = "https://example.com" // Default URL if no input is provided
	}

	fmt.Print("Enter CSS selector to scrape (leave blank for automatic): ")
	selector, _ := reader.ReadString('\n')
	selector = strings.TrimSpace(selector)

	// Create a new Colly collector without domain restrictions
	c := colly.NewCollector()

	// Set a custom User-Agent to avoid being blocked
	c.UserAgent = "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/91.0.4472.124 Safari/537.36"

	// Helper function to send messages to Kafka
	sendToKafka := func(data string) {
		message := &sarama.ProducerMessage{
			Topic: topic,
			Key:   sarama.StringEncoder(consumerID), // Route messages based on the consumer ID
			Value: sarama.StringEncoder(data),
		}
		_, _, err := producer.SendMessage(message)
		if err != nil {
			log.Printf("Failed to send message: %v", err)
		} else {
			log.Println("Message sent to Kafka:", data)
		}
	}

	// Smart detection function using NLP (e.g., extracting keywords)
	analyzeText := func(text string) string {
		doc, err := prose.NewDocument(text)
		if err != nil {
			log.Printf("Error processing text with NLP: %v", err)
			return ""
		}
		// Extract keywords or phrases from the text
		keywords := []string{}
		for _, tok := range doc.Tokens() {
			if tok.Tag == "NN" || tok.Tag == "NNS" { // Only extracting nouns (could improve with more NLP)
				keywords = append(keywords, tok.Text)
			}
		}
		return strings.Join(keywords, ", ")
	}

	// Function for smart data extraction (headings, prices, etc.)
	c.OnHTML("body", func(e *colly.HTMLElement) {
		if selector == "" {
			// If no selector provided, perform smart scraping (headings, prices, etc.)
			e.ForEach("h1, h2, h3", func(i int, el *colly.HTMLElement) {
				text := strings.TrimSpace(el.Text)
				if text != "" {
					// Analyze content for better keywords using NLP
					keywords := analyzeText(text)
					sendToKafka(fmt.Sprintf("Heading: %s, Keywords: %s", text, keywords))
				}
			})

			// Extract paragraphs and descriptions
			e.ForEach("p", func(i int, el *colly.HTMLElement) {
				text := strings.TrimSpace(el.Text)
				if len(text) > 30 { // Assuming meaningful paragraphs are at least 30 characters
					// Analyze the paragraph using NLP
					keywords := analyzeText(text)
					sendToKafka(fmt.Sprintf("Paragraph: %s, Keywords: %s", text, keywords))
				}
			})

			// Extract numeric values (likely prices)
			e.ForEach("span, div", func(i int, el *colly.HTMLElement) {
				text := strings.TrimSpace(el.Text)
				// Regex to match possible prices (e.g., $12.99, ₹500)
				priceRegex := regexp.MustCompile(`\b(\$|€|₹|\d{1,3}(,\d{3})*)\s?\d+(\.\d{2})?\b`)
				if priceRegex.MatchString(text) {
					sendToKafka(fmt.Sprintf("Price: %s", text))
				}
			})

		} else {
			// If user provides a specific CSS selector, scrape based on that
			e.ForEach(selector, func(i int, el *colly.HTMLElement) {
				text := strings.TrimSpace(el.Text)
				if text != "" {
					// Analyze with NLP if necessary
					keywords := analyzeText(text)
					sendToKafka(fmt.Sprintf("Custom Data: %s, Keywords: %s", text, keywords))
				}
			})
		}
	})

	// Handle errors during scraping
	c.OnError(func(r *colly.Response, err error) {
		log.Printf("Request to %s failed: %v", r.Request.URL, err)
	})

	// Start scraping from the URL
	err = c.Visit(url)
	if err != nil {
		log.Fatalf("Failed to visit the website: %v", err)
	}

	log.Println("Scraping complete!")
}
