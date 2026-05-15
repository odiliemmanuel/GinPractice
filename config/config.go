package config

import (
	"fmt"
	"os"
)

type Config struct {
	DBHost     string
	DBPort     string
	DBUser     string
	DBName     string
	DBPassword string
	Port       string	
}



func Load() Config {
	return Config{
		DBHost:     getEnv("DB_HOST", "localhost"),
		DBPort:     getEnv("DB_PORT", "5432"),
		DBUser:     getEnv("DB_USER", "gotask_user"),
		DBName:     getEnv("DB_NAME", "gotask"),
		DBPassword: getEnv("DB_PASSWORD", "password"),
		Port:       getEnv("PORT", "8080"),
	}
}



func (c Config) DSN() string {
	return fmt.Sprintf(
		"host=%s port=%s user=%s dbname=%s password=%s sslmode=disable",
		c.DBHost, c.DBPort, c.DBUser, c.DBName, c.DBPassword,
	)
}



func getEnv(key, fallback string) string {
	if value := os.Getenv(key); value != ""{
		return value
	}

	return fallback
}