package main

import (
	"context"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	web "github.com/hesam-fattahi/trading-signals-platform/services/bff"
	"github.com/hesam-fattahi/trading-signals-platform/services/bff/internal/config"
	"github.com/hesam-fattahi/trading-signals-platform/services/bff/internal/kafka"
	"github.com/hesam-fattahi/trading-signals-platform/services/bff/internal/websocket"
)

func main() {
	cfg := config.Load()
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	hub := websocket.NewHub()

	// Start Kafka to WebSocket pipeline
	kafka.Start(ctx, cfg.KafkaBrokers, cfg.RawTicksTopic, cfg.IndicatorsTopic, cfg.SignalsTopic, hub, logger)

	// HTTP server for WebSocket and static files only
	mux := http.NewServeMux()
	mux.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		websocket.ServeWS(hub, w, r)
	})
	web.RegisterStaticFiles(mux)

	server := &http.Server{Addr: ":8080", Handler: mux}
	go func() {
		logger.Info("BFF listening on :8080")
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Error("server error", "err", err)
			os.Exit(1)
		}
	}()

	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)
	<-sig

	cancel()
	_ = server.Shutdown(context.Background())
}
