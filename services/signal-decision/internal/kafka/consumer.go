package kafka

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"time"

	"github.com/segmentio/kafka-go"
	"google.golang.org/protobuf/proto"

	tradingv1 "github.com/hesam-fattahi/trading-signals-platform/gen/go/trading/v1"
	"github.com/hesam-fattahi/trading-signals-platform/services/signal-decision/internal/engine"
)

type Handler func(ctx context.Context, event *tradingv1.IndicatorEvent) error

type Ingress struct {
	reader *kafka.Reader
	topic  string
	logger *slog.Logger
}

func NewConsumer(brokers, topic, groupID string, logger *slog.Logger) *Ingress {
	return &Ingress{
		reader: kafka.NewReader(kafka.ReaderConfig{
			Brokers:     []string{brokers},
			Topic:       topic,
			GroupID:     groupID,
			StartOffset: kafka.FirstOffset,
			MinBytes:    1,
			MaxBytes:    10e6,
			Dialer: &kafka.Dialer{
				Timeout:   10 * time.Second,
				DualStack: true,
			},
		}),
		topic:  topic,
		logger: logger,
	}
}

func (i *Ingress) Run(ctx context.Context, handler Handler) error {
	for {
		fetchCtx, cancel := context.WithTimeout(ctx, 15*time.Second)
		msg, err := i.reader.FetchMessage(fetchCtx)
		cancel()

		if err != nil {
			if errors.Is(err, context.DeadlineExceeded) {
				// No message in 15s. Loop and try again.
				continue
			}
			if errors.Is(err, context.Canceled) {
				return ctx.Err()
			}
			return fmt.Errorf("fetch message: %w", err)
		}

		var event tradingv1.IndicatorEvent
		if err := proto.Unmarshal(msg.Value, &event); err != nil {
			i.logger.Warn("failed to unmarshal IndicatorEvent", "error", err)
			_ = i.reader.CommitMessages(ctx, msg)
			continue
		}

		i.logger.Info("consumed indicator",
			"topic", i.topic,
			"key", string(msg.Key),
			"partition", msg.Partition,
			"offset", msg.Offset,
			"symbol", event.Symbol,
			"close", event.Close,
			"rsi", event.Rsi,
			"adx", event.Adx,
		)

		if err := handler(ctx, &event); err != nil {
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

// This is a one-shot read: it opens a reader with no group, reads to the
// end of the topic, and returns. It does NOT continue consuming after that.
// The purpose is to bootstrap the StrategyEngine at pod restart.
func RestoreState(ctx context.Context, brokers, topic string, logger *slog.Logger) (map[string]*engine.SymbolContext, error) {
	reader := kafka.NewReader(kafka.ReaderConfig{
		Brokers:     []string{brokers},
		Topic:       topic,
		GroupID:     "", // no group; read once and exit
		StartOffset: kafka.FirstOffset,
		MinBytes:    1,
		MaxBytes:    10e6,
		Dialer: &kafka.Dialer{
			Timeout:   10 * time.Second,
			DualStack: true,
		},
	})
	defer reader.Close()

	states := make(map[string]*engine.SymbolContext)
	logger.Info("restoring position state from signal-events", "topic", topic)

	for {
		fetchCtx, cancel := context.WithTimeout(ctx, 2*time.Second)
		msg, err := reader.FetchMessage(fetchCtx)
		cancel()

		if err != nil {
			if errors.Is(err, context.DeadlineExceeded) {
				logger.Info("signal-events restore caught up to live tail",
					"symbols", len(states))
				break
			}
			if errors.Is(err, context.Canceled) {
				return states, ctx.Err()
			}
			return nil, fmt.Errorf("restore fetch: %w", err)
		}

		var sig tradingv1.SignalEvent
		if err := proto.Unmarshal(msg.Value, &sig); err != nil {
			logger.Warn("failed to unmarshal SignalEvent during restore", "error", err)
			continue
		}

		ctx2, ok := states[sig.Symbol]
		if !ok {
			ctx2 = &engine.SymbolContext{
				Symbol: sig.Symbol,
				State:  engine.StateFlat,
			}
			states[sig.Symbol] = ctx2
		}

		actionName := tradingv1.SignalAction_name[int32(sig.Action)]
		ctx2.State = engine.ApplySignal(ctx2.State, actionName)
	}

	return states, nil
}
