package main

import (
	"backend/internal/repository"
	"backend/pkg/db"
	"context"
	"log"

	"github.com/joho/godotenv"
)

func main() {
	// Read environment variables
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	db := db.Get()
	repo := repository.New(db)

	data, err := repo.GetStockRatings(
		context.Background(),
		repository.GetStockRatingsParams{
			SortBy:      "score",
			SortOrder:   "asc",
			TickerLike:  "",
			CompanyLike: "",
			Offset:      0,
			Limit:       10,
		},
	)
	if err != nil {
		log.Fatal("Error getting stock ratings", err)
	}
	log.Println(data)
}
