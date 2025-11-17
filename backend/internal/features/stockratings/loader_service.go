package stockratings

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"slices"
	"strings"
	"time"

	"backend/internal/repository"

	"github.com/jackc/pgx/v5/pgtype"
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
func (s *LoaderService) clearRawStockRatings() error {
	err := s.repo.ClearRawStockRating(context.Background())
	if err != nil {
		return err
	}
	return nil
}

func (s *LoaderService) clearStockRatings() error {
	err := s.repo.ClearStockRating(context.Background())
	if err != nil {
		return err
	}
	return nil
}

var buyRating = []string{"Strong-Buy", "Overweight", "Outperform", "Outperformer", "Market Outperform", "Sector Outperform", "Buy", "Positive", "Speculative Buy"}
var holdRating = []string{"Market Perform", "In-Line", "Hold", "Neutral", "Equal Weight", "Sector Weight", "Sector Perform", "Peer Perform"}
var sellRating = []string{"Underperform", "Underweight", "Sell", "Negative", "Reduce", "Sector Underperform"}
var pendingRating = []string{""}

// Translate the raw rating to a a normalized rating
func (s *LoaderService) rawRatingToStockRating(rawRating string) (repository.StockRatingType, error) {
	if slices.Contains(buyRating, rawRating) {
		return repository.StockRatingTypeBuy, nil
	} else if slices.Contains(holdRating, rawRating) {
		return repository.StockRatingTypeHold, nil
	} else if slices.Contains(sellRating, rawRating) {
		return repository.StockRatingTypeSell, nil
	} else if slices.Contains(pendingRating, rawRating) {
		return repository.StockRatingTypePending, nil
	} else {
		log.Println("Unknown rating: ", rawRating)
		return "", fmt.Errorf("unknown rating: %s", rawRating)
	}
}

var upAction = []string{"target raised by", "upgraded by"}
var downAction = []string{"target lowered by", "downgraded by"}
var reiteratedAction = []string{"target set by", "reiterated by", "initiated by"}

// Translate the raw action to a a normalized action
func (s *LoaderService) rawActionToStockAction(rawAction string) (repository.StockActionType, error) {
	if slices.Contains(upAction, rawAction) {
		return repository.StockActionTypeUp, nil
	} else if slices.Contains(downAction, rawAction) {
		return repository.StockActionTypeDown, nil
	} else if slices.Contains(reiteratedAction, rawAction) {
		return repository.StockActionTypeReiterated, nil
	} else {
		return "", fmt.Errorf("unknown action: %s", rawAction)
	}
}
func (s *LoaderService) rawTargetToStockTarget(rawTarget string) (pgtype.Numeric, error) {
	// Remove currency symbol, commas and trim spaces
	clean := strings.ReplaceAll(rawTarget, "$", "")
	clean = strings.ReplaceAll(clean, ",", "")
	clean = strings.TrimSpace(clean)

	if clean == "" {
		return pgtype.Numeric{}, fmt.Errorf("empty target value")
	}

	var n pgtype.Numeric
	// pgtype.Numeric implements Set which accepts numeric strings
	if err := n.Scan(clean); err != nil {
		return pgtype.Numeric{}, err
	}
	return n, nil
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
	clearRawStockRatingsError
	clearStockRatingsError
	dataFetchError
	timeParseError
	insertRawStockRatingsError
	insertStockRatingsError
	unknownRatingError
	unknownActionError
	unknownTargetError
)

type InitDataError struct {
	kind initDataErrorKind
	err  error
}

func (e InitDataError) Error() string {
	switch e.kind {
	case clearRawStockRatingsError:
		return fmt.Sprintf("Failed to clear current raw data: %s", e.err.Error())
	case clearStockRatingsError:
		return fmt.Sprintf("Failed to clear current data: %s", e.err.Error())
	case dataFetchError:
		return fmt.Sprintf("Failed to get data from API: %s", e.err.Error())
	case timeParseError:
		return fmt.Sprintf("Failed to parse time from API: %s", e.err.Error())
	case insertRawStockRatingsError:
		return fmt.Sprintf("Failed to insert data into database: %s", e.err.Error())
	case insertStockRatingsError:
		return fmt.Sprintf("Failed to insert data into database: %s", e.err.Error())
	case unknownRatingError:
		return fmt.Sprintf("Failed to parse rating from API: %s", e.err.Error())
	case unknownActionError:
		return fmt.Sprintf("Failed to parse action from API: %s", e.err.Error())
	case unknownTargetError:
		return fmt.Sprintf("Failed to parse target from API: %s", e.err.Error())
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
	ClearRawStockRatingsError  = InitDataError{kind: clearRawStockRatingsError}
	DataFetchError             = InitDataError{kind: dataFetchError}
	TimeParseError             = InitDataError{kind: timeParseError}
	InsertRawStockRatingsError = InitDataError{kind: insertRawStockRatingsError}
	InsertStockRatingsError    = InitDataError{kind: insertStockRatingsError}
	UnknownRatingError         = InitDataError{kind: unknownRatingError}
	UnknownActionError         = InitDataError{kind: unknownActionError}
	UnknownTargetError         = InitDataError{kind: unknownTargetError}
)

// Method ------------------------------------------------------------------------------------------

func (s *LoaderService) InitData() error {

	// Clear the current data in db
	err := s.clearRawStockRatings()
	if err != nil {
		return ClearRawStockRatingsError.From(err)
	}
	err = s.clearStockRatings()
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

			// Insert it into the raw db
			_, err = s.repo.AddRawStockRatings(context.Background(), stocksRatings)
			if err != nil {
				return InsertRawStockRatingsError.From(err)
			}

			var parsedStocksRatings []repository.AddStockRatingsParams
			// Normalize the data
			for _, rating := range stocksRatings {
				ratingFrom, err := s.rawRatingToStockRating(rating.RatingFrom)
				if err != nil {
					log.Println("Error parsing rating: ", rating)
					return UnknownRatingError.From(err)
				}
				ratingTo, err := s.rawRatingToStockRating(rating.RatingTo)
				if err != nil {
					log.Println("Error parsing rating: ", rating)
					return UnknownRatingError.From(err)
				}
				action, err := s.rawActionToStockAction(rating.Action)
				if err != nil {
					log.Println("Error parsing rating: ", rating)
					return UnknownActionError.From(err)
				}
				targetFrom, err := s.rawTargetToStockTarget(rating.TargetFrom)
				if err != nil {
					log.Println("Error parsing rating: ", rating)
					return UnknownTargetError.From(err)
				}
				targetTo, err := s.rawTargetToStockTarget(rating.TargetTo)
				if err != nil {
					log.Println("Error parsing rating: ", rating)
					return UnknownTargetError.From(err)
				}

				parsedStocksRatings = append(parsedStocksRatings, repository.AddStockRatingsParams{
					Ticker:     rating.Ticker,
					Company:    rating.Company,
					TargetFrom: targetFrom,
					TargetTo:   targetTo,
					Action:     action,
					RatingFrom: ratingFrom,
					RatingTo:   ratingTo,
					At:         rating.At,
				})
			}

			// Insert it into the db
			_, err = s.repo.AddStockRatings(context.Background(), parsedStocksRatings)
			if err != nil {
				return InsertStockRatingsError.From(err)
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
