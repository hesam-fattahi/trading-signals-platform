package mapper

import (
	"encoding/json"
	"fmt"

	tradingv1 "github.com/hesam-fattahi/trading-signals-platform/gen/go/trading/v1"
)

// RawMarketBar represents the incoming JSON structure from the WebSocket generator.
type RawMarketBar struct {
	Symbol          string  `json:"symbol"`
	EventTimeUnixMs int64   `json:"timestamp_ms"`
	Open            float64 `json:"open"`
	High            float64 `json:"high"`
	Low             float64 `json:"low"`
	Close           float64 `json:"close"`
	Volume          float64 `json:"volume"`
}

// An idempotent key for Kafka is added too.
func ParseAndMap(payload []byte) (*tradingv1.MarketBarEvent, string, error) {
	var raw RawMarketBar
	if err := json.Unmarshal(payload, &raw); err != nil {
		return nil, "", fmt.Errorf("failed to unmarshal JSON payload: %w", err)
	}

	// Basic validation
	if raw.Symbol == "" {
		return nil, "", fmt.Errorf("validation error: missing symbol")
	}
	if raw.EventTimeUnixMs == 0 {
		return nil, "", fmt.Errorf("validation error: invalid or missing event_time_unix_ms")
	}
	if raw.Close <= 0 {
		return nil, "", fmt.Errorf("validation error: invalid close price for %s", raw.Symbol)
	}

	// Generate a unique event ID or idempotent key
	eventId := fmt.Sprintf("%s-%d", raw.Symbol, raw.EventTimeUnixMs)

	event := &tradingv1.MarketBarEvent{
		EventId:         eventId,
		Symbol:          raw.Symbol,
		EventTimeUnixMs: raw.EventTimeUnixMs,
		Open:            raw.Open,
		High:            raw.High,
		Low:             raw.Low,
		Close:           raw.Close,
		Volume:          raw.Volume,
	}

	return event, eventId, nil
}
