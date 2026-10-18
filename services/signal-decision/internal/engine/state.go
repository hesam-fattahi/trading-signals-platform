package engine

type PositionState string

const (
	StateFlat  PositionState = "FLAT"
	StateLong  PositionState = "LONG"
	StateShort PositionState = "SHORT"
)

type SymbolContext struct {
	Symbol  string
	State   PositionState
	PrevRSI float64
	HasPrev bool
}

func ApplySignal(current PositionState, action string) PositionState {
	switch action {
	case "SIGNAL_ACTION_LONG_ENTRY":
		return StateLong
	case "SIGNAL_ACTION_SHORT_ENTRY":
		return StateShort
	case "SIGNAL_ACTION_LONG_EXIT", "SIGNAL_ACTION_SHORT_EXIT":
		return StateFlat
	default:
		return current
	}
}
