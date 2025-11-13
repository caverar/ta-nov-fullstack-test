package db

import (
	"database/sql"
	"log"
	"os"
	"sync"

	_ "github.com/jackc/pgx/v5/stdlib"
)

var (
	once     sync.Once
	instance *sql.DB
)

func Get() *sql.DB {
	once.Do(func() {
		// Get the database connection string from .env
		dsn := os.Getenv("DATABASE_URL")
		if dsn == "" {
			log.Fatal("DATABASE_URL not set")
		}

		// Connect to the database
		var err error
		instance, err = sql.Open("pgx", dsn)
		if err != nil {
			log.Fatalf("failed to connect to database: %v", err)
		}
		// Check if the connection is alive
		if err := instance.Ping(); err != nil {
			log.Fatalf("ping error: %v", err)
		}
	})
	return instance
}
