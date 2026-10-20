import { useAppStore } from '../store/useAppStore';

export function BottomBar() {
  const bars = useAppStore(state => state.bars);
  const signals = useAppStore(state => state.signals);

  // bars is FIFO newest-first, so bars[0] is the current price.
  const currentPrice = bars.length > 0 ? bars[0].close : 0;

  const { position, entryPrice } = useAppStore.getState().getCurrentPosition();

  const isLong = position === 'LONG';
  const isShort = position === 'SHORT';
  const inPosition = position !== 'FLAT';

  let pnlPct = 0;
  if (inPosition && entryPrice && currentPrice > 0) {
    pnlPct = isLong
      ? ((currentPrice - entryPrice) / entryPrice) * 100
      : ((entryPrice - currentPrice) / entryPrice) * 100;
  }

  const positionText = isLong ? 'Long' : isShort ? 'Short' : 'Flat';
  const positionColor = isLong ? 'text-[#22c55e]' : isShort ? 'text-[#ef4444]' : 'text-white';
  const dotColor = isLong ? 'bg-[#22c55e]' : isShort ? 'bg-[#ef4444]' : 'bg-gray-500';

  // Keep the signals reference visible to force re-render when it changes.
  void signals.length;

  return (
    <div className="w-full bg-[#0d0d0d] border-t border-[#222222] px-6 py-3 flex justify-between items-center text-sm shrink-0">
      <div className="flex items-center gap-3 font-medium text-gray-400">
        <span className="text-white font-mono">
          ${currentPrice.toLocaleString('en-US', { minimumFractionDigits: 2, maximumFractionDigits: 2 })}
        </span>
        <span>·</span>
        <span>Live</span>
      </div>

      <div className="flex items-center gap-2 text-gray-400 font-medium">
        <div className={`w-2 h-2 rounded-full ${dotColor}`} />
        <span>Current position =</span>
        <span className={positionColor}>{positionText}</span>
        {inPosition && entryPrice && (
          <>
            <span>·</span>
            <span>Entry ${entryPrice.toLocaleString('en-US', { minimumFractionDigits: 2 })}</span>
            <span>·</span>
            <span className={`font-mono ${pnlPct >= 0 ? 'text-[#22c55e]' : 'text-[#ef4444]'}`}>
              {pnlPct >= 0 ? '+' : ''}{pnlPct.toFixed(2)}%
            </span>
          </>
        )}
      </div>
    </div>
  );
}