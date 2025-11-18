package main

import (
	"backend/internal/features/stockratings"
	"backend/internal/repository"
	"backend/internal/routes"
	"backend/pkg/db"
	"log"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

var router = gin.Default()

func main() {
	// Read environment variables
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	// DEPENDENCY INJECTION ========================================================================
	db := db.Get()
	repo := repository.New(db)
	service := stockratings.NewService(repo)
	handler := stockratings.NewHandler(service)

	// Start the server
	router = gin.Default()
	routes.GetRoutes(router, handler)
	router.Run(":5000")

	// data, err := repo.GetStockRatings(
	// 	context.Background(),
	// 	repository.GetStockRatingsParams{
	// 		SortBy:      "score",
	// 		SortOrder:   "asc",
	// 		TickerLike:  "",
	// 		CompanyLike: "",
	// 		Offset:      0,
	// 		Limit:       10,
	// 	},
	// )
	// if err != nil {
	// 	log.Fatal("Error getting stock ratings", err)
	// }
	// log.Println(data)
}
