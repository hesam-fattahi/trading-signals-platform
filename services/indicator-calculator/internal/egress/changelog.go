package egress

import (
	"context"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"time"

	"github.com/segmentio/kafka-go"
	"google.golang.org/protobuf/proto"

	tradingv1 "github.com/hesam-fattahi/trading-signals-platform/gen/go/trading/v1"
)

type Changelog struct {
	writer *kafka.Writer
}

func NewChangelog(brokers, topic string) *Changelog {
	return &Changelog{
		writer: &kafka.Writer{
			Addr:         kafka.TCP(brokers),
			Topic:        topic,
			Balancer:     &kafka.Hash{},
			RequiredAcks: kafka.RequireAll,
		},
	}
}

func (c *Changelog) Write(ctx context.Context, snap *tradingv1.SymbolState) error {
	data, err := proto.Marshal(snap)
	if err != nil {
		return fmt.Errorf("marshal symbol state: %w", err)
	}
	return c.writer.WriteMessages(ctx, kafka.Message{
		Key:   []byte(snap.Symbol),
		Value: data,
	})
}

func (c *Changelog) Close() error {
	return c.writer.Close()
}

func Restore(ctx context.Context, brokers, topic string, logger *slog.Logger) ([]*tradingv1.SymbolState, error) {
	reader := kafka.NewReader(kafka.ReaderConfig{
		Brokers:     []string{brokers},
		Topic:       topic,
		GroupID:     "",
		StartOffset: kafka.FirstOffset,
		MinBytes:    1,
		MaxBytes:    10e6,
	})
	defer reader.Close()

	stateMap := make(map[string]*tradingv1.SymbolState)
	logger.Info("restoring state from changelog topic", "topic", topic)

	for {
		fetchCtx, cancel := context.WithTimeout(ctx, 1*time.Second)
		msg, err := reader.FetchMessage(fetchCtx)
		cancel()

		if err != nil {
			if errors.Is(err, context.DeadlineExceeded) {
				logger.Info("changelog restoration caught up to live tail", "restored_symbols", len(stateMap))
				break
			}
			if errors.Is(err, context.Canceled) || errors.Is(err, io.EOF) {
				break
			}
			break
		}

		if len(msg.Value) == 0 {
			if len(msg.Key) > 0 {
				delete(stateMap, string(msg.Key))
			}
			continue
		}

		var snap tradingv1.SymbolState
		if err := proto.Unmarshal(msg.Value, &snap); err != nil {
			logger.Warn("failed to unmarshal changelog record", "error", err)
			continue
		}

		if snap.Symbol != "" {
			stateMap[snap.Symbol] = &snap
		}
	}

	snapshots := make([]*tradingv1.SymbolState, 0, len(stateMap))
	for _, snap := range stateMap {
		snapshots = append(snapshots, snap)
	}
	return snapshots, nil
}
