package main

import (
	"log"

	"github.com/bensiauu/financial-assistance-scheme/internal/router"
	"github.com/bensiauu/financial-assistance-scheme/pkg/db"
)

func main() {
	connAddr := "host=localhost user=govtech password=password123 dbname=financial_assistance sslmode=disable"

	migrationsDir := "pkg/db/migrations"

	// Initialize the database and run migrations
	db.InitDB(connAddr, migrationsDir)

	r := router.SetupRouter()

	if err := r.Run(":8080"); err != nil {
		log.Fatalf("Failed to run server: %v", err)
	}
}
