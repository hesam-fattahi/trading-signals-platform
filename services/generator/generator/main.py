import asyncio
import logging
import websockets

from config import HOST, PORT, TICKER, PERIOD, INTERVAL, STREAM_SPEED_SEC
from fetch_market_data import fetch_market_data
from replay import MarketDataStreamer
from websocket_server import WebSocketServer

logging.basicConfig(
    level=logging.INFO,
    format="%(asctime)s [%(levelname)s] %(message)s",
)


async def main():
    data_points = fetch_market_data(
        ticker=TICKER,
        period=PERIOD,
        interval=INTERVAL,
    )

    websocket_server = WebSocketServer()

    server = await websockets.serve(
        websocket_server.handler,
        HOST,
        PORT,
    )

    streamer = MarketDataStreamer(
        data_points=data_points,
        websocket_server=websocket_server,
        interval_seconds=STREAM_SPEED_SEC,
    )

    logging.info("Generator listening on ws://%s:%s (Streaming ticker=%s)", HOST, PORT, TICKER)

    await asyncio.gather(
        server.wait_closed(),
        streamer.run(),
    )


if __name__ == "__main__":
    asyncio.run(main())