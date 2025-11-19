package stockratings

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type HandlerInterface interface {
	GetStockRatings(c *gin.Context)
}
type Handler struct {
	service ServiceInterface
}

func NewHandler(s *Service) *Handler {
	return &Handler{service: s}
}

type GetStockRatingsResponse struct {
	Ticker      string `json:"ticker"`
	Company     string `json:"company"`
	TargetFrom  string `json:"target_from"`
	TargetTo    string `json:"target_to"`
	Action      string `json:"action"`
	RatingFrom  string `json:"rating_from"`
	RatingTo    string `json:"rating_to"`
	At          string `json:"at"`
	TargetDelta string `json:"target_delta"`
	Score       int32  `json:"score"`
}

func (h *Handler) GetStockRatings(c *gin.Context) {
	// Validate parameters
	sortOrder := c.DefaultQuery("sort_order", "desc")
	sortBy := c.DefaultQuery("sort_by", "score")
	tickerLike := c.DefaultQuery("ticker_like", "")
	companyLike := c.DefaultQuery("company_like", "")
	offset, err := strconv.Atoi(c.DefaultQuery("offset", "0"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid offset"})
		return
	}
	limit, err := strconv.Atoi(c.DefaultQuery("limit", "10"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid limit"})
		return
	}

	// Call the service
	stockRatings, err := h.service.GetStockRatings(GetStockRatingsInput{
		sortOrder:   sortOrder,
		sortBy:      sortBy,
		offset:      int32(offset),
		limit:       int32(limit),
		tickerLike:  tickerLike,
		companyLike: companyLike,
	})
	fmt.Println(stockRatings)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Serialize the output

	resp := make([]GetStockRatingsResponse, len(stockRatings))
	for i, r := range stockRatings {
		resp[i] = GetStockRatingsResponse{
			Ticker:      r.ticker,
			Company:     r.company,
			TargetFrom:  r.targetFrom,
			TargetTo:    r.targetTo,
			Action:      string(r.action),
			RatingFrom:  string(r.ratingFrom),
			RatingTo:    string(r.ratingTo),
			At:          r.at.String(),
			TargetDelta: r.targetDelta,
			Score:       r.score,
		}
	}

	c.JSON(200, gin.H{
		"length":  len(resp),
		"ratings": resp,
	})
}

func AddStockRatingRoutes(rg *gin.RouterGroup, h HandlerInterface) {
	stockRatings := rg.Group("/stock_ratings")
	stockRatings.GET("/", h.GetStockRatings)
}
