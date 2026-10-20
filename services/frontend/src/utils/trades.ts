import { Signal, Trade } from '../types';

interface OpenEntry {
  price: number;
  isLong: boolean;
}

export function computeTrades(signals: Signal[]): Trade[] {
  const chronological = [...signals].sort((a, b) => a.time - b.time);
  const trades: Trade[] = [];
  let open: OpenEntry | null = null;

  for (const signal of chronological) {
    const isEntry = signal.action.includes('ENTRY');
    const isLong = signal.action.includes('LONG');

    if (isEntry) {
      if (open) continue;
      trades.push({
        time: signal.time,
        action: signal.action,
        price: signal.price,
        entryPrice: null,
        pnlPct: null,
        isLong,
      });
      open = { price: signal.price, isLong };
    } else if (open && signal.action.includes('EXIT')) {
      if (open.isLong !== isLong) continue;

      const pnlPct = open.isLong
        ? ((signal.price - open.price) / open.price) * 100
        : ((open.price - signal.price) / open.price) * 100;

      trades.push({
        time: signal.time,
        action: signal.action,
        price: signal.price,
        entryPrice: open.price,
        pnlPct,
        isLong: open.isLong,
      });
      open = null;
    }
  }

  return trades;
}

export function computeCurrentPosition(signals: Signal[]): {
  position: 'FLAT' | 'LONG' | 'SHORT';
  entryPrice: number | null;
} {
  const chronological = [...signals].sort((a, b) => a.time - b.time);
  let open: { price: number; isLong: boolean } | null = null;

  for (const signal of chronological) {
    const isEntry = signal.action.includes('ENTRY');
    const isLong = signal.action.includes('LONG');

    if (isEntry) {
      if (!open) open = { price: signal.price, isLong };
      continue;
    }

    if (open && signal.action.includes('EXIT') && open.isLong === isLong) {
      open = null;
    }
  }

  if (!open) return { position: 'FLAT', entryPrice: null };
  return { position: open.isLong ? 'LONG' : 'SHORT', entryPrice: open.price };
}