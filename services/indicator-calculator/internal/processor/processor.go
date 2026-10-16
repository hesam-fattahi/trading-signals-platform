package processor

import (
	"context"
	"fmt"
	"log/slog"

	tradingv1 "github.com/hesam-fattahi/trading-signals-platform/gen/go/trading/v1"
	"github.com/hesam-fattahi/trading-signals-platform/services/indicator-calculator/internal/calculator"
	"github.com/hesam-fattahi/trading-signals-platform/services/indicator-calculator/internal/config"
	"github.com/hesam-fattahi/trading-signals-platform/services/indicator-calculator/internal/egress"
	"github.com/hesam-fattahi/trading-signals-platform/services/indicator-calculator/internal/state"
)

type Processor struct {
	state     *state.Store
	changelog *egress.Changelog
	indicator *egress.Indicator
	logger    *slog.Logger
	cfg       *config.Config
}

func New(store *state.Store, changelog *egress.Changelog, indicator *egress.Indicator, logger *slog.Logger, cfg *config.Config) *Processor {
	return &Processor{
		state:     store,
		changelog: changelog,
		indicator: indicator,
		logger:    logger,
		cfg:       cfg,
	}
}

func (p *Processor) Handle(ctx context.Context, bar *tradingv1.MarketBarEvent) error {
	// 1. Update in-memory state & get snapshot copy
	snap := p.state.Update(bar.Symbol, bar)

	// 2. Persist state snapshot to changelog topic
	if err := p.changelog.Write(ctx, snap); err != nil {
		return fmt.Errorf("changelog write: %w", err)
	}

	p.logger.Info("wrote changelog",
		"symbol", bar.Symbol,
		"bars_in_buffer", len(snap.RecentCloses),
	)

	// 3. Cold-start check across all required indicator lookbacks
	closes := snap.RecentCloses
	highs := snap.RecentHighs
	lows := snap.RecentLows

	requiredBars := p.cfg.BollingerWindow
	if p.cfg.RsiWindow+1 > requiredBars {
		requiredBars = p.cfg.RsiWindow + 1
	}
	if p.cfg.AdxWindow*2 > requiredBars {
		requiredBars = p.cfg.AdxWindow * 2
	}
	if len(closes) < requiredBars {
		p.logger.Debug("warm-up in progress",
			"symbol", bar.Symbol,
			"have", len(closes),
			"need", requiredBars,
		)
		return nil
	}

	// 4. Compute all indicators
	upper, middle, lower, bbOk := calculator.Bollinger(closes, p.cfg.BollingerWindow, p.cfg.BollingerK)
	rsiVal, rsiOk := calculator.Rsi(closes, p.cfg.RsiWindow)
	adxVal, adxOk := calculator.Adx(highs, lows, closes, p.cfg.AdxWindow)

	if !bbOk || !rsiOk || !adxOk {
		p.logger.Warn("indicator computation failed despite sufficient bars",
			"symbol", bar.Symbol,
			"bars", len(closes),
			"bb_ok", bbOk,
			"rsi_ok", rsiOk,
			"adx_ok", adxOk,
		)
		return nil
	}

	// 5. Emit complete IndicatorEvent with Bollinger, RSI, and ADX populated
	event := &tradingv1.IndicatorEvent{
		EventId:         "IND-" + bar.EventId,
		Symbol:          bar.Symbol,
		EventTimeUnixMs: bar.EventTimeUnixMs,
		Close:           bar.Close,
		BollingerUpper:  upper,
		BollingerMiddle: middle,
		BollingerLower:  lower,
		Rsi:             rsiVal,
		Adx:             adxVal,
	}

	if err := p.indicator.Send(ctx, event); err != nil {
		return fmt.Errorf("indicator send: %w", err)
	}

	p.logger.Info("produced indicator",
		"symbol", bar.Symbol,
		"event_id", event.EventId,
		"close", event.Close,
		"bb_upper", upper,
		"bb_middle", middle,
		"bb_lower", lower,
		"rsi", rsiVal,
		"adx", adxVal,
	)

	return nil
}
