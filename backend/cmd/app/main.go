package main

import (
	"backend/internal/features/stockratings"
	"backend/internal/repository"
	"backend/internal/routes"
	"backend/pkg/db"
	"log"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

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
	router := gin.Default()
	router.Use(cors.Default()) // All origins allowed
	// router.Use(cors.New(cors.Config{
	// 	// AllowOrigins: []string{"http://localhost:5173"},
	// 	AllowOrigins: []string{"*"},
	// 	// AllowAllOrigins:  true,
	// 	AllowMethods: []string{"PUT", "PATCH", "GET", "POST", "DELETE"},
	// 	// AllowHeaders: []string{
	// 	// 	"Origin",
	// 	// 	"Content-Type",
	// 	// 	"Accept",
	// 	// 	"Authorization",
	// 	// 	"Cache-Control",
	// 	// 	"Pragma",
	// 	// },
	// 	AllowHeaders: []string{"*"},
	// 	// ExposeHeaders:    []string{"Content-Length"},
	// 	AllowCredentials: false,
	// 	// MaxAge:           12 * time.Hour,
	// }))
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
