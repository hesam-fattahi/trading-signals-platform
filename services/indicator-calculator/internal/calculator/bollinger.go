package calculator

import (
	"math"

	ta "github.com/TA-Lib/ta-lib-cgo"
)

func Bollinger(closes []float64, window int, k float64) (upper, middle, lower float64, ok bool) {
	if len(closes) < window {
		return 0, 0, 0, false
	}

	upperSlice, middleSlice, lowerSlice := ta.Bbands(closes, window, k, k, 0, nil, nil, nil)

	if len(upperSlice) == 0 || len(middleSlice) == 0 || len(lowerSlice) == 0 {
		return 0, 0, 0, false
	}

	upper = upperSlice[len(upperSlice)-1]
	middle = middleSlice[len(middleSlice)-1]
	lower = lowerSlice[len(lowerSlice)-1]

	if math.IsNaN(upper) || math.IsNaN(middle) || math.IsNaN(lower) {
		return 0, 0, 0, false
	}

	return upper, middle, lower, true
}
