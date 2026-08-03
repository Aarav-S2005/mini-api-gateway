import React from 'react';
import { Link, useNavigate, useLocation } from 'react-router-dom';
import { Network, LogOut, LayoutDashboard, Home, ShieldCheck } from 'lucide-react';
import { useAuth } from '../context/AuthContext';

export const Navbar: React.FC = () => {
  const navigate = useNavigate();
  const location = useLocation();
  const { isLoggedIn, username, logoutUser } = useAuth();

  const handleLogout = async () => {
    await logoutUser();
    navigate('/login');
  };

  const isActive = (path: string) => location.pathname.startsWith(path);

  return (
    <header className="sticky top-0 z-40 bg-white/90 backdrop-blur-md border-b border-stone-200/80 shadow-2xs">
      <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 h-16 flex items-center justify-between">
        {/* Brand Logo */}
        <Link to="/home" className="flex items-center gap-2.5 group">
          <div className="w-9 h-9 rounded-xl bg-gradient-to-tr from-orange-600 via-orange-500 to-amber-500 flex items-center justify-center text-white shadow-sm shadow-orange-500/30 group-hover:scale-105 transition-transform">
            <Network className="w-5 h-5 stroke-[2.2]" />
          </div>
          <div className="flex flex-col">
            <span className="font-heading font-bold text-lg text-stone-900 tracking-tight leading-none group-hover:text-orange-600 transition-colors">
              Mini Gateway
            </span>
            <span className="text-[10px] font-medium text-orange-600 uppercase tracking-widest leading-tight">
              API Gateway
            </span>
          </div>
        </Link>

        {/* Navigation Links */}
        <nav className="flex items-center gap-1 sm:gap-2">
          <Link
            to="/home"
            className={`flex items-center gap-2 px-3 py-1.5 rounded-lg text-sm font-medium transition-colors ${
              location.pathname === '/home'
                ? 'bg-orange-50 text-orange-600 font-semibold'
                : 'text-stone-600 hover:text-stone-900 hover:bg-stone-100/60'
            }`}
          >
            <Home className="w-4 h-4" />
            <span>Home</span>
          </Link>

          <Link
            to="/dashboard"
            className={`flex items-center gap-2 px-3 py-1.5 rounded-lg text-sm font-medium transition-colors ${
              isActive('/dashboard')
                ? 'bg-orange-50 text-orange-600 font-semibold'
                : 'text-stone-600 hover:text-stone-900 hover:bg-stone-100/60'
            }`}
          >
            <LayoutDashboard className="w-4 h-4" />
            <span>Dashboard</span>
          </Link>
        </nav>

        {/* User / Auth Action */}
        <div className="flex items-center gap-3">
          {isLoggedIn ? (
            <div className="flex items-center gap-3">
              <div className="hidden sm:flex items-center gap-2 px-3 py-1 bg-stone-100 border border-stone-200 rounded-full text-xs text-stone-700">
                <ShieldCheck className="w-3.5 h-3.5 text-orange-600" />
                <span className="font-medium">{username || 'User'}</span>
              </div>
              <button
                onClick={handleLogout}
                className="btn-white text-xs py-1.5 px-3"
                title="Logout"
              >
                <LogOut className="w-3.5 h-3.5 text-stone-500" />
                <span className="hidden sm:inline">Logout</span>
              </button>
            </div>
          ) : (
            <div className="flex items-center gap-2">
              <Link to="/login" className="btn-white text-xs py-1.5 px-3.5">
                Log In
              </Link>
              <Link to="/login?mode=signup" className="btn-orange-primary text-xs py-1.5 px-3.5">
                Sign Up
              </Link>
            </div>
          )}
        </div>
      </div>
    </header>
  );
};
