package config

import (
	"os"
	"strconv"
	"sync"
)

type Config struct {
	Server  ServerConfig
	Logger  LoggerConfig
	DataDog DataDogConfig
}

type ServerConfig struct {
	Port string
}

type LoggerConfig struct {
	Level           string
	Encoding        string
	Service         string
	Environment     string
	AsyncBufferSize int
	BatchSize       int
}

type DataDogConfig struct {
	ServiceName string
	Environment string
}

var (
	instance *Config
	once     sync.Once
)

// GetInstance returns the singleton instance of Config
func GetInstance() *Config {
	once.Do(func() {
		instance = &Config{
			Server: ServerConfig{
				Port: getEnv("SERVER_PORT", "8080"),
			},
			Logger: LoggerConfig{
				Level:           getEnv("LOG_LEVEL", "info"),
				Encoding:        getEnv("LOG_ENCODING", "json"),
				Service:         getEnv("SERVICE_NAME", "go-gin-service"),
				Environment:     getEnv("ENVIRONMENT", "production"),
				AsyncBufferSize: getEnvAsInt("LOG_ASYNC_BUFFER_SIZE", 500),
				BatchSize:       getEnvAsInt("LOG_BATCH_SIZE", 20),
			},
			DataDog: DataDogConfig{
				ServiceName: getEnv("DD_SERVICE_NAME", "go-gin-service"),
				Environment: getEnv("DD_ENVIRONMENT", "production"),
			},
		}
	})
	return instance
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func getEnvAsInt(key string, defaultValue int) int {
	if value := os.Getenv(key); value != "" {
		if intValue, err := strconv.Atoi(value); err == nil {
			return intValue
		}
	}
	return defaultValue
}
