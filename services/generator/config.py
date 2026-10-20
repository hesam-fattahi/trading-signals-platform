import os

HOST = os.getenv("HOST", "0.0.0.0")
PORT = int(os.getenv("PORT", "8080"))

TICKER = os.getenv("TICKERS")
PERIOD = os.getenv("PERIOD")
INTERVAL = os.getenv("INTERVAL")
STREAM_SPEED_SEC = float(os.getenv("STREAM_SPEED_SEC", "0.2"))
