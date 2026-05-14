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