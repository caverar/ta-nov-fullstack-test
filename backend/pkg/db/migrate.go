package db

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/cockroachdb"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

func RunMigrations(command string) {
	dsn := os.Getenv("DATABASE_URL")
	if dsn == "" {
		log.Fatal("DATABASE_URL not set")
	}

	db := GetStd()
	driver, err := cockroachdb.WithInstance(db, &cockroachdb.Config{})
	if err != nil {
		log.Fatalf("failed to create migration driver: %v", err)
	}

	m, err := migrate.NewWithDatabaseInstance(
		"file://pkg/db/migrations",
		"postgres",
		driver,
	)
	if err != nil {
		log.Fatalf("failed to create migrate instance: %v", err)
	}
	commands := strings.Fields(command)
	if len(commands) == 0 {
		log.Fatal("no migration command provided")
	}

	switch commands[0] {
	case "up":
		if err := m.Steps(1); err != nil && err != migrate.ErrNoChange {
			log.Fatalf("migration up failed: %v", err)
		}
		fmt.Println("✅ 1 Migration applied successfully")
	case "up-last":
		if err := m.Up(); err != nil {
			log.Fatalf("migration up-last failed: %v", err)
		}
		fmt.Println("✅ All Mmgrations applied successfully")
	case "down":
		if err := m.Steps(-1); err != nil {
			log.Fatalf("migration down failed: %v", err)
		}
		fmt.Println("✅ 1 Migration reverted successfully")
	case "down-first":
		if err := m.Down(); err != nil {
			log.Fatalf("migration down-first failed: %v", err)
		}
		fmt.Println("✅ All migrations reverted successfully")
	case "force":
		if len(commands) < 2 {
			log.Fatalf("force requires a version argument")
		}
		version, err := strconv.Atoi(commands[1])
		if err != nil {
			log.Fatalf("invalid force version: %v", err)
		}
		if err := m.Force(version); err != nil {
			log.Fatalf("migration force failed: %v", err)
		}
		fmt.Println("✅ Migrations forced successfully")

	default:
		log.Fatalf("unknown direction: %s (use 'up' or 'down')", command)
	}
}
