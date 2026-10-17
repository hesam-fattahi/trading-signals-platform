package indicator;

import java.util.ArrayDeque;
import java.util.Deque;

public final class SymbolState {
  private static final int MAX_BARS = 100;
  private final Deque<Double> closes = new ArrayDeque<>(MAX_BARS);
  private final Deque<Double> highs = new ArrayDeque<>(MAX_BARS);
  private final Deque<Double> lows = new ArrayDeque<>(MAX_BARS);
  private long lastEventTimeMs = 0;

  public void append(double close, double high, double low, long eventTimeMs) {
    closes.addLast(close);
    highs.addLast(high);
    lows.addLast(low);
    lastEventTimeMs = eventTimeMs;
    if (closes.size() > MAX_BARS) closes.removeFirst();
    if (highs.size() > MAX_BARS) highs.removeFirst();
    if (lows.size() > MAX_BARS) lows.removeFirst();
  }

  public double[] closesArray() {
    return closes.stream().mapToDouble(Double::doubleValue).toArray();
  }

  public double[] highsArray() {
    return highs.stream().mapToDouble(Double::doubleValue).toArray();
  }

  public double[] lowsArray() {
    return lows.stream().mapToDouble(Double::doubleValue).toArray();
  }

  public long lastEventTimeMs() {
    return lastEventTimeMs;
  }
}
