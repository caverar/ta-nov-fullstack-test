package main

import (
	"backend/internal/features/stockratings"
	"backend/internal/repository"
	"backend/pkg/db"
	"log"

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
	initializer := stockratings.NewLoaderService(repo)

	// INITIALIZE THE DATA =========================================================================
	err = initializer.InitData()
	if err != nil {
		log.Fatal("Error initializing stock data from API", err)
	}
}
