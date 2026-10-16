package main

import (
	"context"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"github.com/hesam-fattahi/trading-signals-platform/services/indicator-calculator/internal/config"
	"github.com/hesam-fattahi/trading-signals-platform/services/indicator-calculator/internal/egress"
	"github.com/hesam-fattahi/trading-signals-platform/services/indicator-calculator/internal/ingress"
	"github.com/hesam-fattahi/trading-signals-platform/services/indicator-calculator/internal/processor"
	"github.com/hesam-fattahi/trading-signals-platform/services/indicator-calculator/internal/state"
)

func main() {
	cfg := config.Load()
	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: parseLogLevel(cfg.LogLevel),
	}))

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		sig := <-sigCh
		logger.Info("received shutdown signal", "signal", sig)
		cancel()
	}()

	logger.Info("restoring state from changelog", "topic", cfg.StateChangelogTopic)
	snapshots, err := egress.Restore(ctx, cfg.Brokers, cfg.StateChangelogTopic, logger)
	if err != nil {
		logger.Error("state restore failed", "error", err)
		os.Exit(1)
	}
	logger.Info("state restoration complete", "symbols_restored", len(snapshots))

	store := state.New(cfg.MaxBars)
	store.Replace(snapshots)

	changelog := egress.NewChangelog(cfg.Brokers, cfg.StateChangelogTopic)
	defer changelog.Close()

	indicator := egress.NewIndicator(cfg.Brokers, cfg.IndicatorTopic)
	defer indicator.Close()

	proc := processor.New(store, changelog, indicator, logger, cfg)

	ing := ingress.New(cfg.Brokers, cfg.RawTicksTopic, cfg.GroupID, proc.Handle, logger)
	defer ing.Close()

	logger.Info("starting indicator-calculator ingress loop", "topic", cfg.RawTicksTopic)
	if err := ing.Run(ctx); err != nil && err != context.Canceled {
		logger.Error("ingress stopped with error", "error", err)
		os.Exit(1)
	}

	logger.Info("indicator-calculator gracefully shut down")
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
