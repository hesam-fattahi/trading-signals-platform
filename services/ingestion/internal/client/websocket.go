package client

import (
	"context"
	"log/slog"
	"math"
	"math/rand"
	"time"

	"github.com/gorilla/websocket"
)

// MessageHandler defines the callback signature for processing received raw message payloads.
type MessageHandler func(ctx context.Context, data []byte) error

type WSClient struct {
	url     string
	handler MessageHandler
	logger  *slog.Logger
}

func NewWSClient(url string, handler MessageHandler, logger *slog.Logger) *WSClient {
	return &WSClient{
		url:     url,
		handler: handler,
		logger:  logger,
	}
}

// Start connects to the WebSocket stream and runs an infinite reconnection loop with exponential backoff & jitter.
func (c *WSClient) Start(ctx context.Context) error {
	initialBackoff := 1 * time.Second
	maxBackoff := 30 * time.Second
	backoff := initialBackoff

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
		}

		c.logger.Info("connecting to market data websocket", "url", c.url)
		conn, _, err := websocket.DefaultDialer.DialContext(ctx, c.url, nil)
		if err != nil {
			c.logger.Error("failed to connect to websocket, retrying...", "error", err, "backoff_sec", backoff.Seconds())
			sleepWithContext(ctx, backoff)

			// Exponential increase with jitter
			backoff = time.Duration(math.Min(float64(maxBackoff), float64(backoff)*1.5))
			jitter := time.Duration(rand.Int63n(int64(500 * time.Millisecond)))
			backoff += jitter
			continue
		}

		// Reset backoff on successful connection
		backoff = initialBackoff
		c.logger.Info("successfully established websocket connection")

		errChan := make(chan error, 1)
		go func() {
			defer conn.Close()
			for {
				_, message, err := conn.ReadMessage()
				if err != nil {
					errChan <- err
					return
				}
				if err := c.handler(ctx, message); err != nil {
					c.logger.Error("failed to process incoming message", "error", err)
				}
			}
		}()

		select {
		case <-ctx.Done():
			conn.Close()
			return ctx.Err()
		case err := <-errChan:
			c.logger.Warn("websocket connection dropped", "error", err)
			// Loop continues and reconnects automatically
		}
	}
}

func sleepWithContext(ctx context.Context, d time.Duration) {
	select {
	case <-ctx.Done():
	case <-time.After(d):
	}
}
