import React, { useEffect } from 'react';
import { BrowserRouter, Routes, Route, Navigate } from 'react-router-dom';
import { LoadingProvider, useLoading } from './context/LoadingContext';
import { AuthProvider } from './context/AuthContext';
import { CustomLoader } from './components/CustomLoader';
import { Navbar } from './components/Navbar';
import { setApiLoaderTrigger } from './services/api';

import { HomePage } from './pages/HomePage';
import { LoginPage } from './pages/LoginPage';
import { DashboardPage } from './pages/DashboardPage';
import { ProjectDetailPage } from './pages/ProjectDetailPage';
import { UpstreamsPage } from './pages/UpstreamsPage';
import { RoutesPage } from './pages/RoutesPage';

// Inner component to hook up loading context to API service
const AppContent: React.FC = () => {
  const { withLoading } = useLoading();

  useEffect(() => {
    setApiLoaderTrigger(withLoading);
  }, [withLoading]);

  return (
    <div className="min-h-screen bg-stone-50/50 flex flex-col font-body selection:bg-orange-500 selection:text-white">
      <CustomLoader />
      <Navbar />

      <main className="flex-1">
        <Routes>
          <Route path="/" element={<Navigate to="/home" replace />} />
          <Route path="/home" element={<HomePage />} />
          <Route path="/login" element={<LoginPage />} />
          <Route path="/dashboard" element={<DashboardPage />} />
          <Route path="/dashboard/project/:projectID" element={<ProjectDetailPage />} />
          <Route path="/dashboard/project/:projectID/upstreams" element={<UpstreamsPage />} />
          <Route path="/dashboard/project/:projectID/routes" element={<RoutesPage />} />
          <Route path="*" element={<Navigate to="/home" replace />} />
        </Routes>
      </main>
    </div>
  );
};

export function App() {
  return (
    <AuthProvider>
      <LoadingProvider>
        <BrowserRouter>
          <AppContent />
        </BrowserRouter>
      </LoadingProvider>
    </AuthProvider>
  );
}

export default App;
