package indicator;

import com.tictactec.ta.lib.Core;
import com.tictactec.ta.lib.MInteger;
import com.tictactec.ta.lib.RetCode;

public final class RsiCalculator {
  public static double compute(double[] closes, int period) {
    if (closes.length < period + 1) return Double.NaN;

    Core core = new Core();
    MInteger begin = new MInteger();
    MInteger length = new MInteger();
    double[] out = new double[closes.length - period];
    RetCode rc = core.rsi(0, closes.length - 1, closes, period, begin, length, out);
    if (rc != RetCode.Success || length.value == 0) return Double.NaN;
    return out[length.value - 1];
  }
}
