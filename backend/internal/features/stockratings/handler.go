package stockratings

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	service *Service
}

func NewHandler(s *Service) *Handler {
	return &Handler{service: s}
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
	res, err := h.service.GetStockRatings(GetStockRatingsInput{
		SortOrder:   sortOrder,
		SortBy:      sortBy,
		Offset:      int32(offset),
		Limit:       int32(limit),
		TickerLike:  tickerLike,
		CompanyLike: companyLike,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, gin.H{"data": res})
}

func AddStockRatingRoutes(rg *gin.RouterGroup, h *Handler) {
	stockRatings := rg.Group("/stock_ratings")
	stockRatings.GET("/", h.GetStockRatings)
}
