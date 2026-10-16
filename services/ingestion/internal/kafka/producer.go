package kafka

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"github.com/segmentio/kafka-go"
	"google.golang.org/protobuf/proto"

	tradingv1 "github.com/hesam-fattahi/trading-signals-platform/gen/go/trading/v1"
)

type Producer struct {
	writer *kafka.Writer
	topic  string
	logger *slog.Logger
}

func NewProducer(broker string, topic string, logger *slog.Logger) *Producer {
	writer := &kafka.Writer{
		Addr:         kafka.TCP(broker),
		Topic:        topic,
		Balancer:     &kafka.Hash{},
		BatchTimeout: 10 * time.Millisecond,
		RequiredAcks: kafka.RequireAll,
	}

	return &Producer{
		writer: writer,
		logger: logger,
	}
}

func (p *Producer) Produce(ctx context.Context, event *tradingv1.MarketBarEvent, idempotentKey string) error {
	payload, err := proto.Marshal(event)
	if err != nil {
		return fmt.Errorf("failed to marshal proto MarketBarEvent: %w", err)
	}

	msg := kafka.Message{
		Key:   []byte(idempotentKey),
		Value: payload,
		Time:  time.Now(),
	}

	if err := p.writer.WriteMessages(ctx, msg); err != nil {
		return fmt.Errorf("failed to write kafka message: %w", err)
	}

	// p.logger.Debug("produced market bar event to kafka", "symbol", event.Symbol, "key", idempotentKey)
	p.logger.Info("produced market bar event",
		"topic", p.topic,
		"key", idempotentKey,
		"symbol", event.Symbol,
		"close", event.Close,
	)
	return nil
}

func (p *Producer) Close() error {
	return p.writer.Close()
}
