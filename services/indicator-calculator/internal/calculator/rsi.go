package calculator

import (
	"math"

	ta "github.com/TA-Lib/ta-lib-cgo"
)

func Rsi(closes []float64, window int) (float64, bool) {
	if len(closes) < window+1 {
		return 0, false
	}

	out := ta.Rsi(closes, window, nil)
	if len(out) == 0 {
		return 0, false
	}

	val := out[len(out)-1]
	if math.IsNaN(val) {
		return 0, false
	}

	return val, true
}
