package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	MongoURI     string
	JWTSecretKey string
	Port         string
	BaseURL      string
	RedisURL     string
}

var AppConfig *Config

func LoadConfig() {
	if err := godotenv.Load(); err != nil {
		log.Println("Warning: .env file not found, using system environment variables.")
	}

	AppConfig = &Config{
		MongoURI:     os.Getenv("MONGO_URI"),
		JWTSecretKey: os.Getenv("JWT_SECRET_KEY"),
		Port:         os.Getenv("PORT"),
		BaseURL:      os.Getenv("BASE_URL"),
		RedisURL:     os.Getenv("REDIS_URL"),
	}
}

func GetJWTSecret() []byte {
	return []byte(AppConfig.JWTSecretKey)
}
