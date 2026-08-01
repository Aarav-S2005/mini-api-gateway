import React, { useState, useEffect } from 'react';
import { useParams, Link } from 'react-router-dom';
import {
  ArrowLeft,
  Route as RouteIcon,
  Plus,
  Trash2,
  Edit2,
  X,
  AlertCircle,
  CheckCircle,
  ShieldCheck,
  ShieldOff,
  Server,
} from 'lucide-react';
import { api } from '../services/api';
import type {
  GetRouteResponseDTO,
  GetUpstreamResponseDTO,
  PathType,
  AuthMode,
  CreateOrUpdateRouteRequestDTO,
} from '../types/api';

export const RoutesPage: React.FC = () => {
  const { projectID } = useParams<{ projectID: string }>();

  const [routes, setRoutes] = useState<GetRouteResponseDTO[]>([]);
  const [upstreams, setUpstreams] = useState<GetUpstreamResponseDTO[]>([]);
  const [error, setError] = useState<string | null>(null);
  const [successMsg, setSuccessMsg] = useState<string | null>(null);

  // Modal / Form state
  const [isModalOpen, setIsModalOpen] = useState(false);
  const [editingRouteId, setEditingRouteId] = useState<string | null>(null);

  const [path, setPath] = useState('/orders');
  const [pathType, setPathType] = useState<PathType>('prefix');
  const [method, setMethod] = useState('GET');
  const [upstreamName, setUpstreamName] = useState('');
  const [authMode, setAuthMode] = useState<AuthMode>('required');
  const [enabled, setEnabled] = useState(true);

  const fetchData = async () => {
    if (!projectID) return;
    try {
      const [routesRes, upstreamsRes] = await Promise.all([
        api.listRoutes(projectID),
        api.listUpstreams(projectID),
      ]);
      setRoutes(routesRes.routes || []);
      setUpstreams(upstreamsRes.upstreams || []);
      if (upstreamsRes.upstreams?.length > 0 && !upstreamName) {
        setUpstreamName(upstreamsRes.upstreams[0].name);
      }
      setError(null);
    } catch (err: any) {
      setError(err.message || 'Failed to fetch routes or upstreams.');
    }
  };

  useEffect(() => {
    fetchData();
  }, [projectID]);

  const openCreateModal = () => {
    setEditingRouteId(null);
    setPath('/orders');
    setPathType('prefix');
    setMethod('GET');
    setUpstreamName(upstreams[0]?.name || '');
    setAuthMode('required');
    setEnabled(true);
    setIsModalOpen(true);
  };

  const openEditModal = (route: GetRouteResponseDTO) => {
    setEditingRouteId(route.route_id);
    setPath(route.path);
    setPathType(route.path_type);
    setMethod(route.method);
    setUpstreamName(route.upstream_name);
    setAuthMode(route.auth_mode);
    setEnabled(route.enabled);
    setIsModalOpen(true);
  };

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    if (!projectID || !path.trim() || !upstreamName) return;

    const payload: CreateOrUpdateRouteRequestDTO = {
      path: path.trim(),
      path_type: pathType,
      method,
      upstream_name: upstreamName,
      auth_mode: authMode,
      enabled,
    };

    try {
      if (editingRouteId) {
        await api.updateRoute(projectID, editingRouteId, payload);
        setSuccessMsg(`Route '${path}' updated.`);
      } else {
        await api.createRoute(projectID, payload);
        setSuccessMsg(`Route '${path}' created.`);
      }
      setIsModalOpen(false);
      setTimeout(() => setSuccessMsg(null), 3000);
      fetchData();
    } catch (err: any) {
      setError(err.message || 'Failed to save route.');
    }
  };

  const handleDeleteRoute = async (routeId: string, routePath: string) => {
    if (!projectID) return;
    if (!window.confirm(`Delete route '${routePath}'?`)) return;

    try {
      await api.deleteRoute(projectID, routeId);
      setRoutes((prev) => prev.filter((r) => r.route_id !== routeId));
      setSuccessMsg(`Route '${routePath}' deleted.`);
      setTimeout(() => setSuccessMsg(null), 3000);
    } catch (err: any) {
      setError(err.message || 'Failed to delete route.');
    }
  };

  const getMethodBadgeClass = (m: string) => {
    switch (m.toUpperCase()) {
      case 'GET':
        return 'bg-blue-100 text-blue-700 border-blue-200';
      case 'POST':
        return 'bg-emerald-100 text-emerald-700 border-emerald-200';
      case 'PUT':
      case 'PATCH':
        return 'bg-amber-100 text-amber-700 border-amber-200';
      case 'DELETE':
        return 'bg-red-100 text-red-700 border-red-200';
      default:
        return 'bg-stone-100 text-stone-700 border-stone-200';
    }
  };

  return (
    <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-8 space-y-8">
      {/* Header */}
      <div className="flex flex-col sm:flex-row sm:items-center justify-between gap-4 pb-6 border-b border-stone-200">
        <div>
          <Link
            to={`/dashboard/project/${projectID}`}
            className="inline-flex items-center gap-1.5 text-xs font-semibold text-stone-500 hover:text-orange-600 mb-2 transition-colors"
          >
            <ArrowLeft className="w-4 h-4" />
            <span>Back to Project Overview</span>
          </Link>
          <h1 className="font-heading text-3xl font-bold text-stone-900 flex items-center gap-3">
            <RouteIcon className="w-7 h-7 text-orange-600" />
            <span>Gateway Routes</span>
          </h1>
          <p className="font-body text-sm text-stone-500 mt-1">
            Map incoming HTTP request paths and methods to upstream microservices.
          </p>
        </div>
        <button onClick={openCreateModal} className="btn-orange-primary shadow-sm">
          <Plus className="w-4 h-4" />
          <span>New Route</span>
        </button>
      </div>

      {/* Notifications */}
      {error && (
        <div className="flex items-center gap-3 p-4 rounded-xl bg-red-50 border border-red-200 text-red-700 text-sm">
          <AlertCircle className="w-5 h-5 shrink-0" />
          <span>{error}</span>
        </div>
      )}

      {successMsg && (
        <div className="flex items-center gap-3 p-4 rounded-xl bg-emerald-50 border border-emerald-200 text-emerald-700 text-sm">
          <CheckCircle className="w-5 h-5 shrink-0" />
          <span>{successMsg}</span>
        </div>
      )}

      {/* Routes List */}
      {routes.length === 0 ? (
        <div className="card-white p-12 text-center space-y-4">
          <div className="w-12 h-12 rounded-full bg-orange-50 flex items-center justify-center text-orange-500 mx-auto">
            <RouteIcon className="w-6 h-6" />
          </div>
          <div className="space-y-1">
            <h3 className="font-heading text-base font-semibold text-stone-800">
              No routes configured
            </h3>
            <p className="font-body text-xs text-stone-500 max-w-sm mx-auto">
              Add your first HTTP path route mapping requests to an upstream backend.
            </p>
          </div>
          <button onClick={openCreateModal} className="btn-orange-primary text-xs">
            <Plus className="w-3.5 h-3.5" />
            <span>Create Route</span>
          </button>
        </div>
      ) : (
        <div className="space-y-4">
          {routes.map((route) => (
            <div
              key={route.route_id}
              className="card-white p-5 flex flex-col sm:flex-row sm:items-center justify-between gap-4 hover:border-orange-300 transition-all"
            >
              <div className="flex items-center gap-4 flex-wrap">
                <span
                  className={`px-2.5 py-1 text-xs font-bold font-mono-url rounded border ${getMethodBadgeClass(
                    route.method
                  )}`}
                >
                  {route.method}
                </span>

                <div className="space-y-0.5">
                  <div className="flex items-center gap-2">
                    <code className="font-mono-url font-bold text-sm text-stone-900">
                      {route.path}
                    </code>
                    <span className="badge-gray text-[10px] lowercase">{route.path_type}</span>
                  </div>

                  <div className="flex items-center gap-3 text-xs text-stone-500">
                    <span className="flex items-center gap-1">
                      <Server className="w-3.5 h-3.5 text-stone-400" />
                      <span>Target:</span>
                      <strong className="text-stone-800 font-mono-url">{route.upstream_name}</strong>
                    </span>

                    <span className="flex items-center gap-1">
                      {route.auth_mode === 'required' ? (
                        <span className="text-orange-600 flex items-center gap-1 font-semibold">
                          <ShieldCheck className="w-3.5 h-3.5" />
                          <span>Auth Required</span>
                        </span>
                      ) : (
                        <span className="text-stone-400 flex items-center gap-1">
                          <ShieldOff className="w-3.5 h-3.5" />
                          <span>Public</span>
                        </span>
                      )}
                    </span>
                  </div>
                </div>
              </div>

              <div className="flex items-center justify-between sm:justify-end gap-3 pt-3 sm:pt-0 border-t sm:border-t-0 border-stone-100">
                <span className={route.enabled ? 'badge-green' : 'badge-gray'}>
                  {route.enabled ? 'Enabled' : 'Disabled'}
                </span>

                <div className="flex items-center gap-1">
                  <button
                    onClick={() => openEditModal(route)}
                    className="p-1.5 rounded-lg text-stone-400 hover:text-stone-700 hover:bg-stone-100 transition-colors"
                    title="Edit route"
                  >
                    <Edit2 className="w-4 h-4" />
                  </button>
                  <button
                    onClick={() => handleDeleteRoute(route.route_id, route.path)}
                    className="p-1.5 rounded-lg text-stone-400 hover:text-red-600 hover:bg-red-50 transition-colors"
                    title="Delete route"
                  >
                    <Trash2 className="w-4 h-4" />
                  </button>
                </div>
              </div>
            </div>
          ))}
        </div>
      )}

      {/* Modal: Create / Edit Route */}
      {isModalOpen && (
        <div className="fixed inset-0 z-50 bg-stone-900/40 backdrop-blur-xs flex items-center justify-center p-4">
          <div className="bg-white rounded-2xl border border-stone-200 shadow-xl max-w-md w-full p-6 space-y-6">
            <div className="flex items-center justify-between pb-4 border-b border-stone-100">
              <h3 className="font-heading text-lg font-bold text-stone-900">
                {editingRouteId ? 'Edit Gateway Route' : 'Create Gateway Route'}
              </h3>
              <button onClick={() => setIsModalOpen(false)} className="text-stone-400 hover:text-stone-600">
                <X className="w-4 h-4" />
              </button>
            </div>

            <form onSubmit={handleSubmit} className="space-y-4">
              <div>
                <label className="block text-xs font-semibold text-stone-700 mb-1.5">
                  Route Path
                </label>
                <input
                  type="text"
                  required
                  value={path}
                  onChange={(e) => setPath(e.target.value)}
                  placeholder="e.g. /orders or /users/{id}"
                  className="input-field font-mono-url"
                />
              </div>

              <div className="grid grid-cols-2 gap-3">
                <div>
                  <label className="block text-xs font-semibold text-stone-700 mb-1.5">
                    Path Type
                  </label>
                  <select
                    value={pathType}
                    onChange={(e) => setPathType(e.target.value as PathType)}
                    className="input-field"
                  >
                    <option value="prefix">Prefix</option>
                    <option value="exact">Exact</option>
                    <option value="regex">Regex</option>
                  </select>
                </div>

                <div>
                  <label className="block text-xs font-semibold text-stone-700 mb-1.5">
                    HTTP Method
                  </label>
                  <select
                    value={method}
                    onChange={(e) => setMethod(e.target.value)}
                    className="input-field font-mono-url"
                  >
                    <option value="GET">GET</option>
                    <option value="POST">POST</option>
                    <option value="PUT">PUT</option>
                    <option value="DELETE">DELETE</option>
                    <option value="PATCH">PATCH</option>
                    <option value="*">ANY (*)</option>
                  </select>
                </div>
              </div>

              <div>
                <label className="block text-xs font-semibold text-stone-700 mb-1.5">
                  Target Upstream Name
                </label>
                {upstreams.length === 0 ? (
                  <div className="p-3 bg-amber-50 border border-amber-200 rounded-lg text-amber-800 text-xs">
                    No upstreams available. Please create an upstream first.
                  </div>
                ) : (
                  <select
                    value={upstreamName}
                    onChange={(e) => setUpstreamName(e.target.value)}
                    className="input-field font-mono-url"
                  >
                    {upstreams.map((u) => (
                      <option key={u.upstream_id} value={u.name}>
                        {u.name}
                      </option>
                    ))}
                  </select>
                )}
              </div>

              <div className="grid grid-cols-2 gap-3 pt-2">
                <div>
                  <label className="block text-xs font-semibold text-stone-700 mb-1.5">
                    Auth Mode
                  </label>
                  <select
                    value={authMode}
                    onChange={(e) => setAuthMode(e.target.value as AuthMode)}
                    className="input-field"
                  >
                    <option value="required">Required</option>
                    <option value="none">None (Public)</option>
                  </select>
                </div>

                <div>
                  <label className="block text-xs font-semibold text-stone-700 mb-1.5">
                    Route Status
                  </label>
                  <div className="flex items-center gap-2 pt-2">
                    <input
                      type="checkbox"
                      checked={enabled}
                      onChange={(e) => setEnabled(e.target.checked)}
                      className="w-5 h-5 accent-orange-600 cursor-pointer"
                    />
                    <span className="text-xs font-medium text-stone-700">Enabled</span>
                  </div>
                </div>
              </div>

              <div className="pt-4 flex items-center justify-end gap-3 border-t border-stone-100">
                <button
                  type="button"
                  onClick={() => setIsModalOpen(false)}
                  className="btn-white text-xs"
                >
                  Cancel
                </button>
                <button
                  type="submit"
                  disabled={upstreams.length === 0}
                  className="btn-orange-primary text-xs"
                >
                  {editingRouteId ? 'Save Changes' : 'Create Route'}
                </button>
              </div>
            </form>
          </div>
        </div>
      )}
    </div>
  );
};
