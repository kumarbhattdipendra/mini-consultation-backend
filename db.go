package main

import (
	"backend/models"
	"fmt"
	"log"
	"os"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var DB *gorm.DB

func InitDB() {
	dbURI := os.Getenv("DB_URI")
	if dbURI == "" {
		log.Fatal("DB_URI environment variable is required")
	}

	var err error
	DB, err = gorm.Open(postgres.Open(dbURI), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	sqlDB, err := DB.DB()
	if err != nil {
		log.Fatal("Could not get generic DB object:", err)
	}
	if err := sqlDB.Ping(); err != nil {
		log.Fatal("Database not reachable:", err)
	}

	fmt.Println("Connected to Database successfully")

	// Run migrations
	err = DB.AutoMigrate(&models.User{}, &models.Guide{}, &models.Booking{})
	if err != nil {
		log.Fatal("Migration failed:", err)
	}

	fmt.Println(" Database migrated successfully")
}

func AutoMigrate(db *gorm.DB) error {
	return db.AutoMigrate(&models.User{}, &models.Guide{}, &models.Booking{})
}

func CloseDB() {
	if DB != nil {
		sqlDB, err := DB.DB()
		if err != nil {
			log.Println("Could not get generic DB object:", err)
			return
		}
		if err := sqlDB.Close(); err != nil {
			log.Println("Error closing database connection:", err)
		} else {
			fmt.Println("Database connection closed successfully")
		}
	}
}
