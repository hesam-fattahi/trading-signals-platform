package calculator

import (
	"testing"
)

func TestCalculators(t *testing.T) {
	closes := make([]float64, 50)
	highs := make([]float64, 50)
	lows := make([]float64, 50)

	for i := 0; i < 50; i++ {
		base := 10.0 + float64(i)*0.1
		closes[i] = base
		highs[i] = base + 0.2
		lows[i] = base - 0.2
	}

	// Test Bollinger Bands
	upper, middle, lower, bbOk := Bollinger(closes, 20, 2.0)
	if !bbOk || !(upper > middle && middle > lower) {
		t.Errorf("Bollinger calculation failed: ok=%v", bbOk)
	}

	// Test RSI
	rsiVal, rsiOk := Rsi(closes, 14)
	if !rsiOk || rsiVal < 0 || rsiVal > 100 {
		t.Errorf("RSI calculation failed: ok=%v, val=%f", rsiOk, rsiVal)
	}

	// Test ADX
	adxVal, adxOk := Adx(highs, lows, closes, 14)
	if !adxOk || adxVal < 0 || adxVal > 100 {
		t.Errorf("ADX calculation failed: ok=%v, val=%f", adxOk, adxVal)
	}
}
