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
	logger *slog.Logger
}

func NewProducer(broker string, topic string, logger *slog.Logger) *Producer {
	writer := &kafka.Writer{
		Addr:         kafka.TCP(broker),
		Topic:        topic,
		Balancer:     &kafka.Hash{}, // To make sure messages with the same key go to the same partition
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

	p.logger.Debug("produced market bar event to kafka", "symbol", event.Symbol, "key", idempotentKey)
	return nil
}

func (p *Producer) Close() error {
	return p.writer.Close()
}
