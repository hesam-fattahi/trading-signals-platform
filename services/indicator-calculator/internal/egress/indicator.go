package egress

import (
	"context"

	"github.com/segmentio/kafka-go"
	"google.golang.org/protobuf/proto"

	tradingv1 "github.com/hesam-fattahi/trading-signals-platform/gen/go/trading/v1"
)

type Indicator struct {
	writer *kafka.Writer
}

func NewIndicator(brokers, topic string) *Indicator {
	return &Indicator{
		writer: &kafka.Writer{
			Addr:         kafka.TCP(brokers),
			Topic:        topic,
			Balancer:     &kafka.Hash{},
			RequiredAcks: kafka.RequireAll,
		},
	}
}

func (p *Indicator) Send(ctx context.Context, event *tradingv1.IndicatorEvent) error {
	data, err := proto.Marshal(event)
	if err != nil {
		return err
	}
	return p.writer.WriteMessages(ctx, kafka.Message{
		Key:   []byte(event.Symbol),
		Value: data,
	})
}

func (p *Indicator) Close() error {
	return p.writer.Close()
}
