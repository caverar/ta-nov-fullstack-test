package stockratings

import (
	"context"
	"encoding/json"
	"fmt"
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

// SERVICE =========================================================================================

type LoaderService struct {
	client *http.Client
	token  string
	host   string
	repo   *repository.Queries
}

func NewLoaderService(r *repository.Queries) *LoaderService {
	// Get the host and token
	host := os.Getenv("DATA_HOST")
	token := os.Getenv("DATA_TOKEN")

	if host == "" || token == "" {
		log.Fatal("Missing host/token for stockRatingsLoaders")
	}

	// Create a new HTTP client
	client := &http.Client{CheckRedirect: http.DefaultClient.CheckRedirect}

	// Return the initializer
	return &LoaderService{
		client: client,
		token:  token,
		host:   host,
		repo:   r,
	}
}

// utils ===========================================================================================
func (s *LoaderService) clearStockRatings() error {
	err := s.repo.DeleteRawStockRatings(context.Background())
	if err != nil {
		return err
	}
	return nil
}

// getData =========================================================================================

// Errors ------------------------------------------------------------------------------------------
type getDataErrorKind int

const (
	_ getDataErrorKind = iota
	apiError
	jSONParseError
)

type GetDataError struct {
	kind getDataErrorKind
	err  error
}

func (e GetDataError) Error() string {
	switch e.kind {
	case apiError:
		return fmt.Sprintf("Failed to request external API: %s", e.err.Error())
	case jSONParseError:
		return fmt.Sprintf("Failed to parse data from API: %s", e.err.Error())
	default:
		return "Unknown error"
	}
}

func (e GetDataError) From(err error) GetDataError {
	e1 := e
	e1.err = err
	return e1
}

var (
	APIError       = GetDataError{kind: apiError}
	JSONParseError = GetDataError{kind: jSONParseError}
)

// Method ------------------------------------------------------------------------------------------
// Get the data from the API
func (s *LoaderService) getData(cursor string) (APIResponse, error) {
	// Config the request
	var host string
	if cursor == "" {
		host = s.host
	} else {
		host = s.host + "?next_page=" + cursor
	}
	req, err := http.NewRequest("GET", host, nil)
	if err != nil {
		log.Fatal("Bad request build to API", err)
	}
	req.Header.Add("Authorization", "Bearer "+s.token)

	// Make the request
	resp, err := s.client.Do(req)
	if err != nil {
		return APIResponse{}, APIError.From(err)
	}
	defer resp.Body.Close()

	// Verify if the request was successful
	if resp.StatusCode != http.StatusOK {
		var err = fmt.Errorf("status code: %d", resp.StatusCode)
		return APIResponse{}, APIError.From(err)
	}

	// Parse the body
	var result APIResponse
	err = json.NewDecoder(resp.Body).Decode(&result)
	if err != nil {
		return APIResponse{}, JSONParseError.From(err)
	}
	return result, nil
}

// InitData ========================================================================================
// Errors ------------------------------------------------------------------------------------------
type initDataErrorKind int

const (
	_ initDataErrorKind = iota
	clearStockRatingsError
	dataFetchError
	timeParseError
	insertRawStockRatingsError
)

type InitDataError struct {
	kind initDataErrorKind
	err  error
}

func (e InitDataError) Error() string {
	switch e.kind {
	case clearStockRatingsError:
		return fmt.Sprintf("Failed to clear current data: %s", e.err.Error())
	case dataFetchError:
		return fmt.Sprintf("Failed to get data from API: %s", e.err.Error())
	case timeParseError:
		return fmt.Sprintf("Failed to parse time from API: %s", e.err.Error())
	case insertRawStockRatingsError:
		return fmt.Sprintf("Failed to insert data into database: %s", e.err.Error())
	default:
		return "Unknown error"
	}
}

func (e InitDataError) From(err error) InitDataError {
	e1 := e
	e1.err = err
	return e1
}
func (e InitDataError) Unwrap() error {
	return e.err
}

var (
	ClearStockRatingsError     = InitDataError{kind: clearStockRatingsError}
	DataFetchError             = InitDataError{kind: dataFetchError}
	TimeParseError             = InitDataError{kind: timeParseError}
	InsertRawStockRatingsError = InitDataError{kind: insertRawStockRatingsError}
)

// Method ------------------------------------------------------------------------------------------

func (s *LoaderService) InitData() error {

	// Clear the current data
	err := s.clearStockRatings()
	if err != nil {
		return ClearStockRatingsError.From(err)
	}

	// Download and insert by chunks
	var nextPage string
	var counter int
	for {
		// Get the data
		resp, err := s.getData(nextPage)
		if err != nil {
			return DataFetchError.From(err)
		}

		// If not void insert it into the database
		if len(resp.Items) > 0 {

			var stocksRatings []repository.AddRawStockRatingsParams
			for _, rating := range resp.Items {
				at, err := time.Parse(time.RFC3339Nano, rating.Time)
				if err != nil {
					return TimeParseError.From(err)
				}
				stocksRatings = append(stocksRatings, repository.AddRawStockRatingsParams{
					Ticker:     rating.Ticker,
					TargetFrom: rating.TargetFrom,
					TargetTo:   rating.TargetTo,
					Company:    rating.Company,
					Action:     rating.Action,
					Brokerage:  rating.Brokerage,
					RatingFrom: rating.RatingFrom,
					RatingTo:   rating.RatingTo,
					At:         at,
				})
			}

			_, err = s.repo.AddRawStockRatings(context.Background(), stocksRatings)
			if err != nil {
				return InsertRawStockRatingsError.From(err)
			}
		}

		nextPage = resp.NextPage
		if nextPage == "" || len(resp.Items) == 0 {
			break
		}
		log.Println("Chunk ", counter, "length: ", len(resp.Items), "next page: ", nextPage)
		counter++
	}
	log.Println("Data initialized")

	return nil

}
