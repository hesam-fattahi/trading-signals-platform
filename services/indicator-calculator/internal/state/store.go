package state

import (
	"sync"

	"google.golang.org/protobuf/proto"

	tradingv1 "github.com/hesam-fattahi/trading-signals-platform/gen/go/trading/v1"
)

type Store struct {
	mu      sync.RWMutex
	symbols map[string]*tradingv1.SymbolState
	maxBars int
}

func New(maxBars int) *Store {
	return &Store{
		symbols: make(map[string]*tradingv1.SymbolState),
		maxBars: maxBars,
	}
}

func (s *Store) Update(symbol string, bar *tradingv1.MarketBarEvent) *tradingv1.SymbolState {
	s.mu.Lock()
	defer s.mu.Unlock()

	state, exists := s.symbols[symbol]
	if !exists {
		state = &tradingv1.SymbolState{
			Symbol: symbol,
		}
		s.symbols[symbol] = state
	}

	state.RecentCloses = append(state.RecentCloses, bar.Close)
	state.RecentHighs = append(state.RecentHighs, bar.High)
	state.RecentLows = append(state.RecentLows, bar.Low)
	state.LastEventTimeUnixMs = bar.EventTimeUnixMs

	if len(state.RecentCloses) > s.maxBars {
		state.RecentCloses = state.RecentCloses[len(state.RecentCloses)-s.maxBars:]
	}
	if len(state.RecentHighs) > s.maxBars {
		state.RecentHighs = state.RecentHighs[len(state.RecentHighs)-s.maxBars:]
	}
	if len(state.RecentLows) > s.maxBars {
		state.RecentLows = state.RecentLows[len(state.RecentLows)-s.maxBars:]
	}

	return proto.Clone(state).(*tradingv1.SymbolState)
}

func (s *Store) Replace(snapshots []*tradingv1.SymbolState) {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.symbols = make(map[string]*tradingv1.SymbolState, len(snapshots))
	for _, snap := range snapshots {
		if snap != nil && snap.Symbol != "" {
			s.symbols[snap.Symbol] = snap
		}
	}
}
