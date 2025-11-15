package stockloader

import (
	"context"
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"os"
	"time"

	"backend/internal/repository"
)

// TYPES ===========================================================================================
type RawStockEvent struct {
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
	Items    []RawStockEvent `json:"items"`
	NextPage string          `json:"next_page"`
}

// ERRORS ==========================================================================================

// TODO Make good error handling using this:
var (
	ErrMissingCredentials = errors.New("missing host/token for stockloader")
	ErrAPIRequest         = errors.New("failed to request external API")
	ErrBadStatus          = errors.New("unexpected status from API")
	ErrDecode             = errors.New("failed to decode response")
	ErrTimeParse          = errors.New("invalid time format")
)

// SERVICE =========================================================================================

type Service struct {
	client *http.Client
	token  string
	host   string
	repo   *repository.Queries
}

func NewService(r *repository.Queries) (*Service, error) {
	// Get the host and token
	host := os.Getenv("DATA_HOST")
	token := os.Getenv("DATA_TOKEN")

	if host == "" || token == "" {
		return nil, errors.New("data credentials not available")
	}

	// Create a new HTTP client
	client := &http.Client{CheckRedirect: http.DefaultClient.CheckRedirect}

	// Return the initializer
	return &Service{
		client: client,
		token:  token,
		host:   host,
		repo:   r,
	}, nil
}

// Get the data from the API
func (s *Service) getData(cursor string) (APIResponse, error) {
	// Config the request
	var host string
	if cursor == "" {
		host = s.host
	} else {
		host = s.host + "?next_page=" + cursor
	}
	req, err := http.NewRequest("GET", host, nil)
	if err != nil {
		log.Fatal(err)
	}
	req.Header.Add("Authorization", "Bearer "+s.token)

	// Make the request
	resp, err := s.client.Do(req)
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

func (s *Service) clearEvents() error {
	err := s.repo.DeleteAllStockEvents(context.Background())
	if err != nil {
		return err
	}
	return nil
}

func (s *Service) InitData() error {

	// Clear the current data
	err := s.clearEvents()
	if err != nil {
		return err
	}

	// Download and insert by chunks
	var RawData []RawStockEvent
	var nextPage string
	var counter int
	for {
		// Get the data
		resp, err := s.getData(nextPage)
		if err != nil {
			log.Fatal(err)
		}

		// If not void insert it into the database
		if len(resp.Items) > 0 {

			var stocksEvents []repository.AddStockEventsParams
			for _, stockEvent := range resp.Items {
				at, err := time.Parse(time.RFC3339Nano, stockEvent.Time)
				if err != nil {
					log.Fatal(err)
				}
				stocksEvents = append(stocksEvents, repository.AddStockEventsParams{
					Ticker:     stockEvent.Ticker,
					TargetFrom: stockEvent.TargetFrom,
					TargetTo:   stockEvent.TargetTo,
					Company:    stockEvent.Company,
					Action:     stockEvent.Action,
					Brokerage:  stockEvent.Brokerage,
					RatingFrom: stockEvent.RatingFrom,
					RatingTo:   stockEvent.RatingTo,
					At:         at,
				})
			}

			_, err = s.repo.AddStockEvents(context.Background(), stocksEvents)
			if err != nil {
				log.Fatal(err)
			}
		}

		nextPage = resp.NextPage
		if nextPage == "" || len(resp.Items) == 0 {
			break
		}
		log.Println("Chunk ", counter, "length: ", len(resp.Items), "next page: ", nextPage)
		counter++
	}
	log.Println("Data initialized", RawData)

	return nil

}
