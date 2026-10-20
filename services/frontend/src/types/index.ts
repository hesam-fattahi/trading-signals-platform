export interface Bar {
  time: number;
  open: number;
  high: number;
  low: number;
  close: number;
  volume: number;
}

export interface Indicator {
  time: number;
  bb_upper: number;
  bb_middle: number;
  bb_lower: number;
  rsi: number;
  adx: number;
  di_plus: number;
  di_minus: number;
}

export interface Signal {
  time: number;
  action: string;
  price: number;
}

export interface Trade {
  time: number;
  action: string;
  price: number;
  entryPrice: number | null;
  pnlPct: number | null;
  isLong: boolean;
}

export type WsMessage =
  | { type: 'tick'; data: Bar }
  | { type: 'indicator'; data: Indicator }
  | { type: 'signal'; data: Signal };