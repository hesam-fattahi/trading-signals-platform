package indicator;

import com.tictactec.ta.lib.Core;
import com.tictactec.ta.lib.MInteger;
import com.tictactec.ta.lib.RetCode;

public final class AdxCalculator {
  public static double[] compute(double[] highs, double[] lows, double[] closes, int period) {
    if (highs.length != lows.length || lows.length != closes.length) return null;
    if (highs.length < period * 2 + 1) return null;

    Core core = new Core();

    int adxLookback = core.adxLookback(period);
    int diLookback = core.plusDILookback(period);

    int adxSize = highs.length - adxLookback;
    int diSize = highs.length - diLookback;

    if (adxSize <= 0 || diSize <= 0) return null;

    MInteger adxBegin = new MInteger();
    MInteger adxLength = new MInteger();
    double[] adx = new double[adxSize];

    MInteger plusBegin = new MInteger();
    MInteger plusLength = new MInteger();
    double[] plusDI = new double[diSize];

    MInteger minusBegin = new MInteger();
    MInteger minusLength = new MInteger();
    double[] minusDI = new double[diSize];

    RetCode rcAdx =
        core.adx(0, highs.length - 1, highs, lows, closes, period, adxBegin, adxLength, adx);
    if (rcAdx != RetCode.Success || adxLength.value == 0) return null;

    RetCode rcPlus =
        core.plusDI(
            0, highs.length - 1, highs, lows, closes, period, plusBegin, plusLength, plusDI);
    if (rcPlus != RetCode.Success || plusLength.value == 0) return null;

    RetCode rcMinus =
        core.minusDI(
            0, highs.length - 1, highs, lows, closes, period, minusBegin, minusLength, minusDI);
    if (rcMinus != RetCode.Success || minusLength.value == 0) return null;

    return new double[] {
      adx[adxLength.value - 1], plusDI[plusLength.value - 1], minusDI[minusLength.value - 1]
    };
  }
}
