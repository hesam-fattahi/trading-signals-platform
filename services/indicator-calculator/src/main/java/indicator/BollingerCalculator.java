package indicator;

import com.tictactec.ta.lib.Core;
import com.tictactec.ta.lib.MAType;
import com.tictactec.ta.lib.MInteger;
import com.tictactec.ta.lib.RetCode;

public final class BollingerCalculator {
  public static double[] compute(double[] closes, int period, double k) {
    if (closes.length < period) return null;

    Core core = new Core();
    MInteger begin = new MInteger();
    MInteger length = new MInteger();
    double[] upper = new double[closes.length - period + 1];
    double[] middle = new double[closes.length - period + 1];
    double[] lower = new double[closes.length - period + 1];
    RetCode rc =
        core.bbands(
            0,
            closes.length - 1,
            closes,
            period,
            k,
            k,
            MAType.Sma,
            begin,
            length,
            upper,
            middle,
            lower);
    if (rc != RetCode.Success || length.value == 0) return null;
    return new double[] {
      upper[length.value - 1], middle[length.value - 1], lower[length.value - 1]
    };
  }
}
