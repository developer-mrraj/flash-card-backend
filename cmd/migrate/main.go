package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/joho/godotenv"
)

func main() {
	// Define flags
	up := flag.Bool("up", false, "Run all pending migrations")
	down := flag.Bool("down", false, "Rollback the last migration")
	steps := flag.Int("steps", 0, "Migrate up/down by N steps (positive=up, negative=down)")
	version := flag.Bool("version", false, "Show current migration version")
	force := flag.Int("force", -1, "Force set migration version (use to fix dirty state)")
	flag.Parse()

	// Load .env file
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using system environment variables")
	}

	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		log.Fatal("DATABASE_URL environment variable is not set")
	}

	// Create migrate instance
	m, err := migrate.New("file://migrations", dbURL)
	if err != nil {
		log.Fatalf("Failed to create migrate instance: %v", err)
	}
	defer m.Close()

	// Show current version if no flags provided or --version flag
	if *version || (!*up && !*down && *steps == 0 && *force == -1) {
		v, dirty, err := m.Version()
		if err != nil {
			fmt.Printf("Current version: none (fresh database)\n")
		} else {
			fmt.Printf("Current version: %d (dirty: %v)\n", v, dirty)
		}
		if !*version {
			fmt.Println("\nUsage:")
			fmt.Println("  go run ./cmd/migrate -up        Run all pending migrations")
			fmt.Println("  go run ./cmd/migrate -down       Rollback last migration")
			fmt.Println("  go run ./cmd/migrate -steps N    Migrate N steps (positive=up, negative=down)")
			fmt.Println("  go run ./cmd/migrate -version    Show current migration version")
			fmt.Println("  go run ./cmd/migrate -force V    Force set version V (fix dirty state)")
		}
		return
	}

	// Force version (to fix dirty database)
	if *force >= 0 {
		if err := m.Force(*force); err != nil {
			log.Fatalf("Failed to force version: %v", err)
		}
		fmt.Printf("✅ Forced migration version to %d\n", *force)
		return
	}

	// Run migrations by steps
	if *steps != 0 {
		if err := m.Steps(*steps); err != nil && err != migrate.ErrNoChange {
			log.Fatalf("Failed to migrate steps: %v", err)
		}
		fmt.Printf("✅ Migrated %d steps successfully\n", *steps)
		return
	}

	// Migrate up
	if *up {
		if err := m.Up(); err != nil && err != migrate.ErrNoChange {
			log.Fatalf("Failed to migrate up: %v", err)
		}
		v, _, _ := m.Version()
		fmt.Printf("✅ All migrations applied successfully (version: %d)\n", v)
		return
	}

	// Migrate down
	if *down {
		if err := m.Steps(-1); err != nil && err != migrate.ErrNoChange {
			log.Fatalf("Failed to migrate down: %v", err)
		}
		v, _, _ := m.Version()
		fmt.Printf("✅ Rolled back one migration (version: %d)\n", v)
		return
	}
}
