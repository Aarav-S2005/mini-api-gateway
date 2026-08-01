import React, { createContext, useContext, useState } from 'react';
import { api } from '../services/api';

interface AuthContextType {
  isLoggedIn: boolean;
  username: string | null;
  loginUser: (user: string) => void;
  logoutUser: () => Promise<void>;
  setIsLoggedIn: (status: boolean) => void;
}

const AuthContext = createContext<AuthContextType | undefined>(undefined);

export const AuthProvider: React.FC<{ children: React.ReactNode }> = ({ children }) => {
  const [username, setUsername] = useState<string | null>(() => {
    return localStorage.getItem('gateway_user');
  });

  const [isLoggedInState, setIsLoggedInState] = useState<boolean>(() => {
    return localStorage.getItem('gateway_is_logged_in') === 'true';
  });

  const loginUser = (user: string) => {
    setUsername(user);
    setIsLoggedInState(true);
    localStorage.setItem('gateway_user', user);
    localStorage.setItem('gateway_is_logged_in', 'true');
  };

  const logoutUser = async () => {
    try {
      await api.logout();
    } catch {
      // Ignore logout error
    } finally {
      setUsername(null);
      setIsLoggedInState(false);
      localStorage.removeItem('gateway_user');
      localStorage.removeItem('gateway_is_logged_in');
    }
  };

  const setIsLoggedIn = (status: boolean) => {
    setIsLoggedInState(status);
    localStorage.setItem('gateway_is_logged_in', status ? 'true' : 'false');
    if (!status) {
      setUsername(null);
      localStorage.removeItem('gateway_user');
    }
  };

  const isLoggedIn = isLoggedInState || !!username;

  return (
    <AuthContext.Provider value={{ isLoggedIn, username, loginUser, logoutUser, setIsLoggedIn }}>
      {children}
    </AuthContext.Provider>
  );
};

export const useAuth = (): AuthContextType => {
  const context = useContext(AuthContext);
  if (!context) {
    throw new Error('useAuth must be used within an AuthProvider');
  }
  return context;
};
