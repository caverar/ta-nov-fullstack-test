package db

import (
	"context"
	"database/sql"
	"log"
	"os"
	"sync"

	"github.com/jackc/pgx/v5"
	_ "github.com/jackc/pgx/v5/stdlib"
)

var (
	once         sync.Once
	instance     *pgx.Conn
	raw_instance *sql.DB
)

func Get() *pgx.Conn {
	once.Do(func() {
		// Get the database connection string from .env
		dsn := os.Getenv("DATABASE_URL")
		if dsn == "" {
			log.Fatal("DATABASE_URL not set")
		}

		// Connect to the database
		var err error
		instance, err = pgx.Connect(context.Background(), dsn)
		if err != nil {
			log.Fatalf("failed to connect to database: %v", err)
		}
		// Check if the connection is alive
		if err := instance.Ping(context.Background()); err != nil {
			log.Fatalf("ping error: %v", err)
		}
	})
	return instance
}

func GetStd() *sql.DB {
	once.Do(func() {
		// Get the database connection string from .env
		dsn := os.Getenv("DATABASE_URL")
		if dsn == "" {
			log.Fatal("DATABASE_URL not set")
		}

		// Connect to the database
		var err error
		raw_instance, err = sql.Open("pgx", dsn)
		if err != nil {
			log.Fatalf("failed to connect to database: %v", err)
		}
		// Check if the connection is alive
		if err := raw_instance.Ping(); err != nil {
			log.Fatalf("ping error: %v", err)
		}
	})
	return raw_instance
}
