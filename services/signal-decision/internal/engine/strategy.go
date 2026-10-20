package engine

import (
	"fmt"

	tradingv1 "github.com/hesam-fattahi/trading-signals-platform/gen/go/trading/v1"
)

type StrategyEngine struct {
	states map[string]*SymbolContext
}

func NewStrategyEngine() *StrategyEngine {
	return &StrategyEngine{
		states: make(map[string]*SymbolContext),
	}
}

// Restore replaces the in-memory state with a snapshot from signal-events.
// Called once at startup, after RestoreState has replayed the topic.
func (s *StrategyEngine) Restore(snapshot map[string]*SymbolContext) {
	if snapshot == nil {
		snapshot = make(map[string]*SymbolContext)
	}
	s.states = snapshot
}

func (s *StrategyEngine) Evaluate(event *tradingv1.IndicatorEvent) (*tradingv1.SignalEvent, bool) {
	symbol := event.Symbol
	ctx, exists := s.states[symbol]
	if !exists {
		ctx = &SymbolContext{
			Symbol: symbol,
			State:  StateFlat,
		}
		s.states[symbol] = ctx
	}

	currentRSI := event.Rsi

	action := tradingv1.SignalAction_SIGNAL_ACTION_HOLD
	actionTriggered := false

	// ==========================================
	// LONG EXIT
	// ==========================================
	switch ctx.State {
	case StateLong:
		if event.Adx < 20 || event.DiMinus > event.DiPlus ||
			(event.Rsi > 70 && event.Close >= event.BollingerUpper) {
			ctx.State = StateFlat
			action = tradingv1.SignalAction_SIGNAL_ACTION_LONG_EXIT
			actionTriggered = true
		}
	// ==========================================
	// SHORT EXIT
	// ==========================================
	case StateShort:
		if event.Adx < 20 || event.DiPlus > event.DiMinus ||
			(event.Rsi < 30 && event.Close <= event.BollingerLower) {
			ctx.State = StateFlat
			action = tradingv1.SignalAction_SIGNAL_ACTION_SHORT_EXIT
			actionTriggered = true
		}
	}

	if ctx.State == StateFlat && !actionTriggered && event.Adx > 25 {
		// ==========================================
		// LONG ENTRY
		// ==========================================
		isUptrend := event.DiPlus > event.DiMinus
		isPullbackZone := event.Close <= event.BollingerMiddle
		rsiBullish := currentRSI > 50

		if isUptrend && isPullbackZone && rsiBullish {
			ctx.State = StateLong
			action = tradingv1.SignalAction_SIGNAL_ACTION_LONG_ENTRY
			actionTriggered = true
		} else {
			// ==========================================
			// SHORT ENTRY
			// ==========================================
			isDowntrend := event.DiMinus > event.DiPlus
			isRallyZone := event.Close >= event.BollingerMiddle
			rsiBearish := currentRSI < 50

			if isDowntrend && isRallyZone && rsiBearish {
				ctx.State = StateShort
				action = tradingv1.SignalAction_SIGNAL_ACTION_SHORT_ENTRY
				actionTriggered = true
			}
		}
	}

	if !actionTriggered {
		return nil, false
	}

	return s.createSignal(event, action), true
}

func (s *StrategyEngine) createSignal(event *tradingv1.IndicatorEvent, action tradingv1.SignalAction) *tradingv1.SignalEvent {
	return &tradingv1.SignalEvent{
		EventId:         fmt.Sprintf("sig-%s-%d", event.Symbol, event.EventTimeUnixMs),
		Symbol:          event.Symbol,
		EventTimeUnixMs: event.EventTimeUnixMs,
		Action:          action,
		Price:           event.Close,
		Strategy:        "ADX_RSI_Bollinger_v1",
	}
}
