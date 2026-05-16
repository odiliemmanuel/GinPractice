package db

import (
	"log"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gotask/config"
	"gotask/models"
)

var DB *gorm.DB

func Connect(configs config.Config){
	var err error

	DB, err = gorm.Open(postgres.Open(configs.DSN()), &gorm.Config{})

	if err != nil {
		log.Fatal("Failed to connect to database: ", err)
	}

	if err := DB.AutoMigrate(&models.Task{}, &models.User{}); err != nil {
		log.Fatal("Migration failed:", err)
	}

	log.Println("Database connected and migrated. ")
}
