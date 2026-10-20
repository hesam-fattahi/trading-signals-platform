import { BrowserRouter, Routes, Route } from 'react-router-dom';
import { useEffect } from 'react';
import { TopBar } from './components/TopBar';
import { PlotPage } from './pages/PlotPage';
import { LogPage } from './pages/LogPage';
import { DisconnectOverlay } from './components/DisconnectOverlay';
import { connectWebSocket } from './api/websocket';

export default function App() {
  useEffect(() => {
    const disconnect = connectWebSocket();
    return () => disconnect();
  }, []);

  return (
    <BrowserRouter>
      <div className="w-full h-screen bg-[#0a0a0a] text-[#e0e0e0] overflow-hidden flex flex-col">
        <TopBar />
        <DisconnectOverlay />
        <div className="flex-1 relative">
          <Routes>
            <Route path="/" element={<PlotPage />} />
            <Route path="/log" element={<LogPage />} />
          </Routes>
        </div>
      </div>
    </BrowserRouter>
  );
}