import { Link, useLocation } from 'react-router-dom';

export function TopBar() {
  const location = useLocation();

  return (
    <div className="absolute top-6 left-1/2 -translate-x-1/2 z-50">
      <div className="flex items-center gap-3 bg-[#121212]/10 backdrop-blur-md border border-[#222222] rounded-xl px-2 py-1.5 shadow-2xl">        
        <Link 
          to="/" 
          className={`px-6 py-1.5 rounded-lg text-sm font-medium transition-colors ${location.pathname === '/' ? 'bg-[#222222] text-white' : 'text-gray-400 hover:text-white'}`}
        >
          Plot
        </Link>
        <Link 
          to="/log" 
          className={`px-6 py-1.5 rounded-lg text-sm font-medium transition-colors ${location.pathname === '/log' ? 'bg-[#222222] text-white' : 'text-gray-400 hover:text-white'}`}
        >
          Log
        </Link>
      </div>
    </div>
  );
}