/** @type {import('tailwindcss').Config} */
export default {
  content: [
    "./index.html",
    "./src/**/*.{js,ts,jsx,tsx}",
  ],
  theme: {
    extend: {
      colors: {
        background: '#0a0a0a',
        panel: '#141414',
        border: '#222222',
        textMain: '#e0e0e0',
        textMuted: '#888888',
        up: '#22c55e',
        down: '#ef4444',
        accent: '#3b82f6',
      },
      fontFamily: {
        sans: ['Inter', '-apple-system', 'sans-serif'],
      },
    },
  },
  plugins: [],
}