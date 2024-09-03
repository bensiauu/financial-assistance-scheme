package main

import (
	"log"
	"os"

	"github.com/bensiauu/financial-assistance-scheme/internal/middleware"
	"github.com/bensiauu/financial-assistance-scheme/internal/router"
	"github.com/bensiauu/financial-assistance-scheme/pkg/db"
)

func main() {
	connAddr := "host=localhost user=govtech password=password123 dbname=financial_assistance sslmode=disable"

	migrationsDir := "pkg/db/migrations"

	secretKey := os.Getenv("JWT_SECRET_KEY")
	if secretKey == "" {
		var err error
		secretKey, err = middleware.GenerateSecureKey(32)
		if err != nil {
			log.Fatalf("Failed to generate secure key: %v", err)
		}
		log.Println("Generated Secret Key:", secretKey)
	}
	middleware.JWTSecret = []byte(secretKey)

	// Initialize the database and run migrations
	db.InitDB(connAddr, migrationsDir)

	r := router.SetupRouter()

	if err := r.Run(":8080"); err != nil {
		log.Fatalf("Failed to run server: %v", err)
	}
}
