import { useEffect, useRef } from 'react';
import { createChart, ColorType, CandlestickSeries, HistogramSeries, LineSeries } from 'lightweight-charts';
import { useAppStore } from '../store/useAppStore';

export function Chart() {
  const chartContainerRef = useRef<HTMLDivElement>(null);

  const candleSeriesRef = useRef<any>(null);
  const volumeSeriesRef = useRef<any>(null);
  const bbUpperRef = useRef<any>(null);
  const bbMiddleRef = useRef<any>(null);
  const bbLowerRef = useRef<any>(null);
  const rsiSeriesRef = useRef<any>(null);
  const adxSeriesRef = useRef<any>(null);

  const bars = useAppStore(state => state.bars);
  const indicators = useAppStore(state => state.indicators);
  const visibleIndicators = useAppStore(state => state.visibleIndicators);
  const toggleIndicator = useAppStore(state => state.toggleIndicator);
  const latestIndicator = indicators.length > 0 ? indicators[0] : null; // FIFO: newest at [0]

  useEffect(() => {
    if (!chartContainerRef.current) return;

    const chart = createChart(chartContainerRef.current, {
      layout: {
        background: { type: ColorType.Solid, color: '#0a0a0a' },
        textColor: '#666666',
        fontFamily: 'Inter, sans-serif',
      },
      grid: {
        vertLines: { color: '#141414' },
        horzLines: { color: '#141414' },
      },
      timeScale: {
        borderColor: '#222222',
        timeVisible: true,
        secondsVisible: false,
      },
      rightPriceScale: {
        borderColor: '#222222',
      },
      crosshair: {
        mode: 1,
        vertLine: { color: '#333333', style: 3 },
        horzLine: { color: '#333333', style: 3 },
      },
    });

    const candleSeries = chart.addSeries(CandlestickSeries, {
      upColor: '#22c55e',
      downColor: '#ef4444',
      borderVisible: false,
      wickUpColor: '#22c55e',
      wickDownColor: '#ef4444',
    });
    candleSeriesRef.current = candleSeries;

    bbUpperRef.current = chart.addSeries(LineSeries, { color: '#3b82f6', lineWidth: 1, crosshairMarkerVisible: false, lastValueVisible: false, priceLineVisible: false });
    bbMiddleRef.current = chart.addSeries(LineSeries, { color: '#888888', lineWidth: 1, lineStyle: 2, crosshairMarkerVisible: false, lastValueVisible: false, priceLineVisible: false });
    bbLowerRef.current = chart.addSeries(LineSeries, { color: '#3b82f6', lineWidth: 1, crosshairMarkerVisible: false, lastValueVisible: false, priceLineVisible: false });

    const volumeSeries = chart.addSeries(HistogramSeries, {
      color: '#22c55e',
      priceFormat: { type: 'volume' },
      priceScaleId: '',
    });
    volumeSeries.priceScale().applyOptions({ scaleMargins: { top: 0.8, bottom: 0 } });
    volumeSeriesRef.current = volumeSeries;

    adxSeriesRef.current = chart.addSeries(LineSeries, { color: '#f59e0b', lineWidth: 1, priceLineVisible: false, lastValueVisible: false }, 2);
    rsiSeriesRef.current = chart.addSeries(LineSeries, { color: '#06b6d4', lineWidth: 1, priceLineVisible: false, lastValueVisible: false }, 3);

    const handleResize = () => {
      if (chartContainerRef.current) {
        chart.applyOptions({ width: chartContainerRef.current.clientWidth, height: chartContainerRef.current.clientHeight });
      }
    };
    window.addEventListener('resize', handleResize);

    return () => {
      window.removeEventListener('resize', handleResize);
      chart.remove();
    };
  }, []);

  useEffect(() => {
    if (!candleSeriesRef.current || bars.length === 0) return;

    // Reverse FIFO to chronological for lightweight-charts.
    const chronological = [...bars].reverse();

    const formattedBars = chronological.map(b => ({
      time: b.time as any, open: b.open, high: b.high, low: b.low, close: b.close,
    }));
    candleSeriesRef.current.setData(formattedBars);

    const formattedVolume = chronological.map(b => ({
      time: b.time as any, value: b.volume, color: b.close >= b.open ? '#22c55e33' : '#ef444433',
    }));
    volumeSeriesRef.current.setData(formattedVolume);
  }, [bars]);

  useEffect(() => {
    if (indicators.length === 0) return;

    const chronological = [...indicators].reverse();

    if (bbUpperRef.current) bbUpperRef.current.setData(visibleIndicators.bbUpper ? chronological.map(i => ({ time: i.time as any, value: i.bb_upper })) : []);
    if (bbMiddleRef.current) bbMiddleRef.current.setData(visibleIndicators.bbMiddle ? chronological.map(i => ({ time: i.time as any, value: i.bb_middle })) : []);
    if (bbLowerRef.current) bbLowerRef.current.setData(visibleIndicators.bbLower ? chronological.map(i => ({ time: i.time as any, value: i.bb_lower })) : []);
    if (adxSeriesRef.current) adxSeriesRef.current.setData(visibleIndicators.adx ? chronological.map(i => ({ time: i.time as any, value: i.adx })) : []);
    if (rsiSeriesRef.current) rsiSeriesRef.current.setData(visibleIndicators.rsi ? chronological.map(i => ({ time: i.time as any, value: i.rsi })) : []);
  }, [indicators, visibleIndicators]);

  return (
    <div className="relative w-full h-full">
      <div ref={chartContainerRef} className="w-full h-full" />

      <div className="absolute top-24 left-4 z-20 bg-[#121212]/10 backdrop-blur-md border border-[#222222] rounded-lg p-3 text-xs space-y-1.5 shadow-xl">
        <button type="button" aria-pressed={visibleIndicators.bbUpper} onClick={() => toggleIndicator('bbUpper')} className={`flex items-center gap-2 text-left ${visibleIndicators.bbUpper ? '' : 'opacity-40'}`}>
          <div className="w-3 h-0.5 bg-[#3b82f6]" /><span className="text-gray-300">Bollinger upper</span>
        </button>
        <button type="button" aria-pressed={visibleIndicators.bbMiddle} onClick={() => toggleIndicator('bbMiddle')} className={`flex items-center gap-2 text-left ${visibleIndicators.bbMiddle ? '' : 'opacity-40'}`}>
          <div className="w-3 h-0.5 bg-[#888888] border-dashed border-t" /><span className="text-gray-300">Bollinger middle</span>
        </button>
        <button type="button" aria-pressed={visibleIndicators.bbLower} onClick={() => toggleIndicator('bbLower')} className={`flex items-center gap-2 text-left ${visibleIndicators.bbLower ? '' : 'opacity-40'}`}>
          <div className="w-3 h-0.5 bg-[#3b82f6]" /><span className="text-gray-300">Bollinger lower</span>
        </button>
        <button type="button" aria-pressed={visibleIndicators.adx} onClick={() => toggleIndicator('adx')} className={`flex items-center gap-2 text-left ${visibleIndicators.adx ? '' : 'opacity-40'}`}>
          <div className="w-3 h-0.5 bg-[#f59e0b]" /><span className="text-gray-300">ADX · 14</span>
        </button>
        <button type="button" aria-pressed={visibleIndicators.rsi} onClick={() => toggleIndicator('rsi')} className={`flex items-center gap-2 text-left ${visibleIndicators.rsi ? '' : 'opacity-40'}`}>
          <div className="w-3 h-0.5 bg-[#06b6d4]" /><span className="text-gray-300">RSI · 14</span>
        </button>
      </div>

      {visibleIndicators.adx && latestIndicator && (
        <div className="absolute top-[68%] left-6 z-20 text-xs font-mono text-gray-400 pointer-events-none">
          <span className="text-[#f59e0b] font-semibold">ADX 14</span> {latestIndicator.adx.toFixed(2)}
        </div>
      )}

      {visibleIndicators.rsi && latestIndicator && (
        <div className="absolute top-[84%] left-6 z-20 text-xs font-mono text-gray-400 pointer-events-none">
          <span className="text-[#06b6d4] font-semibold">RSI 14</span> {latestIndicator.rsi.toFixed(2)}
        </div>
      )}
    </div>
  );
}