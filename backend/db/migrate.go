package db

import (
	"fmt"
	"log"
	"os"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/cockroachdb"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

// CockroachDB config for migrations
type cockroachDBConfig struct {
	migrationsTable string
	lockingDisabled bool
}

func RunMigrations(direction string) {
	dsn := os.Getenv("DATABASE_URL")
	if dsn == "" {
		log.Fatal("DATABASE_URL not set")
	}

	db := Get().DB
	// CockroachDB compatible config - disable locking since pg_advisory_lock is not supported
	driver, err := cockroachdb.WithInstance(db, &cockroachdb.Config{})
	if err != nil {
		log.Fatalf("failed to create migration driver: %v", err)
	}

	m, err := migrate.NewWithDatabaseInstance(
		"file://db/migrations",
		"postgres",
		driver,
	)
	if err != nil {
		log.Fatalf("failed to create migrate instance: %v", err)
	}

	switch direction {
	case "up":
		if err := m.Up(); err != nil && err != migrate.ErrNoChange {
			log.Fatalf("migration up failed: %v", err)
		}
		fmt.Println("✅ Migrations applied successfully")
	case "down":
		if err := m.Down(); err != nil {
			log.Fatalf("migration down failed: %v", err)
		}
		fmt.Println("✅ Migrations reverted successfully")
	default:
		log.Fatalf("unknown direction: %s (use 'up' or 'down')", direction)
	}
}
