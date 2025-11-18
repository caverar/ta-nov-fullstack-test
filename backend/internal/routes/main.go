package routes

import (
	stockratings "backend/internal/features/stockratings"
	"net/http"

	"github.com/gin-gonic/gin"
)

func GetRoutes(rg *gin.Engine, h *stockratings.Handler) {
	ping := rg.Group("/ping")
	ping.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, "pong")
	})

	v1 := rg.Group("/v1")
	stockratings.AddStockRatingRoutes(v1, h)

}
