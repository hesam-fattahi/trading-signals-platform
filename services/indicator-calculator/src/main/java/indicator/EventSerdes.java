package indicator;

import com.google.protobuf.InvalidProtocolBufferException;
import org.apache.kafka.common.serialization.Serde;
import org.apache.kafka.common.serialization.Serdes;
import trading.v1.Indicators.IndicatorEvent;
import trading.v1.MarketData.MarketBarEvent;

public final class EventSerdes {

  private EventSerdes() {}

  public static Serde<MarketBarEvent> marketBarSerde() {
    return Serdes.serdeFrom(
        (topic, data) -> data.toByteArray(),
        (topic, bytes) -> {
          try {
            return MarketBarEvent.parseFrom(bytes);
          } catch (InvalidProtocolBufferException e) {
            throw new RuntimeException("Failed to deserialize MarketBarEvent", e);
          }
        });
  }

  public static Serde<IndicatorEvent> indicatorSerde() {
    return Serdes.serdeFrom(
        (topic, data) -> data.toByteArray(),
        (topic, bytes) -> {
          try {
            return IndicatorEvent.parseFrom(bytes);
          } catch (InvalidProtocolBufferException e) {
            throw new RuntimeException("Failed to deserialize IndicatorEvent", e);
          }
        });
  }
}
