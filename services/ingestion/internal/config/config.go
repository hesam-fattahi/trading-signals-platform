package config

import (
	"os"
)

type Config struct {
	WebSocketURL string
	KafkaBrokers string
	KafkaTopic   string
	LogLevel     string
}

// Load reads configuration from environment variables with fallback defaults.
func Load() *Config {
	return &Config{
		WebSocketURL: getEnv("WS_URL", "ws://localhost:8080/ws"),
		KafkaBrokers: getEnv("KAFKA_BROKERS", "localhost:9092"),
		KafkaTopic:   getEnv("KAFKA_TOPIC", "raw-ticks"),
		LogLevel:     getEnv("LOG_LEVEL", "info"),
	}
}

func getEnv(key, fallback string) string {
	if val, ok := os.LookupEnv(key); ok {
		return val
	}
	return fallback
}
