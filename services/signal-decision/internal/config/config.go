package config

import (
	"os"
)

type Config struct {
	KafkaBrokers  string
	InputTopic    string
	OutputTopic   string
	ConsumerGroup string
	LogLevel      string
}

func Load() *Config {
	return &Config{
		KafkaBrokers:  getEnv("KAFKA_BROKERS", "localhost:9092"),
		InputTopic:    getEnv("INPUT_TOPIC", "indicator-events"),
		OutputTopic:   getEnv("OUTPUT_TOPIC", "signal-events"),
		ConsumerGroup: getEnv("CONSUMER_GROUP", "signal-decision-service"),
		LogLevel:      getEnv("LOG_LEVEL", "info"),
	}
}

func getEnv(key, fallback string) string {
	if val, ok := os.LookupEnv(key); ok {
		return val
	}
	return fallback
}
