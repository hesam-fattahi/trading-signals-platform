import { useAppStore } from '../store/useAppStore';

export function DisconnectOverlay() {
  const isConnected = useAppStore(state => state.isConnected);
  if (isConnected) return null;

  return (
    <div className="fixed inset-0 z-[100] backdrop-blur-md bg-black/40 flex items-center justify-center">
      <div className="bg-[#141414] border border-[#2a2a2a] px-8 py-6 rounded-2xl shadow-2xl flex flex-col items-center">
        <div className="w-8 h-8 border-4 border-blue-500 border-t-transparent rounded-full animate-spin mb-4"></div>
        <h2 className="text-white font-semibold text-lg">WebSocket Disconnected</h2>
        <p className="text-gray-400 text-sm mt-1">Retrying connection...</p>
      </div>
    </div>
  );
}