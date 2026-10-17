package indicator;

import com.google.protobuf.InvalidProtocolBufferException;
import org.apache.kafka.common.serialization.Serde;
import org.apache.kafka.common.serialization.Serdes;
import trading.v1.SymbolStateOuterClass.SymbolState;

public final class SymbolStateSerde {

  public static Serde<indicator.SymbolState> instance() {
    return Serdes.serdeFrom(
        (topic, state) -> {
          SymbolState.Builder b =
              SymbolState.newBuilder().setLastEventTimeUnixMs(state.lastEventTimeMs());
          for (double d : state.closesArray()) b.addRecentCloses(d);
          for (double d : state.highsArray()) b.addRecentHighs(d);
          for (double d : state.lowsArray()) b.addRecentLows(d);
          return b.build().toByteArray();
        },
        (topic, bytes) -> {
          try {
            SymbolState p = SymbolState.parseFrom(bytes);
            indicator.SymbolState s = new indicator.SymbolState();
            for (int i = 0; i < p.getRecentClosesCount(); i++) {
              s.append(
                  p.getRecentCloses(i),
                  p.getRecentHighs(i),
                  p.getRecentLows(i),
                  p.getLastEventTimeUnixMs());
            }
            return s;
          } catch (InvalidProtocolBufferException e) {
            throw new RuntimeException("Failed to deserialize SymbolState", e);
          }
        });
  }
}
