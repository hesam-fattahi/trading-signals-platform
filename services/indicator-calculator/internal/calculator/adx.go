package calculator

import (
	"math"

	ta "github.com/TA-Lib/ta-lib-cgo"
)

func Adx(highs, lows, closes []float64, window int) (float64, bool) {
	if len(highs) != len(lows) || len(lows) != len(closes) {
		return 0, false
	}
	lb := ta.AdxLookback(window)
	if len(highs) < lb+1 {
		return 0, false
	}

	out := ta.Adx(highs, lows, closes, window, nil)
	if len(out) == 0 {
		return 0, false
	}

	val := out[len(out)-1]
	if math.IsNaN(val) {
		return 0, false
	}
	return val, true
}
