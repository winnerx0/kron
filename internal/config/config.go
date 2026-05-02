package config

import (
	"fmt"
	"os"
)

type Config struct {
	DBHost        string
	DBUser        string
	DBPassword    string
	DBPort        string
	DBName        string
	EncryptionKey string
}

func Load() Config {

	cfg := Config{
		DBHost:        getEnv("DB_HOST"),
		DBUser:        getEnv("DB_USER"),
		DBPassword:    getEnv("DB_PASSWORD"),
		DBPort:        getEnv("DB_PORT"),
		DBName:        getEnv("DB_NAME"),
		EncryptionKey: getEnv("KRON_ENCRYPTION_KEY"),
	}

	return cfg
}

func getEnv(key string) (env string) {

	if env = os.Getenv(key); env == "" {
		panic(fmt.Errorf("%s not found in .env", key))
	}
	return
}
