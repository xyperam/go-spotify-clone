package utils

import (
	"fmt"
	"log"

	"github.com/xyperam/go-spotify-clone/config"
	"github.com/xyperam/go-spotify-clone/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func ConnectDatabase() {
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%d sslmode=%s", config.Host, config.User, config.Password, config.DBName, config.Port, config.SSLMode)

	database, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})

	//auto migrate
	if err := database.AutoMigrate(&models.User{}); err != nil {
		log.Fatal("Failed auto migrate database:", err)
	}
	if err != nil {
		log.Fatal("Failed to connect database:", err)
	}
	DB = database
	fmt.Println("Connected to database")
}

func CheckConnection() error {
	sqlDB, err := DB.DB()
	if err != nil {
		log.Fatal("Failed to get DB instance:", err)
	}
	if err := sqlDB.Ping(); err != nil {
		log.Fatal("Failed to ping database:", err)
	}
	return nil
}
