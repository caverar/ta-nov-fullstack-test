package main

import (
	"backend/db"
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"os"

	"github.com/jmoiron/sqlx"
	"github.com/joho/godotenv"
)

// TYPES ===========================================================================================
type RawItem struct {
	Ticker     string `json:"ticker"`
	TargetFrom string `json:"target_from"`
	TargetTo   string `json:"target_to"`
	Company    string `json:"company"`
	Action     string `json:"action"`
	Brokerage  string `json:"brokerage"`
	RatingFrom string `json:"rating_from"`
	RatingTo   string `json:"rating_to"`
	Time       string `json:"time"`
}
type APIResponse struct {
	Items    []RawItem `json:"items"`
	NextPage string    `json:"next_page"`
}

// DATA INITIALIZER ================================================================================

type DataInitializer struct {
	client *http.Client
	token  string
	host   string
	db     *sqlx.DB
}

func newDataInitializer() (*DataInitializer, error) {
	// Read environment variables
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	// Get the host and token
	host := os.Getenv("DATA_HOST")
	token := os.Getenv("DATA_TOKEN")

	if host == "" || token == "" {
		return nil, errors.New("data credentials not available")
	}

	// Connect to the database
	db := db.Get()

	// Create a new HTTP client
	client := &http.Client{CheckRedirect: http.DefaultClient.CheckRedirect}

	// Return the initializer
	return &DataInitializer{
		client: client,
		token:  token,
		host:   host,
		db:     db,
	}, nil
}

func (di *DataInitializer) getData(next string) (APIResponse, error) {
	// Config the request
	var host string
	if next == "" {
		host = di.host
	} else {
		host = di.host + "?next_page=" + next
	}
	req, err := http.NewRequest("GET", host, nil)
	if err != nil {
		log.Fatal(err)
	}
	req.Header.Add("Authorization", "Bearer "+di.token)

	// Make the request
	resp, err := di.client.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()

	// Verify if the request was successful
	if resp.StatusCode != http.StatusOK {
		log.Fatal("Data request failed", resp.Status)
	}

	// Parse the body
	var result APIResponse
	err = json.NewDecoder(resp.Body).Decode(&result)
	if err != nil {
		log.Fatal(err)
	}
	return result, nil
}

func main() {

	// Start the data initializer
	initializer, err := newDataInitializer()
	if err != nil {
		log.Fatal(err)
	}

	var RawData []RawItem
	var nextPage string
	var counter int
	for {
		resp, err := initializer.getData(nextPage)
		if err != nil {
			log.Fatal(err)
		}
		RawData = append(RawData, resp.Items...)
		nextPage = resp.NextPage
		if nextPage == "" || len(resp.Items) == 0 {
			break
		}
		log.Println("Chunk ", counter, "length: ", len(resp.Items), "next page: ", nextPage)
		counter++
	}
	log.Println("Data initialized", RawData)
}
