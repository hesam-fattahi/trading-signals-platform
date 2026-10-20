package config

import (
	"os"
)

type Config struct {
	KafkaBrokers    string
	RawTicksTopic   string
	IndicatorsTopic string
	SignalsTopic    string
}

func Load() *Config {
	return &Config{
		KafkaBrokers:    getEnv("KAFKA_BROKERS", "localhost:9092"),
		RawTicksTopic:   getEnv("RAW_TICKS_TOPIC", "raw-ticks"),
		IndicatorsTopic: getEnv("INDICATORS_TOPIC", "indicator-events"),
		SignalsTopic:    getEnv("SIGNALS_TOPIC", "signal-events"),
	}
}

func getEnv(key, fallback string) string {
	if val, ok := os.LookupEnv(key); ok {
		return val
	}
	return fallback
}
