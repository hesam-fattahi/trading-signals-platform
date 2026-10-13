# Market Data Generator Service

A mock service that fetches historical market data using `yfinance` and replays it as a real-time WebSocket tick/bar stream, simulating an external data provider (such as Binance or Alpaca). It is not part of the core trading signal app.

## Configuration (Environment Variables)
- `HOST`: Server bind address (Default: `0.0.0.0`)
- `PORT`: Server port (Default: `8080`)
- `TICKERS`: Market symbol to fetch (Default: `BTC-USD`)
- `PERIOD`: Historical data lookback period (Default: `5d`)
- `INTERVAL`: Historical bar granularity (Default: `1m`)
- `STREAM_SPEED_SEC`: Replay interval per bar (Default: `0.2` seconds)

## How to Run with Docker

Build the container locally:

```bash
# Build the image
docker build -t trading-generator:latest .
```
or in the root project directory, run:
```bash
make build-generator
```

Then run the container with custom options:

```bash
docker run --rm -p 8080:8080 \
  -e TICKERS="BTC-USD" \
  -e STREAM_SPEED_SEC="0.2" \
  trading-generator:latest
```

## Testing the WebSocket Stream

Once running, you can connect to the websocket stream to verify data output using `wscat`:

```bash
npx wscat -c ws://localhost:8080
```