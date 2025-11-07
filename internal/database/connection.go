package database

import (
	"log"
	"os"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func ConnectGORM(connStr string) (*gorm.DB, error) {
	db, err := gorm.Open(postgres.Open(connStr), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	sqlDB, err := db.DB()
	if err != nil {
		return nil, err
	}

	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(100)
	sqlDB.SetConnMaxLifetime(time.Hour)

	err = applySQLMigration(db, "migrations/001_init.sql")
	if err != nil {
		return nil, err
	}

	log.Println("Successfully connected to PostgreSQL")
	return db, nil
}

func applySQLMigration(db *gorm.DB, migrationPath string) error {
	sqlContent, err := os.ReadFile(migrationPath)
	if err != nil {
		return err
	}

	result := db.Exec(string(sqlContent))
	if result.Error != nil {
		return result.Error
	}

	log.Printf("Applied migration: %s", migrationPath)
	return nil
}
