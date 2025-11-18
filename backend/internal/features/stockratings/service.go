package stockratings

import (
	"backend/internal/repository"
	"context"
	"fmt"
)

// SERVICE =========================================================================================

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
	SortOrder   string
	SortBy      string
	Offset      int32
	Limit       int32
	TickerLike  string
	CompanyLike string
}
type GetStockRatingsOutput = []repository.GetStockRatingsRow
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
		SortOrder:   input.SortOrder,
		SortBy:      input.SortBy,
		Offset:      input.Offset,
		Limit:       input.Limit,
		TickerLike:  input.TickerLike,
		CompanyLike: input.CompanyLike,
	})
	if err != nil {
		return nil, GetStockRatingsErrorUnexpectedError.From(err)
	}
	return res, nil
}
