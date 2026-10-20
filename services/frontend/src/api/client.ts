// import { Bar, Indicator, Signal, Info } from '../types';

// const API_BASE = '/';

// async function getJSON<T>(url: string, signal?: AbortSignal): Promise<T> {
//   const res = await fetch(url, { signal });
//   if (!res.ok) {
//     throw new Error(`Request failed (${res.status} ${res.statusText}): ${url}`);
//   }
//   return res.json() as Promise<T>;
// }

// export function fetchInfo(signal?: AbortSignal): Promise<Info> {
//   return getJSON<Info>(`${API_BASE}info`, signal);
// }

// export function fetchHistoryBars(limit = 100, before?: number): Promise<Bar[]> {
//   const params = new URLSearchParams({ limit: String(limit) });
//   if (before !== undefined) {
//     params.set('before', String(before));
//   }
//   return getJSON<Bar[]>(`${API_BASE}history/bars?${params.toString()}`);
// }

// export function fetchHistoryIndicators(limit = 100): Promise<Indicator[]> {
//   return getJSON<Indicator[]>(`${API_BASE}history/indicators?limit=${limit}`);
// }

// export function fetchHistorySignals(limit = 50): Promise<Signal[]> {
//   return getJSON<Signal[]>(`${API_BASE}history/signals?limit=${limit}`);
// }

export {};
