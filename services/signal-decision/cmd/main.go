package main

import (
	"context"
	"errors"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
	"time"

	tradingv1 "github.com/hesam-fattahi/trading-signals-platform/gen/go/trading/v1"
	"github.com/hesam-fattahi/trading-signals-platform/services/signal-decision/internal/config"
	"github.com/hesam-fattahi/trading-signals-platform/services/signal-decision/internal/engine"
	"github.com/hesam-fattahi/trading-signals-platform/services/signal-decision/internal/kafka"
)

func main() {
	cfg := config.Load()
	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: parseLogLevel(cfg.LogLevel),
	}))

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Graceful shutdown on SIGINT/SIGTERM.
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		sig := <-sigCh
		logger.Info("received shutdown signal", "signal", sig)
		cancel()
	}()

	// --- 1. Restore position state from signal-events ---
	logger.Info("restoring position state", "topic", cfg.OutputTopic)
	snapshot, err := kafka.RestoreState(ctx, cfg.KafkaBrokers, cfg.OutputTopic, logger)
	if err != nil {
		logger.Error("state restore failed", "error", err)
		os.Exit(1)
	}
	logger.Info("state restoration complete", "symbols", len(snapshot))

	// --- 2. Build the strategy engine and hand it the restored state ---
	strat := engine.NewStrategyEngine()
	strat.Restore(snapshot)

	// --- 3. Producer for signal-events ---
	producer := kafka.NewProducer(cfg.KafkaBrokers, cfg.OutputTopic, logger)
	defer func() {
		if err := producer.Close(); err != nil {
			logger.Error("failed to close producer", "error", err)
		}
	}()

	// --- 4. Consumer for indicator-events ---
	ing := kafka.NewConsumer(cfg.KafkaBrokers, cfg.InputTopic, cfg.ConsumerGroup, logger)
	defer func() {
		if err := ing.Close(); err != nil {
			logger.Error("failed to close consumer", "error", err)
		}
	}()

	// --- 5. Main loop ---
	logger.Info("starting signal-decision ingress loop", "topic", cfg.InputTopic)
	backoff := 5 * time.Second
	for {
		err := ing.Run(ctx, func(ctx context.Context, ind *tradingv1.IndicatorEvent) error {
			sig, produced := strat.Evaluate(ind)
			if !produced {
				return nil
			}
			return producer.Produce(ctx, sig, sig.Symbol)
		})

		if errors.Is(err, context.Canceled) {
			break
		}
		logger.Error("ingress stopped, retrying",
			"error", err,
			"backoff_sec", backoff.Seconds(),
		)
		select {
		case <-ctx.Done():
			return
		case <-time.After(backoff):
		}
		if backoff < 60*time.Second {
			backoff *= 2
		} else {
			backoff = 60 * time.Second
		}
	}

	logger.Info("signal-decision gracefully shut down")
}

func parseLogLevel(level string) slog.Level {
	switch level {
	case "debug":
		return slog.LevelDebug
	case "warn":
		return slog.LevelWarn
	case "error":
		return slog.LevelError
	default:
		return slog.LevelInfo
	}
}
