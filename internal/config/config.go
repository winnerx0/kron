package config

import (
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	DBHost     string
	DBUser     string
	DBPassword string
	DBPort     string
	DBName     string
}

func Load() Config {
	
	err := godotenv.Load()
	
	if err != nil {
		log.Println("Failed to load .env")
	}

	cfg := Config{
		DBHost:     getEnv("DB_HOST"),
		DBUser:     getEnv("DB_USER"),
		DBPassword: getEnv("DB_PASSWORD"),
		DBPort:     getEnv("DB_PORT"),
		DBName:     getEnv("DB_NAME"),
	}

	return cfg
}

func getEnv(key string) (env string) {

	if env = os.Getenv(key); env == "" {
		panic(fmt.Errorf("%s not found in .env", key))
	}
	return
}
