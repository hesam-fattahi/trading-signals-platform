package kafka

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/segmentio/kafka-go"
	"google.golang.org/protobuf/proto"

	tradingv1 "github.com/hesam-fattahi/trading-signals-platform/gen/go/trading/v1"
)

type Producer struct {
	writer *kafka.Writer
	topic  string
	logger *slog.Logger
}

func NewProducer(brokers, topic string, logger *slog.Logger) *Producer {
	return &Producer{
		writer: &kafka.Writer{
			Addr:         kafka.TCP(brokers),
			Topic:        topic,
			Balancer:     &kafka.Hash{},
			RequiredAcks: kafka.RequireAll,
		},
		topic:  topic,
		logger: logger,
	}
}

func (p *Producer) Produce(ctx context.Context, event *tradingv1.SignalEvent, key string) error {
	payload, err := proto.Marshal(event)
	if err != nil {
		return fmt.Errorf("marshal signal event: %w", err)
	}

	if err := p.writer.WriteMessages(ctx, kafka.Message{
		Key:   []byte(key),
		Value: payload,
	}); err != nil {
		return fmt.Errorf("write signal event: %w", err)
	}

	p.logger.Info("produced signal",
		"topic", p.topic,
		"key", key,
		"symbol", event.Symbol,
		"action", event.Action.String(),
		"price", event.Price,
		"strategy", event.Strategy,
	)
	return nil
}

func (p *Producer) Close() error {
	return p.writer.Close()
}
