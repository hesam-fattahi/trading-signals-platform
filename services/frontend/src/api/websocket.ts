import { useAppStore } from '../store/useAppStore';
import { WsMessage } from '../types';

let activeCleanup: (() => void) | null = null;
let generation = 0;

export function connectWebSocket(): () => void {
  activeCleanup?.();
  const connectionGeneration = ++generation;
  let disposed = false;
  let currentSocket: WebSocket | null = null;
  let reconnectTimer: number | null = null;
  let reconnectDelay = 1000;

  const openSocket = () => {
    if (disposed || connectionGeneration !== generation) return;

    const protocol = window.location.protocol === 'https:' ? 'wss:' : 'ws:';
    const socket = new WebSocket(`${protocol}//${window.location.host}/ws`);
    currentSocket = socket;

    socket.onopen = () => {
      if (disposed || connectionGeneration !== generation || currentSocket !== socket) return;
      console.log('WebSocket connected');
      useAppStore.getState().setConnected(true);
      reconnectDelay = 1000;
    };

    socket.onmessage = (event) => {
      if (disposed || connectionGeneration !== generation || currentSocket !== socket) return;
      try {
        const msg = JSON.parse(event.data) as WsMessage;
        const store = useAppStore.getState();
        if (msg.type === 'tick') store.addBar(msg.data);
        else if (msg.type === 'indicator') store.addIndicator(msg.data);
        else if (msg.type === 'signal') store.addSignal(msg.data);
      } catch (err) {
        console.error('Failed to parse WebSocket message', err);
      }
    };

    socket.onclose = () => {
      if (disposed || connectionGeneration !== generation || currentSocket !== socket) return;
      console.log('WebSocket disconnected');
      useAppStore.getState().setConnected(false);

      reconnectTimer = window.setTimeout(() => {
        reconnectTimer = null;
        reconnectDelay = Math.min(reconnectDelay * 2, 30000);
        openSocket();
      }, reconnectDelay);
    };

    socket.onerror = () => socket.close();
  };

  const cleanup = () => {
    if (disposed) return;
    disposed = true;
    if (connectionGeneration === generation) generation++;

    if (reconnectTimer !== null) {
      window.clearTimeout(reconnectTimer);
      reconnectTimer = null;
    }

    if (currentSocket) {
      currentSocket.onopen = null;
      currentSocket.onmessage = null;
      currentSocket.onclose = null;
      currentSocket.onerror = null;
      if (currentSocket.readyState === WebSocket.CONNECTING || currentSocket.readyState === WebSocket.OPEN) {
        currentSocket.close();
      }
      currentSocket = null;
    }
    useAppStore.getState().setConnected(false);
  };

  const dispose = () => {
    if (activeCleanup === dispose) activeCleanup = null;
    cleanup();
  };

  activeCleanup = dispose;
  openSocket();
  return dispose;
}