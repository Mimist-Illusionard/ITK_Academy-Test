package config

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/joho/godotenv"
)

type DatabaseConfig interface {
	Print() string
}

type PostgresConfig struct {
	Host            string
	Port            string
	User            string
	Password        string
	DatabaseName    string
	MaxOpenConns    int
	MaxIdleConns    int
	ConnMaxLifetime time.Duration
	ConnMaxIdleTime time.Duration
}

func (*PostgresConfig) Load() PostgresConfig {
	err := godotenv.Load("config.env")
	if err != nil {
		log.Fatalln("No .env file found or failed to load. Using system env.")
	}

	return PostgresConfig{
		Host:            getEnv("DB_HOST"),
		Port:            getEnv("DB_PORT"),
		User:            getEnv("DB_USER"),
		Password:        getEnv("DB_PASSWORD"),
		DatabaseName:    getEnv("DB_NAME"),
		MaxOpenConns:    getEnvAsInt("DB_MAX_OPEN_CONNS"),
		MaxIdleConns:    getEnvAsInt("DB_MAX_IDLE_CONNS"),
		ConnMaxLifetime: time.Duration(getEnvAsInt("DB_CONN_MAX_LIFETIME")) * time.Second,
		ConnMaxIdleTime: time.Duration(getEnvAsInt("DB_CONN_MAX_IDLE_TIME")) * time.Second,
	}
}

func getEnv(key string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}

	return ""
}

func getEnvAsInt(key string) int {
	valueStr := getEnv(key)
	if value, err := strconv.Atoi(valueStr); err == nil {
		return value
	}
	return 0
}

func (c *PostgresConfig) Print() string {
	return fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=UTC",
		c.Host, c.User, c.Password, c.DatabaseName, c.Port,
	)
}
