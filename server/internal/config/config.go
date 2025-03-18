package config

import (
	"os"
	"strconv"
	"time"
)

type Config struct {
	ServerPort        string
	Difficulty        int
	ConnectionTimeout time.Duration
	QuotesFilePath    string
}

func LoadConfig() *Config {
	serverPort := getEnvString("SERVER_PORT", "8080")
	difficulty := getEnvInt("DIFFICULTY", 5)
	connectionTimeout := getEnvDuration("CONNECTION_TIMEOUT", 10*time.Second)
	quotesFilePath := getEnvString("QUOTES_FILE_PATH", "quotes.txt")

	return &Config{
		ServerPort:        serverPort,
		Difficulty:        difficulty,
		ConnectionTimeout: connectionTimeout,
		QuotesFilePath:    quotesFilePath,
	}
}

func getEnvString(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}

func getEnvInt(key string, defaultValue int) int {
	if value, exists := os.LookupEnv(key); exists {
		if intValue, err := strconv.Atoi(value); err == nil {
			return intValue
		}
	}
	return defaultValue
}

func getEnvDuration(key string, defaultValue time.Duration) time.Duration {
	if value, exists := os.LookupEnv(key); exists {
		if durationValue, err := time.ParseDuration(value); err == nil {
			return durationValue
		}
	}
	return defaultValue
}
