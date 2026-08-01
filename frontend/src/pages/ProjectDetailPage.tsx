import React, { useState, useEffect } from 'react';
import { useParams, Link } from 'react-router-dom';
import {
  ArrowLeft,
  Server,
  Route as RouteIcon,
  Sliders,
  Users,
  Trash2,
  Save,
  Plus,
  AlertCircle,
  CheckCircle,
  Globe,
  Filter,
  Clock,
  Key,
} from 'lucide-react';
import { api } from '../services/api';
import type { GetProjectResponse, Middleware, Access, Permission } from '../types/api';

export const ProjectDetailPage: React.FC = () => {
  const { projectID } = useParams<{ projectID: string }>();

  const [project, setProject] = useState<GetProjectResponse | null>(null);
  const [error, setError] = useState<string | null>(null);
  const [successMsg, setSuccessMsg] = useState<string | null>(null);

  // Middlewares state
  const [middlewares, setMiddlewares] = useState<Middleware[]>([]);
  const [activePluginTab, setActivePluginTab] = useState<'cors' | 'ip-filter' | 'rate-limit' | 'jwt-auth'>('cors');

  // Access List state
  const [accessList, setAccessList] = useState<Access[]>([]);
  const [newUsername, setNewUsername] = useState('');
  const [newPermission, setNewPermission] = useState<Permission>('viewing');

  const fetchProjectDetails = async () => {
    if (!projectID) return;
    try {
      const data = await api.getProject(projectID);
      setProject(data);
      setMiddlewares(data.middlewares || []);
      setError(null);
    } catch (err: any) {
      setError(err.message || 'Failed to fetch project details.');
    }
  };

  useEffect(() => {
    fetchProjectDetails();
  }, [projectID]);

  // Helper to get active middleware config or default empty object
  const getPluginConfig = (pluginName: string) => {
    const found = middlewares.find((m) => m.name === pluginName);
    return found ? found.config : null;
  };

  // Helper to update a plugin's config in local state
  const updateLocalMiddlewareConfig = (pluginName: string, newConfig: Record<string, any>) => {
    setMiddlewares((prev) => {
      const exists = prev.some((m) => m.name === pluginName);
      if (exists) {
        return prev.map((m) => (m.name === pluginName ? { ...m, config: newConfig } : m));
      } else {
        return [...prev, { name: pluginName, config: newConfig }];
      }
    });
  };

  const handleSaveMiddlewares = async () => {
    if (!projectID) return;
    try {
      await api.updateMiddlewares(projectID, middlewares);
      setSuccessMsg('Middleware configurations updated successfully!');
      setTimeout(() => setSuccessMsg(null), 3000);
      fetchProjectDetails();
    } catch (err: any) {
      setError(err.message || 'Failed to update middlewares.');
    }
  };

  const handleDeleteMiddleware = async (pluginName: string) => {
    if (!projectID) return;
    if (!window.confirm(`Delete configuration for plugin '${pluginName}'?`)) return;

    try {
      await api.deleteMiddleware(projectID, pluginName);
      setMiddlewares((prev) => prev.filter((m) => m.name !== pluginName));
      setSuccessMsg(`Middleware '${pluginName}' deleted.`);
      setTimeout(() => setSuccessMsg(null), 3000);
    } catch (err: any) {
      setError(err.message || 'Failed to delete middleware.');
    }
  };

  const handleSaveAccessList = async (updatedList: Access[]) => {
    if (!projectID) return;
    try {
      await api.updateAccessList(projectID, { access_list: updatedList });
      setAccessList(updatedList);
      setSuccessMsg('Access list updated.');
      setTimeout(() => setSuccessMsg(null), 3000);
    } catch (err: any) {
      setError(err.message || 'Failed to update access list.');
    }
  };

  const handleAddAccessUser = (e: React.FormEvent) => {
    e.preventDefault();
    if (!newUsername.trim()) return;

    const updated = [...accessList, { username: newUsername.trim(), permission: newPermission }];
    handleSaveAccessList(updated);
    setNewUsername('');
  };

  const handleRemoveAccessUser = (username: string) => {
    const updated = accessList.filter((a) => a.username !== username);
    handleSaveAccessList(updated);
  };

  if (!project) {
    return (
      <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-12 text-center text-stone-500 font-body">
        Loading project metadata...
      </div>
    );
  }

  const corsConfig = getPluginConfig('cors') || {
    enabled: false,
    allowed_origins: ['*'],
    allowed_methods: ['GET', 'POST', 'PUT', 'DELETE'],
    allowed_headers: ['Content-Type', 'Authorization'],
  };

  const ipFilterConfig = getPluginConfig('ip-filter') || {
    enabled: false,
    black_listed_ips: [],
    white_listed_ips: [],
  };

  const rateLimitConfig = getPluginConfig('rate-limit') || {
    enabled: false,
    strategy: 'token_bucket',
    key_by: 'ip',
    token_bucket: { capacity: 100, refill_rate: 10 },
    fixed_window: { limit: 1000, window_seconds: 60 },
  };

  const jwtConfig = getPluginConfig('jwt-auth') || {
    enabled: false,
    algorithm: 'RS256',
    public_key: '',
    user_id_claim: 'sub',
    token_source: 'header',
    header_name: 'Authorization',
    prefix: 'Bearer',
    cookie_name: 'auth-token',
  };

  return (
    <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-8 space-y-8">
      {/* Top Navigation / Breadcrumbs */}
      <div className="flex items-center justify-between">
        <Link
          to="/dashboard"
          className="inline-flex items-center gap-1.5 text-xs font-semibold text-stone-500 hover:text-orange-600 transition-colors"
        >
          <ArrowLeft className="w-4 h-4" />
          <span>Back to Projects</span>
        </Link>
        <span className="badge-orange uppercase tracking-wider text-[10px]">
          Permission: {project.permission || 'editing'}
        </span>
      </div>

      {/* Project Overview Header */}
      <div className="card-white p-6 sm:p-8 space-y-6">
        <div className="flex flex-col sm:flex-row sm:items-center justify-between gap-4">
          <div>
            <h1 className="font-heading text-3xl font-bold text-stone-900">{project.name}</h1>
            <div className="flex items-center gap-2 mt-2">
              <span className="text-xs text-stone-400 font-semibold uppercase tracking-wider">
                Project ID:
              </span>
              <code className="font-mono-url text-xs font-medium text-stone-800 bg-stone-100 px-2 py-0.5 rounded border border-stone-200">
                {project.project_id}
              </code>
            </div>
          </div>

          {/* Quick Sub-route Buttons */}
          <div className="flex items-center gap-3">
            <Link
              to={`/dashboard/project/${projectID}/upstreams`}
              className="btn-white text-xs py-2 px-4"
            >
              <Server className="w-4 h-4 text-orange-600" />
              <span>Upstreams</span>
            </Link>

            <Link
              to={`/dashboard/project/${projectID}/routes`}
              className="btn-white text-xs py-2 px-4"
            >
              <RouteIcon className="w-4 h-4 text-orange-600" />
              <span>Routes</span>
            </Link>
          </div>
        </div>
      </div>

      {/* Alerts */}
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

      {/* Plugins / Middlewares Editor */}
      <div className="card-white p-6 sm:p-8 space-y-6">
        <div className="flex flex-col sm:flex-row sm:items-center justify-between gap-4 pb-4 border-b border-stone-100">
          <div>
            <h2 className="font-heading text-xl font-bold text-stone-900 flex items-center gap-2">
              <Sliders className="w-5 h-5 text-orange-600" />
              <span>Project Middlewares & Plugins</span>
            </h2>
            <p className="font-body text-xs text-stone-500 mt-1">
              Configure edge request handling for CORS, IP filtering, rate limiting, and JWT authentication.
            </p>
          </div>
          <button onClick={handleSaveMiddlewares} className="btn-orange-primary text-xs shadow-sm">
            <Save className="w-4 h-4" />
            <span>Save Middlewares</span>
          </button>
        </div>

        {/* Plugin Tabs */}
        <div className="flex border-b border-stone-200 overflow-x-auto gap-2">
          {[
            { id: 'cors', label: 'CORS', icon: Globe },
            { id: 'ip-filter', label: 'IP Filter', icon: Filter },
            { id: 'rate-limit', label: 'Rate Limit', icon: Clock },
            { id: 'jwt-auth', label: 'JWT Auth', icon: Key },
          ].map((tab) => {
            const Icon = tab.icon;
            const isEnabled = getPluginConfig(tab.id)?.enabled;
            return (
              <button
                key={tab.id}
                onClick={() => setActivePluginTab(tab.id as any)}
                className={`flex items-center gap-2 px-4 py-2.5 text-xs font-semibold border-b-2 transition-all cursor-pointer whitespace-nowrap ${
                  activePluginTab === tab.id
                    ? 'border-orange-600 text-orange-600 bg-orange-50/40'
                    : 'border-transparent text-stone-500 hover:text-stone-800'
                }`}
              >
                <Icon className="w-4 h-4" />
                <span>{tab.label}</span>
                {isEnabled && <span className="w-2 h-2 rounded-full bg-emerald-500" />}
              </button>
            );
          })}
        </div>

        {/* Plugin Editor Form Panels */}
        <div className="pt-2">
          {/* CORS Panel */}
          {activePluginTab === 'cors' && (
            <div className="space-y-6">
              <div className="flex items-center justify-between p-4 rounded-xl bg-stone-50 border border-stone-200">
                <div>
                  <h4 className="font-heading font-bold text-sm text-stone-900">Enable CORS Plugin</h4>
                  <p className="font-body text-xs text-stone-500">
                    Inject Cross-Origin Resource Sharing headers into responses.
                  </p>
                </div>
                <input
                  type="checkbox"
                  checked={corsConfig.enabled}
                  onChange={(e) =>
                    updateLocalMiddlewareConfig('cors', { ...corsConfig, enabled: e.target.checked })
                  }
                  className="w-5 h-5 accent-orange-600 cursor-pointer"
                />
              </div>

              <div className="grid grid-cols-1 gap-4">
                <div>
                  <label className="block text-xs font-semibold text-stone-700 mb-1">
                    Allowed Origins (comma separated)
                  </label>
                  <input
                    type="text"
                    value={corsConfig.allowed_origins.join(', ')}
                    onChange={(e) =>
                      updateLocalMiddlewareConfig('cors', {
                        ...corsConfig,
                        allowed_origins: e.target.value.split(',').map((s) => s.trim()).filter(Boolean),
                      })
                    }
                    placeholder="https://app.example.com, *"
                    className="input-field font-mono-url"
                  />
                </div>

                <div>
                  <label className="block text-xs font-semibold text-stone-700 mb-1">
                    Allowed Methods (comma separated)
                  </label>
                  <input
                    type="text"
                    value={corsConfig.allowed_methods.join(', ')}
                    onChange={(e) =>
                      updateLocalMiddlewareConfig('cors', {
                        ...corsConfig,
                        allowed_methods: e.target.value.split(',').map((s) => s.trim()).filter(Boolean),
                      })
                    }
                    placeholder="GET, POST, PUT, DELETE, OPTIONS"
                    className="input-field font-mono-url"
                  />
                </div>

                <div>
                  <label className="block text-xs font-semibold text-stone-700 mb-1">
                    Allowed Headers (comma separated)
                  </label>
                  <input
                    type="text"
                    value={corsConfig.allowed_headers.join(', ')}
                    onChange={(e) =>
                      updateLocalMiddlewareConfig('cors', {
                        ...corsConfig,
                        allowed_headers: e.target.value.split(',').map((s) => s.trim()).filter(Boolean),
                      })
                    }
                    placeholder="Content-Type, Authorization"
                    className="input-field font-mono-url"
                  />
                </div>
              </div>

              <div className="flex justify-end">
                <button
                  onClick={() => handleDeleteMiddleware('cors')}
                  className="btn-danger text-xs"
                >
                  <Trash2 className="w-3.5 h-3.5" />
                  <span>Delete CORS Plugin</span>
                </button>
              </div>
            </div>
          )}

          {/* IP Filter Panel */}
          {activePluginTab === 'ip-filter' && (
            <div className="space-y-6">
              <div className="flex items-center justify-between p-4 rounded-xl bg-stone-50 border border-stone-200">
                <div>
                  <h4 className="font-heading font-bold text-sm text-stone-900">Enable IP Filter</h4>
                  <p className="font-body text-xs text-stone-500">
                    Allow or block requests based on client IP addresses.
                  </p>
                </div>
                <input
                  type="checkbox"
                  checked={ipFilterConfig.enabled}
                  onChange={(e) =>
                    updateLocalMiddlewareConfig('ip-filter', { ...ipFilterConfig, enabled: e.target.checked })
                  }
                  className="w-5 h-5 accent-orange-600 cursor-pointer"
                />
              </div>

              <div className="grid grid-cols-1 sm:grid-cols-2 gap-4">
                <div>
                  <label className="block text-xs font-semibold text-stone-700 mb-1">
                    Blacklisted IPs (comma separated)
                  </label>
                  <textarea
                    rows={4}
                    value={ipFilterConfig.black_listed_ips.join(', ')}
                    onChange={(e) =>
                      updateLocalMiddlewareConfig('ip-filter', {
                        ...ipFilterConfig,
                        black_listed_ips: e.target.value.split(',').map((s) => s.trim()).filter(Boolean),
                      })
                    }
                    placeholder="203.0.113.10, 2001:db8::10"
                    className="input-field font-mono-url"
                  />
                </div>

                <div>
                  <label className="block text-xs font-semibold text-stone-700 mb-1">
                    Whitelisted IPs (comma separated)
                  </label>
                  <textarea
                    rows={4}
                    value={ipFilterConfig.white_listed_ips.join(', ')}
                    onChange={(e) =>
                      updateLocalMiddlewareConfig('ip-filter', {
                        ...ipFilterConfig,
                        white_listed_ips: e.target.value.split(',').map((s) => s.trim()).filter(Boolean),
                      })
                    }
                    placeholder="203.0.113.20, 2001:db8::20"
                    className="input-field font-mono-url"
                  />
                </div>
              </div>

              <div className="flex justify-end">
                <button
                  onClick={() => handleDeleteMiddleware('ip-filter')}
                  className="btn-danger text-xs"
                >
                  <Trash2 className="w-3.5 h-3.5" />
                  <span>Delete IP Filter Plugin</span>
                </button>
              </div>
            </div>
          )}

          {/* Rate Limit Panel */}
          {activePluginTab === 'rate-limit' && (
            <div className="space-y-6">
              <div className="flex items-center justify-between p-4 rounded-xl bg-stone-50 border border-stone-200">
                <div>
                  <h4 className="font-heading font-bold text-sm text-stone-900">Enable Rate Limiting</h4>
                  <p className="font-body text-xs text-stone-500">
                    Throttle excessive requests from clients.
                  </p>
                </div>
                <input
                  type="checkbox"
                  checked={rateLimitConfig.enabled}
                  onChange={(e) =>
                    updateLocalMiddlewareConfig('rate-limit', { ...rateLimitConfig, enabled: e.target.checked })
                  }
                  className="w-5 h-5 accent-orange-600 cursor-pointer"
                />
              </div>

              <div className="grid grid-cols-1 sm:grid-cols-2 gap-4">
                <div>
                  <label className="block text-xs font-semibold text-stone-700 mb-1">Strategy</label>
                  <select
                    value={rateLimitConfig.strategy}
                    onChange={(e) =>
                      updateLocalMiddlewareConfig('rate-limit', {
                        ...rateLimitConfig,
                        strategy: e.target.value as any,
                      })
                    }
                    className="input-field"
                  >
                    <option value="token_bucket">Token Bucket</option>
                    <option value="fixed_window">Fixed Window</option>
                  </select>
                </div>

                <div>
                  <label className="block text-xs font-semibold text-stone-700 mb-1">Key By</label>
                  <select
                    value={rateLimitConfig.key_by || 'ip'}
                    onChange={(e) =>
                      updateLocalMiddlewareConfig('rate-limit', {
                        ...rateLimitConfig,
                        key_by: e.target.value as any,
                      })
                    }
                    className="input-field"
                  >
                    <option value="ip">Client IP Address</option>
                    <option value="api_key">API Key</option>
                  </select>
                </div>
              </div>

              {rateLimitConfig.strategy === 'token_bucket' ? (
                <div className="grid grid-cols-1 sm:grid-cols-2 gap-4 p-4 rounded-xl bg-orange-50/30 border border-orange-200/60">
                  <div>
                    <label className="block text-xs font-semibold text-stone-700 mb-1">Capacity</label>
                    <input
                      type="number"
                      value={rateLimitConfig.token_bucket?.capacity || 100}
                      onChange={(e) =>
                        updateLocalMiddlewareConfig('rate-limit', {
                          ...rateLimitConfig,
                          token_bucket: {
                            ...(rateLimitConfig.token_bucket || { refill_rate: 10 }),
                            capacity: parseInt(e.target.value) || 0,
                          },
                        })
                      }
                      className="input-field font-mono-url"
                    />
                  </div>
                  <div>
                    <label className="block text-xs font-semibold text-stone-700 mb-1">Refill Rate (tokens/sec)</label>
                    <input
                      type="number"
                      value={rateLimitConfig.token_bucket?.refill_rate || 10}
                      onChange={(e) =>
                        updateLocalMiddlewareConfig('rate-limit', {
                          ...rateLimitConfig,
                          token_bucket: {
                            ...(rateLimitConfig.token_bucket || { capacity: 100 }),
                            refill_rate: parseInt(e.target.value) || 0,
                          },
                        })
                      }
                      className="input-field font-mono-url"
                    />
                  </div>
                </div>
              ) : (
                <div className="grid grid-cols-1 sm:grid-cols-2 gap-4 p-4 rounded-xl bg-orange-50/30 border border-orange-200/60">
                  <div>
                    <label className="block text-xs font-semibold text-stone-700 mb-1">Limit (max requests)</label>
                    <input
                      type="number"
                      value={rateLimitConfig.fixed_window?.limit || 1000}
                      onChange={(e) =>
                        updateLocalMiddlewareConfig('rate-limit', {
                          ...rateLimitConfig,
                          fixed_window: {
                            ...(rateLimitConfig.fixed_window || { window_seconds: 60 }),
                            limit: parseInt(e.target.value) || 0,
                          },
                        })
                      }
                      className="input-field font-mono-url"
                    />
                  </div>
                  <div>
                    <label className="block text-xs font-semibold text-stone-700 mb-1">Window (seconds)</label>
                    <input
                      type="number"
                      value={rateLimitConfig.fixed_window?.window_seconds || 60}
                      onChange={(e) =>
                        updateLocalMiddlewareConfig('rate-limit', {
                          ...rateLimitConfig,
                          fixed_window: {
                            ...(rateLimitConfig.fixed_window || { limit: 1000 }),
                            window_seconds: parseInt(e.target.value) || 0,
                          },
                        })
                      }
                      className="input-field font-mono-url"
                    />
                  </div>
                </div>
              )}

              <div className="flex justify-end">
                <button
                  onClick={() => handleDeleteMiddleware('rate-limit')}
                  className="btn-danger text-xs"
                >
                  <Trash2 className="w-3.5 h-3.5" />
                  <span>Delete Rate Limit Plugin</span>
                </button>
              </div>
            </div>
          )}

          {/* JWT Auth Panel */}
          {activePluginTab === 'jwt-auth' && (
            <div className="space-y-6">
              <div className="flex items-center justify-between p-4 rounded-xl bg-stone-50 border border-stone-200">
                <div>
                  <h4 className="font-heading font-bold text-sm text-stone-900">Enable JWT Auth</h4>
                  <p className="font-body text-xs text-stone-500">
                    Verify signed JWT tokens before forwarding requests downstream.
                  </p>
                </div>
                <input
                  type="checkbox"
                  checked={jwtConfig.enabled}
                  onChange={(e) =>
                    updateLocalMiddlewareConfig('jwt-auth', { ...jwtConfig, enabled: e.target.checked })
                  }
                  className="w-5 h-5 accent-orange-600 cursor-pointer"
                />
              </div>

              <div className="grid grid-cols-1 sm:grid-cols-3 gap-4">
                <div>
                  <label className="block text-xs font-semibold text-stone-700 mb-1">Algorithm</label>
                  <select
                    value={jwtConfig.algorithm}
                    onChange={(e) =>
                      updateLocalMiddlewareConfig('jwt-auth', { ...jwtConfig, algorithm: e.target.value })
                    }
                    className="input-field"
                  >
                    <option value="RS256">RS256</option>
                    <option value="ES256">ES256</option>
                  </select>
                </div>

                <div>
                  <label className="block text-xs font-semibold text-stone-700 mb-1">User ID Claim</label>
                  <input
                    type="text"
                    value={jwtConfig.user_id_claim}
                    onChange={(e) =>
                      updateLocalMiddlewareConfig('jwt-auth', { ...jwtConfig, user_id_claim: e.target.value })
                    }
                    placeholder="sub or user_id"
                    className="input-field font-mono-url"
                  />
                </div>

                <div>
                  <label className="block text-xs font-semibold text-stone-700 mb-1">Token Source</label>
                  <select
                    value={jwtConfig.token_source}
                    onChange={(e) =>
                      updateLocalMiddlewareConfig('jwt-auth', {
                        ...jwtConfig,
                        token_source: e.target.value as any,
                      })
                    }
                    className="input-field"
                  >
                    <option value="header">HTTP Header</option>
                    <option value="cookie">Cookie</option>
                  </select>
                </div>
              </div>

              {jwtConfig.token_source === 'header' ? (
                <div className="grid grid-cols-1 sm:grid-cols-2 gap-4">
                  <div>
                    <label className="block text-xs font-semibold text-stone-700 mb-1">Header Name</label>
                    <input
                      type="text"
                      value={jwtConfig.header_name || 'Authorization'}
                      onChange={(e) =>
                        updateLocalMiddlewareConfig('jwt-auth', { ...jwtConfig, header_name: e.target.value })
                      }
                      placeholder="Authorization"
                      className="input-field font-mono-url"
                    />
                  </div>
                  <div>
                    <label className="block text-xs font-semibold text-stone-700 mb-1">Prefix</label>
                    <input
                      type="text"
                      value={jwtConfig.prefix || 'Bearer'}
                      onChange={(e) =>
                        updateLocalMiddlewareConfig('jwt-auth', { ...jwtConfig, prefix: e.target.value })
                      }
                      placeholder="Bearer"
                      className="input-field font-mono-url"
                    />
                  </div>
                </div>
              ) : (
                <div>
                  <label className="block text-xs font-semibold text-stone-700 mb-1">Cookie Name</label>
                  <input
                    type="text"
                    value={jwtConfig.cookie_name || 'auth-token'}
                    onChange={(e) =>
                      updateLocalMiddlewareConfig('jwt-auth', { ...jwtConfig, cookie_name: e.target.value })
                    }
                    placeholder="auth-token"
                    className="input-field font-mono-url"
                  />
                </div>
              )}

              <div>
                <label className="block text-xs font-semibold text-stone-700 mb-1">Public Key (PEM)</label>
                <textarea
                  rows={4}
                  value={jwtConfig.public_key}
                  onChange={(e) =>
                    updateLocalMiddlewareConfig('jwt-auth', { ...jwtConfig, public_key: e.target.value })
                  }
                  placeholder="-----BEGIN PUBLIC KEY-----\n...\n-----END PUBLIC KEY-----"
                  className="input-field font-mono-url"
                />
              </div>

              <div className="flex justify-end">
                <button
                  onClick={() => handleDeleteMiddleware('jwt-auth')}
                  className="btn-danger text-xs"
                >
                  <Trash2 className="w-3.5 h-3.5" />
                  <span>Delete JWT Auth Plugin</span>
                </button>
              </div>
            </div>
          )}
        </div>
      </div>

      {/* Access List Management Section */}
      <div className="card-white p-6 sm:p-8 space-y-6">
        <div className="pb-4 border-b border-stone-100">
          <h2 className="font-heading text-xl font-bold text-stone-900 flex items-center gap-2">
            <Users className="w-5 h-5 text-orange-600" />
            <span>Project Access List</span>
          </h2>
          <p className="font-body text-xs text-stone-500 mt-1">
            Grant permission levels (<code className="font-mono-url">editing</code> or <code className="font-mono-url">viewing</code>) to team members.
          </p>
        </div>

        {/* Add User Form */}
        <form onSubmit={handleAddAccessUser} className="flex flex-col sm:flex-row gap-3">
          <input
            type="text"
            required
            value={newUsername}
            onChange={(e) => setNewUsername(e.target.value)}
            placeholder="Teammate username"
            className="input-field flex-1"
          />
          <select
            value={newPermission}
            onChange={(e) => setNewPermission(e.target.value as Permission)}
            className="input-field w-full sm:w-40"
          >
            <option value="editing">Editing</option>
            <option value="viewing">Viewing</option>
          </select>
          <button type="submit" className="btn-orange-primary text-xs shrink-0">
            <Plus className="w-4 h-4" />
            <span>Add Member</span>
          </button>
        </form>

        {/* User Table */}
        {accessList.length > 0 && (
          <div className="overflow-x-auto border border-stone-200 rounded-xl">
            <table className="w-full text-left text-xs">
              <thead className="bg-stone-50 border-b border-stone-200 font-heading text-stone-600">
                <tr>
                  <th className="px-4 py-3 font-semibold">Username</th>
                  <th className="px-4 py-3 font-semibold">Permission</th>
                  <th className="px-4 py-3 font-semibold text-right">Actions</th>
                </tr>
              </thead>
              <tbody className="divide-y divide-stone-100 font-body">
                {accessList.map((user) => (
                  <tr key={user.username} className="hover:bg-stone-50/60">
                    <td className="px-4 py-3 font-medium text-stone-900">{user.username}</td>
                    <td className="px-4 py-3">
                      <span
                        className={
                          user.permission === 'editing' ? 'badge-orange' : 'badge-gray'
                        }
                      >
                        {user.permission}
                      </span>
                    </td>
                    <td className="px-4 py-3 text-right">
                      <button
                        onClick={() => handleRemoveAccessUser(user.username)}
                        className="text-stone-400 hover:text-red-600 transition-colors"
                        title="Remove user"
                      >
                        <Trash2 className="w-4 h-4 inline" />
                      </button>
                    </td>
                  </tr>
                ))}
              </tbody>
            </table>
          </div>
        )}
      </div>
    </div>
  );
};
