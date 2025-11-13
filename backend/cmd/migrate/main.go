package main

import (
	"backend/pkg/db"
	"log"
	"os"

	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	if len(os.Args) < 2 {
		log.Fatal("Usage: go run ./cmd/migrate [up|down]")
	}

	direction := os.Args[1]
	db.RunMigrations(direction)
}
