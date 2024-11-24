package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

const apiKey = "4f7785b7427e48a9a4931f480d14183e"
const apiURL = "https://newsapi.org/v2/"

type NewsAPIResponse struct {
	Status       string `json:"status"`
	TotalResults int    `json:"totalResults"`
	Articles     []struct {
		Source struct {
			ID   string `json:"id"`
			Name string `json:"name"`
		} `json:"source"`
		Author      string `json:"author"`
		Title       string `json:"title"`
		Description string `json:"description"`
		URL         string `json:"url"`
		PublishedAt string `json:"publishedAt"`
		Content     string `json:"content"`
	} `json:"articles"`
}

// Fetch top headlines with specific query
func getTopHeadlines(query string, sources string, category string, language string, country string) (*NewsAPIResponse, error) {
	url := fmt.Sprintf("%stop-headlines?q=%s&sources=%s&category=%s&language=%s&country=%s&apiKey=%s",
		apiURL, query, sources, category, language, country, apiKey)

	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("error: failed to fetch data, status code: %d", resp.StatusCode)
	}

	// Read response body
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	// Parse JSON
	var result NewsAPIResponse
	err = json.Unmarshal(body, &result)
	if err != nil {
		return nil, err
	}

	return &result, nil
}

func main() {
	// Example: Get top headlines for 'bitcoin' from BBC News and The Verge in business category
	query := "bitcoin"
	sources := "bbc-news,the-verge"
	category := "business"
	language := "en"
	country := "us"

	// Fetch top headlines
	headlines, err := getTopHeadlines(query, sources, category, language, country)
	if err != nil {
		log.Fatalf("Error fetching data: %v", err)
	}

	// Print articles
	for _, article := range headlines.Articles {
		fmt.Printf("Title: %s\n", article.Title)
		fmt.Printf("Description: %s\n", article.Description)
		fmt.Printf("Published At: %s\n", article.PublishedAt)
		fmt.Printf("URL: %s\n\n", article.URL)
	}

	// Fetch all articles with 'bitcoin' from specific sources, within date range
	// Add similar logic to fetch articles
	// You can use the same structure but modify the URL with `getEverything` or other endpoints.
}
