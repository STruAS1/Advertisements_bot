package db

import (
	"fmt"
	"log"
	"tgbotBARAHOLKA/config"
	"tgbotBARAHOLKA/db/models"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func Connect(cfg *config.Config) {
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%d sslmode=disable",
		cfg.Database.Host, cfg.Database.User, cfg.Database.Password, cfg.Database.Name, cfg.Database.Port)
	var err error

	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	log.Println("Connected to database.")

	if err := DB.AutoMigrate(&models.User{}, &models.Advertisement{}, &models.AdvertisementInputs{}, &models.AdvertisementType{}, &models.Cities{}); err != nil {
		log.Fatalf("Error creating tables: %v", err)
	} else {
		log.Println("Tables created successfully.")
	}
}
