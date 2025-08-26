package config

import (
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
)

type DatabaseConfig interface {
	Print() string
}

type PostgresConfig struct {
	Host         string
	Port         string
	User         string
	Password     string
	DatabaseName string
}

func (*PostgresConfig) Load() PostgresConfig {
	err := godotenv.Load("config.env")
	if err != nil {
		log.Fatalln("No .env file found or failed to load. Using system env.")
	}

	return PostgresConfig{
		Host:         getEnv("DB_HOST"),
		Port:         getEnv("DB_PORT"),
		User:         getEnv("DB_USER"),
		Password:     getEnv("DB_PASSWORD"),
		DatabaseName: getEnv("DB_NAME"),
	}
}

func getEnv(key string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}

	return ""
}

func (c *PostgresConfig) Print() string {
	return fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=UTC",
		c.Host, c.User, c.Password, c.DatabaseName, c.Port,
	)
}
