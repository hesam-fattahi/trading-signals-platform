import { Chart } from '../components/Chart';
import { BottomBar } from '../components/BottomBar';

export function PlotPage() {
  return (
    <div className="relative w-full h-full bg-[#0a0a0a] flex flex-col">
      <div className="absolute top-8 left-8 z-40 pointer-events-none">
        <h1 className="text-2xl font-bold text-white tracking-wide">BTCUSD</h1>
        <p className="text-sm text-gray-500">Bitcoin / US Dollar · 1m</p>
      </div>

      <div className="flex-1 w-full relative overflow-hidden">
        <Chart />
      </div>

      <BottomBar />
    </div>
  );
}