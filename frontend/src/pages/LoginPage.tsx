import React, { useState, useEffect } from 'react';
import { useSearchParams, useNavigate } from 'react-router-dom';
import { LogIn, UserPlus, Lock, User, AlertCircle, CheckCircle } from 'lucide-react';
import { api } from '../services/api';
import { useAuth } from '../context/AuthContext';

export const LoginPage: React.FC = () => {
  const [searchParams, setSearchParams] = useSearchParams();
  const navigate = useNavigate();
  const { loginUser } = useAuth();

  const isSignup = searchParams.get('mode') === 'signup';

  const [username, setUsername] = useState('');
  const [password, setPassword] = useState('');
  const [error, setError] = useState<string | null>(null);
  const [successMsg, setSuccessMsg] = useState<string | null>(null);
  const [isSubmitting, setIsSubmitting] = useState(false);

  useEffect(() => {
    setError(null);
    setSuccessMsg(null);
  }, [isSignup]);

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    setError(null);
    setSuccessMsg(null);
    setIsSubmitting(true);

    try {
      if (isSignup) {
        await api.register({ username, password });
        setSuccessMsg('Account created successfully! Logging you in...');
        await api.login({ username, password });
      } else {
        await api.login({ username, password });
      }

      loginUser(username);
      navigate('/dashboard');
    } catch (err: any) {
      setError(err.message || 'Authentication failed. Please check your credentials.');
    } finally {
      setIsSubmitting(false);
    }
  };

  const toggleMode = (mode: 'login' | 'signup') => {
    if (mode === 'signup') {
      setSearchParams({ mode: 'signup' });
    } else {
      setSearchParams({});
    }
  };

  return (
    <div className="min-h-[calc(100vh-4rem)] flex items-center justify-center p-4 bg-stone-50/50">
      <div className="w-full max-w-md">
        {/* Card Header & Toggle */}
        <div className="card-white p-8 shadow-lg shadow-orange-500/5">
          <div className="text-center mb-8">
            <div className="w-12 h-12 rounded-2xl bg-orange-100 flex items-center justify-center text-orange-600 mx-auto mb-3">
              {isSignup ? <UserPlus className="w-6 h-6" /> : <LogIn className="w-6 h-6" />}
            </div>
            <h1 className="font-heading text-2xl font-bold text-stone-900">
              {isSignup ? 'Create Gateway Account' : 'Welcome Back'}
            </h1>
            <p className="font-body text-xs text-stone-500 mt-1">
              {isSignup
                ? 'Register to start managing projects and routes'
                : 'Sign in to access your Gateway Dashboard'}
            </p>
          </div>

          {/* Tab Switcher */}
          <div className="flex p-1 bg-stone-100 rounded-xl mb-6">
            <button
              type="button"
              onClick={() => toggleMode('login')}
              className={`flex-1 py-2 text-xs font-semibold rounded-lg transition-all cursor-pointer ${
                !isSignup
                  ? 'bg-white text-stone-900 shadow-xs'
                  : 'text-stone-500 hover:text-stone-800'
              }`}
            >
              Sign In
            </button>
            <button
              type="button"
              onClick={() => toggleMode('signup')}
              className={`flex-1 py-2 text-xs font-semibold rounded-lg transition-all cursor-pointer ${
                isSignup
                  ? 'bg-white text-stone-900 shadow-xs'
                  : 'text-stone-500 hover:text-stone-800'
              }`}
            >
              Sign Up
            </button>
          </div>

          {/* Feedback messages */}
          {error && (
            <div className="flex items-start gap-2.5 p-3.5 mb-6 rounded-lg bg-red-50 border border-red-200 text-red-700 text-xs">
              <AlertCircle className="w-4 h-4 shrink-0 mt-0.5" />
              <span>{error}</span>
            </div>
          )}

          {successMsg && (
            <div className="flex items-start gap-2.5 p-3.5 mb-6 rounded-lg bg-emerald-50 border border-emerald-200 text-emerald-700 text-xs">
              <CheckCircle className="w-4 h-4 shrink-0 mt-0.5" />
              <span>{successMsg}</span>
            </div>
          )}

          {/* Form */}
          <form onSubmit={handleSubmit} className="space-y-4">
            <div>
              <label className="block text-xs font-semibold text-stone-700 mb-1.5">
                Username
              </label>
              <div className="relative">
                <User className="w-4 h-4 text-stone-400 absolute left-3.5 top-1/2 -translate-y-1/2" />
                <input
                  type="text"
                  required
                  value={username}
                  onChange={(e) => setUsername(e.target.value)}
                  placeholder="e.g. john_doe"
                  className="input-field pl-10"
                />
              </div>
            </div>

            <div>
              <label className="block text-xs font-semibold text-stone-700 mb-1.5">
                Password
              </label>
              <div className="relative">
                <Lock className="w-4 h-4 text-stone-400 absolute left-3.5 top-1/2 -translate-y-1/2" />
                <input
                  type="password"
                  required
                  value={password}
                  onChange={(e) => setPassword(e.target.value)}
                  placeholder="••••••••••••"
                  className="input-field pl-10"
                />
              </div>
            </div>

            <button
              type="submit"
              disabled={isSubmitting}
              className="btn-orange-primary w-full py-2.5 mt-2"
            >
              {isSignup ? 'Create Account' : 'Sign In'}
            </button>
          </form>
        </div>
      </div>
    </div>
  );
};
