package db

import (
	"fmt"
	"log"

	"github.com/golang-migrate/migrate/v4"
	migratePostgres "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	_ "github.com/lib/pq"
)

var DB *gorm.DB

func InitDB(connAddr, migrationsDir string) {
	db, err := gorm.Open(postgres.Open(connAddr), &gorm.Config{})
	if err != nil {
		log.Fatalf("Failed to connect to the final database: %v", err)
	}

	err = runMigrations(db, migrationsDir)
	if err != nil {
		log.Fatalf("Subsequent migration failed: %v", err)
	}

	DB = db

	log.Println("Final database connected and all migrations applied successfully!")
}

func runMigrations(gormDB *gorm.DB, migrationsDir string) error {
	db, err := gormDB.DB()
	if err != nil {
		return fmt.Errorf("could not obtain sql.db from gorm.db: %w", err)
	}
	driver, err := migratePostgres.WithInstance(db, &migratePostgres.Config{})
	if err != nil {
		return fmt.Errorf("could not create migration driver: %w", err)
	}

	m, err := migrate.NewWithDatabaseInstance(
		fmt.Sprintf("file://%s", migrationsDir),
		"postgres", driver)
	if err != nil {
		return fmt.Errorf("could not create migrate instance: %w", err)
	}

	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		return fmt.Errorf("could not run up migrations: %w", err)
	}

	log.Println("Migrations applied successfully!")
	return nil
}
