import { create } from 'zustand';
import { Bar, Indicator, Signal, Trade } from '../types';
import { computeCurrentPosition, computeTrades } from '../utils/trades';

const MAX_BARS = 500;
const MAX_INDICATORS = 500;
const MAX_SIGNALS = 100;

interface AppState {
  isConnected: boolean;
  bars: Bar[];
  indicators: Indicator[];
  signals: Signal[];
  visibleIndicators: {
    bbUpper: boolean;
    bbMiddle: boolean;
    bbLower: boolean;
    adx: boolean;
    rsi: boolean;
  };

  setConnected: (status: boolean) => void;
  addBar: (bar: Bar) => void;
  addIndicator: (indicator: Indicator) => void;
  addSignal: (signal: Signal) => void;
  toggleIndicator: (key: keyof AppState['visibleIndicators']) => void;
  getTrades: () => Trade[];
  getCurrentPosition: () => { position: 'FLAT' | 'LONG' | 'SHORT'; entryPrice: number | null };
}

export const useAppStore = create<AppState>((set, get) => ({
  isConnected: false,
  bars: [],
  indicators: [],
  signals: [],
  visibleIndicators: {
    bbUpper: true,
    bbMiddle: true,
    bbLower: true,
    adx: true,
    rsi: true,
  },

  setConnected: (status) => set({ isConnected: status }),

  // FIFO: newest first, trim the tail when we exceed MAX
  addBar: (bar) => set((state) => ({
    bars: [bar, ...state.bars].slice(0, MAX_BARS),
  })),
  addIndicator: (indicator) => set((state) => ({
    indicators: [indicator, ...state.indicators].slice(0, MAX_INDICATORS),
  })),
  addSignal: (signal) => set((state) => ({
    signals: [signal, ...state.signals].slice(0, MAX_SIGNALS),
  })),

  toggleIndicator: (key) => set((state) => ({
    visibleIndicators: {
      ...state.visibleIndicators,
      [key]: !state.visibleIndicators[key],
    },
  })),

  getTrades: () => computeTrades(get().signals),
  getCurrentPosition: () => computeCurrentPosition(get().signals),
}));