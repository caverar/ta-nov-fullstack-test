package raw_stocks

import (
	"backend/pkg/db"
	"database/sql"
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"os"

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

// Service =========================================================================================

type DataInitializer struct {
	client *http.Client
	token  string
	host   string
	db     *sql.DB
}

func NewDataInitializer() (*DataInitializer, error) {
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

func (di *DataInitializer) GetData(next string) (APIResponse, error) {
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
