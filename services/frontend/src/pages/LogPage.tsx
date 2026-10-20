import { useMemo, useState } from 'react';
import { useAppStore } from '../store/useAppStore';
import { computeTrades } from '../utils/trades';

const formatTime = (timestamp: number) => {
  const ms = timestamp > 10000000000 ? timestamp : timestamp * 1000;
  return new Date(ms).toLocaleTimeString('en-GB', {
    hour12: false,
    hour: '2-digit',
    minute: '2-digit',
    second: '2-digit',
    timeZone: 'UTC',
  });
};

const formatPrice = (value: number) =>
  value.toLocaleString('en-US', { minimumFractionDigits: 2, maximumFractionDigits: 2 });

export function LogPage() {
  const signals = useAppStore(state => state.signals);
  const [filter, setFilter] = useState<'all' | 'longs' | 'shorts'>('all');

  const trades = useMemo(() => computeTrades(signals), [signals]);

  const filteredTrades = useMemo(
    () =>
      trades.filter(t => {
        if (filter === 'longs') return t.isLong;
        if (filter === 'shorts') return !t.isLong;
        return true;
      }),
    [trades, filter],
  );

  // computeTrades emits chronological; reverse for newest-trade-on-top.
  const displayTrades = useMemo(() => [...filteredTrades].reverse(), [filteredTrades]);

  return (
    <div className="absolute inset-0 bg-[#0a0a0a] pt-24 px-8 pb-8 flex flex-col overflow-hidden">
      <div className="mb-6 flex justify-between items-end shrink-0">
        <div>
          <h1 className="text-2xl font-bold text-white">Signal logs</h1>
          <p className="text-sm text-gray-500 mt-1">BTCUSD · Live replay session</p>
        </div>
        <div className="text-sm text-gray-400">Live</div>
      </div>

      <div className="flex gap-2 mb-4 shrink-0">
        <button onClick={() => setFilter('all')} className={`px-4 py-1.5 rounded-lg text-sm font-medium transition-colors ${filter === 'all' ? 'bg-[#222222] text-white' : 'text-gray-400 hover:text-white'}`}>All signals</button>
        <button onClick={() => setFilter('shorts')} className={`px-4 py-1.5 rounded-lg text-sm font-medium transition-colors ${filter === 'shorts' ? 'bg-[#222222] text-white' : 'text-gray-400 hover:text-white'}`}>Shorts</button>
        <button onClick={() => setFilter('longs')} className={`px-4 py-1.5 rounded-lg text-sm font-medium transition-colors ${filter === 'longs' ? 'bg-[#222222] text-white' : 'text-gray-400 hover:text-white'}`}>Longs</button>
      </div>

      <div className="flex-1 min-h-0 overflow-y-auto bg-[#0a0a0a] border border-[#222222] rounded-xl">
        <table className="w-full text-left border-collapse">
          <thead className="sticky top-0 bg-[#0e0e0e] z-10">
            <tr className="text-xs text-gray-500 uppercase tracking-wider border-b border-[#222222]">
              <th className="px-6 py-4 font-medium">Time</th>
              <th className="px-6 py-4 font-medium">Symbol</th>
              <th className="px-6 py-4 font-medium">Action</th>
              <th className="px-6 py-4 font-medium text-right">Price</th>
              <th className="px-6 py-4 font-medium text-right">Entry</th>
              <th className="px-6 py-4 font-medium text-right">P&L %</th>
            </tr>
          </thead>
          <tbody className="text-sm">
            {displayTrades.map((t, idx) => (
              <tr key={`${t.time}-${t.action}-${idx}`} className="border-b border-[#161616] hover:bg-[#121212] transition-colors">
                <td className="px-6 py-4 font-mono text-gray-400">{formatTime(t.time)}</td>
                <td className="px-6 py-4 text-gray-300">BTCUSD</td>
                <td className="px-6 py-4">
                  <span className="flex items-center gap-2">
                    {t.isLong ? (
                      <svg className="w-4 h-4 text-[#22c55e]" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                        <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M7 17L17 7M17 7H7M17 7v10" />
                      </svg>
                    ) : (
                      <svg className="w-4 h-4 text-[#ef4444]" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                        <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M7 7l10 10M17 17V7M17 17H7" />
                      </svg>
                    )}
                    <span className="text-white font-medium">
                      {t.isLong ? 'Long' : 'Short'} {t.action.includes('ENTRY') ? 'entry' : 'exit'}
                    </span>
                  </span>
                </td>
                <td className="px-6 py-4 font-mono text-right text-gray-300">${formatPrice(t.price)}</td>
                <td className="px-6 py-4 font-mono text-right text-gray-500">{t.entryPrice ? `$${formatPrice(t.entryPrice)}` : '–'}</td>
                <td className={`px-6 py-4 font-mono text-right ${t.pnlPct !== null ? (t.pnlPct >= 0 ? 'text-[#22c55e]' : 'text-[#ef4444]') : 'text-gray-500'}`}>
                  {t.pnlPct !== null
                    ? `${t.pnlPct >= 0 ? '+' : ''}${t.pnlPct.toFixed(2)}%`
                    : '–'}
                </td>
              </tr>
            ))}
          </tbody>
        </table>
      </div>

      <div className="flex justify-between items-center mt-4 text-xs text-gray-500 shrink-0">
        <span>{displayTrades.length} signals · Live</span>
        <span>UTC · Prices in USD · P&L as % of entry, before fees</span>
      </div>
    </div>
  );
}