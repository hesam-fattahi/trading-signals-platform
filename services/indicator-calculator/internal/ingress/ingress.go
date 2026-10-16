package ingress

import (
	"context"
	"errors"
	"fmt"
	"log/slog"

	"github.com/segmentio/kafka-go"
	"google.golang.org/protobuf/proto"

	tradingv1 "github.com/hesam-fattahi/trading-signals-platform/gen/go/trading/v1"
)

type Handler func(ctx context.Context, bar *tradingv1.MarketBarEvent) error

type Ingress struct {
	reader  *kafka.Reader
	topic   string
	handler Handler
	logger  *slog.Logger
}

func New(brokers, topic, groupID string, handler Handler, logger *slog.Logger) *Ingress {
	return &Ingress{
		reader: kafka.NewReader(kafka.ReaderConfig{
			Brokers:     []string{brokers},
			Topic:       topic,
			GroupID:     groupID,
			StartOffset: kafka.FirstOffset,
			MinBytes:    1,
			MaxBytes:    10e6,
		}),
		topic:   topic,
		handler: handler,
		logger:  logger,
	}
}

func (i *Ingress) Run(ctx context.Context) error {
	for {
		msg, err := i.reader.FetchMessage(ctx)
		if err != nil {
			if errors.Is(err, context.Canceled) {
				return ctx.Err()
			}
			return fmt.Errorf("fetch message: %w", err)
		}

		var bar tradingv1.MarketBarEvent
		if err := proto.Unmarshal(msg.Value, &bar); err != nil {
			i.logger.Warn("failed to unmarshal MarketBarEvent", "error", err)
			_ = i.reader.CommitMessages(ctx, msg)
			continue
		}

		i.logger.Info("consumed bar",
			"topic", i.topic,
			"key", string(msg.Key),
			"partition", msg.Partition,
			"offset", msg.Offset,
			"symbol", bar.Symbol,
			"close", bar.Close,
		)

		if err := i.handler(ctx, &bar); err != nil {
			i.logger.Error("handler failed", "error", err, "offset", msg.Offset)
			continue
		}

		if err := i.reader.CommitMessages(ctx, msg); err != nil {
			i.logger.Error("commit failed", "error", err)
		}
	}
}

func (i *Ingress) Close() error {
	return i.reader.Close()
}
