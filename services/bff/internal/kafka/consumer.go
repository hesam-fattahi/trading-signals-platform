package kafka

import (
	"context"
	"encoding/json"
	"log/slog"
	"time"

	"github.com/segmentio/kafka-go"
	"google.golang.org/protobuf/proto"

	tradingv1 "github.com/hesam-fattahi/trading-signals-platform/gen/go/trading/v1"
	"github.com/hesam-fattahi/trading-signals-platform/services/bff/internal/websocket"
)

// Message is the unified JSON structure sent to the frontend.
type Message struct {
	Topic string          `json:"topic"` // "tick", "indicator", or "signal"
	Data  json.RawMessage `json:"data"`
}

func Start(ctx context.Context, brokers, rawTopic, indTopic, sigTopic string, hub *websocket.Hub, logger *slog.Logger) {
	// Raw Ticks
	go consume(ctx, brokers, rawTopic, "bff-live-tick", logger, func(val []byte) {
		var bar tradingv1.MarketBarEvent
		if err := proto.Unmarshal(val, &bar); err == nil {
			payload, _ := json.Marshal(map[string]interface{}{
				"time":   bar.EventTimeUnixMs / 1000,
				"open":   bar.Open,
				"high":   bar.High,
				"low":    bar.Low,
				"close":  bar.Close,
				"volume": bar.Volume,
			})
			hub.Broadcast("tick", payload)
		}
	})

	// Indicators
	go consume(ctx, brokers, indTopic, "bff-live-ind", logger, func(val []byte) {
		var ind tradingv1.IndicatorEvent
		if err := proto.Unmarshal(val, &ind); err == nil {
			payload, _ := json.Marshal(map[string]interface{}{
				"time":      ind.EventTimeUnixMs / 1000,
				"bb_upper":  ind.BollingerUpper,
				"bb_middle": ind.BollingerMiddle,
				"bb_lower":  ind.BollingerLower,
				"rsi":       ind.Rsi,
				"adx":       ind.Adx,
				"di_plus":   ind.DiPlus,
				"di_minus":  ind.DiMinus,
			})
			hub.Broadcast("indicator", payload)
		}
	})

	// Signals
	go consume(ctx, brokers, sigTopic, "bff-live-sig", logger, func(val []byte) {
		var sig tradingv1.SignalEvent
		if err := proto.Unmarshal(val, &sig); err == nil {
			payload, _ := json.Marshal(map[string]interface{}{
				"time":   sig.EventTimeUnixMs / 1000,
				"action": sig.Action.String(),
				"price":  sig.Price,
			})
			hub.Broadcast("signal", payload)
		}
	})
}

func consume(ctx context.Context, brokers, topic, groupID string, logger *slog.Logger, handler func([]byte)) {
	r := kafka.NewReader(kafka.ReaderConfig{
		Brokers:     []string{brokers},
		Topic:       topic,
		GroupID:     groupID,
		StartOffset: kafka.LastOffset, // Only read new messages from now on
		MinBytes:    1,
		MaxBytes:    10e6,
		MaxWait:     500 * time.Millisecond,
	})
	defer r.Close()

	for {
		msg, err := r.ReadMessage(ctx)
		if err != nil {
			if ctx.Err() != nil {
				return
			}
			logger.Error("kafka read error", "topic", topic, "err", err)
			time.Sleep(time.Second)
			continue
		}
		handler(msg.Value)
	}
}
