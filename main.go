package main

import (
	"gotask/routes"
	"log"
	"gotask/config"
	"gotask/db"
	"gotask/middleware"

	"github.com/joho/godotenv"
	"github.com/gin-gonic/gin"
)

func main() {

	err := godotenv.Load()

	if err != nil {
		log.Fatal("Error loading .env file")
	}

	cfg := config.Load()

	db.Connect(cfg)

	r := gin.Default()

	r.Use(middleware.Logger())
	r.Use(gin.Recovery())

	routes.SetupRoutes(r)

	r.Run(":" + cfg.Port)

}
