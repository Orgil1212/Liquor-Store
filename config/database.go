package config

import (
	"liquor-store/models"
	"log"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func ConnectDatabase() {
	dsn := "host=localhost user=postgres password=1234 dbname=liquor_store port=5432 sslmode=disable"
	database, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})

	if err != nil {
		log.Fatal("Өгөгдлийн сантай холбогдож чадсангүй:", err)
	}

	// AutoMigrate ашиглан хүснэгтийг шинэчлэх
	database.AutoMigrate(&models.Product{})

	DB = database
}
