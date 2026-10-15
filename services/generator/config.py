import os

HOST = os.getenv("HOST", "0.0.0.0")
PORT = int(os.getenv("PORT", "8080"))

TICKER = os.getenv("TICKERS", "BTC-USD").strip()

PERIOD = os.getenv("PERIOD", "5d")
INTERVAL = os.getenv("INTERVAL", "1m")

STREAM_SPEED_SEC = float(os.getenv("STREAM_SPEED_SEC", "0.2"))
