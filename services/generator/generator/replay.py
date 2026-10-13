import asyncio
import json
import logging
import time

logger = logging.getLogger(__name__)


class MarketDataStreamer:

    def __init__(
        self,
        data_points: list[dict],
        websocket_server,
        interval_seconds: float,
    ):
        self.data_points = data_points
        self.websocket_server = websocket_server
        self.interval_seconds = interval_seconds

    async def run(self):
        if not self.data_points:
            raise RuntimeError("Cannot stream empty market data")

        index = 0
        sequence = 0

        while True:
            # Only process and broadcast when a client is actually connected
            if self.websocket_server.has_clients:
                point = self.data_points[index].copy()
                sequence += 1

                point["sequence_id"] = sequence
                point["stream_time_ms"] = int(time.time() * 1000)

                payload = json.dumps(point)
                self.websocket_server.broadcast(payload)

                logger.info(
                    "Broadcasted symbol=%s price=%s sequence=%d",
                    point["symbol"],
                    point["price"],
                    sequence,
                )

                index = (index + 1) % len(self.data_points)

            await asyncio.sleep(self.interval_seconds)