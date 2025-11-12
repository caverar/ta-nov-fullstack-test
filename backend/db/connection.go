package db

import (
	"log"
	"os"
	"sync"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/jmoiron/sqlx"
)

var (
	once     sync.Once
	instance *sqlx.DB
)

func Get() *sqlx.DB {
	once.Do(func() {
		// Get the database connection string from .env
		dsn := os.Getenv("DATABASE_URL")
		if dsn == "" {
			log.Fatal("DATABASE_URL not set")
		}
		log.Println("Connecting to database ...")

		// Connect to the database
		var err error
		instance, err = sqlx.Connect("pgx", dsn)
		if err != nil {
			log.Fatalf("failed to connect to database: %v", err)
		}

		log.Println("✅ Database connection established")
	})
	return instance
}
