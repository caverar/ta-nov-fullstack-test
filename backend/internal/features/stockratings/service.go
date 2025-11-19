package stockratings

import (
	"backend/internal/repository"
	"context"
	"fmt"
	"time"
)

// SERVICE =========================================================================================

type ServiceInterface interface {
	GetStockRatings(input GetStockRatingsInput) (GetStockRatingsOutput, error)
}
type Service struct {
	repo *repository.Queries
}

func NewService(r *repository.Queries) *Service {
	return &Service{
		repo: r,
	}
}

// GetStockRatings ---------------------------------------------------------------------------------
type GetStockRatingsInput struct {
	sortOrder   string
	sortBy      string
	offset      int32
	limit       int32
	tickerLike  string
	companyLike string
}

type rating = struct {
	ticker      string
	company     string
	brokerage   string
	targetFrom  string
	targetTo    string
	action      string
	rawAction   string
	ratingFrom  string
	ratingTo    string
	at          time.Time
	targetDelta string
	score       int32
}
type GetStockRatingsOutput = []rating

type GetStockRatingsErrorKind int

const (
	_ GetStockRatingsErrorKind = iota
	getStockRatingsUnexpectedError
)

type GetStockRatingsError struct {
	kind GetStockRatingsErrorKind
	err  error
}

func (e GetStockRatingsError) Error() string {
	switch e.kind {
	case getStockRatingsUnexpectedError:
		return fmt.Sprintf("Unexpected error: %s", e.err.Error())
	default:
		return "Unknown error"
	}
}

func (e GetStockRatingsError) From(err error) GetStockRatingsError {
	e1 := e
	e1.err = err
	return e1
}
func (e GetStockRatingsError) Unwrap() error {
	return e.err
}

var (
	GetStockRatingsErrorUnexpectedError = GetStockRatingsError{kind: getStockRatingsUnexpectedError}
)

func (s *Service) GetStockRatings(input GetStockRatingsInput) (GetStockRatingsOutput, error) {
	res, err := s.repo.GetStockRatings(context.Background(), repository.GetStockRatingsParams{
		SortOrder:   input.sortOrder,
		SortBy:      input.sortBy,
		Offset:      input.offset,
		Limit:       input.limit,
		TickerLike:  input.tickerLike,
		CompanyLike: input.companyLike,
	})
	if err != nil {
		return nil, GetStockRatingsErrorUnexpectedError.From(err)
	}

	var out GetStockRatingsOutput
	for _, r := range res {
		out = append(out, rating{
			ticker:      r.Ticker,
			company:     r.Company,
			brokerage:   r.Brokerage,
			targetFrom:  r.TargetFrom,
			targetTo:    r.TargetTo,
			action:      string(r.Action),
			rawAction:   r.RawAction,
			ratingFrom:  string(r.RatingFrom),
			ratingTo:    string(r.RatingTo),
			at:          r.At,
			targetDelta: r.TargetDelta,
			score:       r.Score,
		})
	}
	return out, nil
}
