package main

import "github.com/bensiauu/financial-assistance-scheme/pkg/db"

func main() {
	connAddr := "host=localhost user=govtech password=password123 dbname=financial_assistance sslmode=disable"

	migrationsDir := "pkg/db/migrations"

	// Initialize the database and run migrations
	db.InitDB(connAddr, migrationsDir)

	// Continue with the rest of your application logic...
	// Example: Start your Gin server here
	// r := gin.Default()
	// r.Run()
}
