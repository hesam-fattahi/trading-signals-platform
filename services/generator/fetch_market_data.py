import logging

import pandas as pd
import yfinance as yf

logger = logging.getLogger(__name__)


def fetch_market_data(
    ticker: str,
    period: str,
    interval: str,
) -> list[dict]:
    logger.info(
        "Downloading market data: ticker=%s period=%s interval=%s",
        ticker,
        period,
        interval,
    )

    df = yf.download(
        tickers=ticker,
        period=period,
        interval=interval,
        auto_adjust=True,
        prepost=False,
        progress=False,
    )

    if df.empty:
        raise RuntimeError(f"No market data returned for ticker {ticker}")

    data_points = []
    normalized_symbol = ticker.replace("-", "")  # "BTC-USD" -> "BTCUSD"

    for timestamp, row in df.iterrows():
        # Handle single ticker pandas series/scalar values
        close_val = (
            row["Close"].item() if isinstance(row["Close"], pd.Series) else row["Close"]
        )
        if pd.isna(close_val):
            continue

        open_val = (
            row["Open"].item() if isinstance(row["Open"], pd.Series) else row["Open"]
        )
        high_val = (
            row["High"].item() if isinstance(row["High"], pd.Series) else row["High"]
        )
        low_val = row["Low"].item() if isinstance(row["Low"], pd.Series) else row["Low"]
        vol_val = (
            row["Volume"].item()
            if isinstance(row["Volume"], pd.Series)
            else row["Volume"]
        )

        data_points.append(
            {
                "symbol": normalized_symbol,
                "timestamp_ms": int(timestamp.timestamp() * 1000),
                "open": round(float(open_val), 4),
                "high": round(float(high_val), 4),
                "low": round(float(low_val), 4),
                "close": round(float(close_val), 4),
                "price": round(float(close_val), 4),
                "volume": round(float(vol_val), 4),
            }
        )

    logger.info(
        "Successfully fetched %d data points for %s",
        len(data_points),
        normalized_symbol,
    )
    return data_points
