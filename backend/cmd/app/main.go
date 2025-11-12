package main

import (
	"backend/db"
	"log"

	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	_ = db.Get()
	log.Println("🚀 app started")
}
