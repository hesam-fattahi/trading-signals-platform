package main

import (
	"context"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"github.com/hesam-fattahi/trading-signals-platform/services/ingestion/internal/client"
	"github.com/hesam-fattahi/trading-signals-platform/services/ingestion/internal/config"

	"github.com/hesam-fattahi/trading-signals-platform/services/ingestion/internal/kafka"
	"github.com/hesam-fattahi/trading-signals-platform/services/ingestion/internal/mapper"
)

func main() {
	cfg := config.Load()

	var level slog.Level
	switch cfg.LogLevel {
	case "debug":
		level = slog.LevelDebug
	case "warn":
		level = slog.LevelWarn
	case "error":
		level = slog.LevelError
	default:
		level = slog.LevelInfo
	}

	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: level}))
	logger.Info("starting ingestion service", "ws_url", cfg.WebSocketURL, "kafka_topic", cfg.KafkaTopic)

	producer := kafka.NewProducer(cfg.KafkaBrokers, cfg.KafkaTopic, logger)
	defer func() {
		if err := producer.Close(); err != nil {
			logger.Error("failed to close kafka producer", "error", err)
		}
	}()

	// Message handling pipeline: JSON bytes -> Protobuf Mapper -> Kafka Produce
	messageHandler := func(ctx context.Context, data []byte) error {
		event, idempotentKey, err := mapper.ParseAndMap(data)
		if err != nil {
			return err
		}

		return producer.Produce(ctx, event, idempotentKey)
	}

	// Setup context with graceful shutdown listeners
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		sig := <-sigChan
		logger.Info("received shutdown signal", "signal", sig)
		cancel()
	}()

	// Start WebSocket Client
	wsClient := client.NewWSClient(cfg.WebSocketURL, messageHandler, logger)
	if err := wsClient.Start(ctx); err != nil && err != context.Canceled {
		logger.Error("websocket client stopped with error", "error", err)
		os.Exit(1)
	}

	logger.Info("ingestion service gracefully shut down")
}
