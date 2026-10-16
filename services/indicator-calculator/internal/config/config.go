package config

import (
	"os"
	"strconv"
)

type Config struct {
	Brokers             string
	RawTicksTopic       string
	IndicatorTopic      string
	StateChangelogTopic string
	GroupID             string
	LogLevel            string
	MaxBars             int
	BollingerWindow     int
	BollingerK          float64
	RsiWindow           int
	AdxWindow           int
}

func Load() *Config {
	return &Config{
		Brokers:             getEnv("KAFKA_BROKERS", "trading-events-kafka-bootstrap.platform.svc.cluster.local:9092"),
		RawTicksTopic:       getEnv("RAW_TICKS_TOPIC", "raw-ticks"),
		IndicatorTopic:      getEnv("INDICATOR_TOPIC", "indicator-events"),
		StateChangelogTopic: getEnv("STATE_CHANGELOG_TOPIC", "indicator-state-changelog"),
		GroupID:             getEnv("KAFKA_GROUP_ID", "indicator-service-group"),
		LogLevel:            getEnv("LOG_LEVEL", "info"),
		MaxBars:             getEnvAsInt("MAX_BARS", 100),
		BollingerWindow:     getEnvAsInt("BOLLINGER_WINDOW", 20),
		BollingerK:          getEnvAsFloat("BOLLINGER_K", 2.0),
		RsiWindow:           getEnvAsInt("RSI_WINDOW", 14),
		AdxWindow:           getEnvAsInt("ADX_WINDOW", 14),
	}
}

func getEnv(key, fallback string) string {
	if val, ok := os.LookupEnv(key); ok {
		return val
	}
	return fallback
}

func getEnvAsInt(key string, fallback int) int {
	if val, ok := os.LookupEnv(key); ok {
		if parsed, err := strconv.Atoi(val); err == nil {
			return parsed
		}
	}
	return fallback
}

func getEnvAsFloat(key string, fallback float64) float64 {
	if val, ok := os.LookupEnv(key); ok {
		if parsed, err := strconv.ParseFloat(val, 64); err == nil {
			return parsed
		}
	}
	return fallback
}
