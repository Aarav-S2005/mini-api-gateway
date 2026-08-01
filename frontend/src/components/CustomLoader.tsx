import React from 'react';
import { useLoading } from '../context/LoadingContext';

export const CustomLoader: React.FC = () => {
  const { isLoading } = useLoading();

  if (!isLoading) return null;

  return (
    <>
      {/* Top loader bar with orange glow */}
      <div className="fixed top-0 left-0 right-0 z-50 h-1 bg-orange-100 overflow-hidden shadow-xs">
        <div className="h-full bg-gradient-to-r from-orange-400 via-orange-600 to-orange-500 animate-loader-bar shadow-sm shadow-orange-500/50" />
      </div>

      {/* Floating status pill bottom right */}
      <div className="fixed bottom-5 right-5 z-50 flex items-center gap-3 px-4 py-2.5 bg-white/95 backdrop-blur-md border border-orange-200/80 rounded-full shadow-lg shadow-orange-500/10 text-xs font-medium text-stone-700 animate-bounce">
        <div className="relative flex items-center justify-center">
          <div className="w-3.5 h-3.5 border-2 border-orange-200 border-t-orange-600 rounded-full animate-spin" />
        </div>
        <span className="font-heading font-semibold text-orange-600">Gateway Syncing...</span>
      </div>
    </>
  );
};
