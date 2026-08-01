import React, { useState, useEffect } from 'react';
import { useParams, Link } from 'react-router-dom';
import {
  ArrowLeft,
  Server,
  Plus,
  Trash2,
  Edit2,
  X,
  AlertCircle,
  CheckCircle,
} from 'lucide-react';
import { api } from '../services/api';
import type {
  GetUpstreamResponseDTO,
  LoadBalancingStrategy,
  Backend,
  CreateOrUpdateUpstreamRequestDTO,
} from '../types/api';

export const UpstreamsPage: React.FC = () => {
  const { projectID } = useParams<{ projectID: string }>();

  const [upstreams, setUpstreams] = useState<GetUpstreamResponseDTO[]>([]);
  const [error, setError] = useState<string | null>(null);
  const [successMsg, setSuccessMsg] = useState<string | null>(null);

  // Modal / Form state
  const [isModalOpen, setIsModalOpen] = useState(false);
  const [editingUpstreamId, setEditingUpstreamId] = useState<string | null>(null);
  const [name, setName] = useState('');
  const [strategy, setStrategy] = useState<LoadBalancingStrategy>('ROUND_ROBIN');
  const [backends, setBackends] = useState<Backend[]>([{ url: 'http://localhost:8080' }]);

  const fetchUpstreams = async () => {
    if (!projectID) return;
    try {
      const res = await api.listUpstreams(projectID);
      setUpstreams(res.upstreams || []);
      setError(null);
    } catch (err: any) {
      setError(err.message || 'Failed to fetch upstreams.');
    }
  };

  useEffect(() => {
    fetchUpstreams();
  }, [projectID]);

  const openCreateModal = () => {
    setEditingUpstreamId(null);
    setName('');
    setStrategy('ROUND_ROBIN');
    setBackends([{ url: 'http://localhost:8080' }]);
    setIsModalOpen(true);
  };

  const openEditModal = (upstream: GetUpstreamResponseDTO) => {
    setEditingUpstreamId(upstream.upstream_id);
    setName(upstream.name);
    setStrategy(upstream.load_balancing_strategy);
    setBackends(upstream.backends?.length ? upstream.backends : [{ url: 'http://localhost:8080' }]);
    setIsModalOpen(true);
  };

  const handleAddBackend = () => {
    setBackends((prev) => [...prev, { url: 'http://localhost:8080' }]);
  };

  const handleRemoveBackend = (index: number) => {
    setBackends((prev) => prev.filter((_, i) => i !== index));
  };

  const handleBackendChange = (index: number, field: keyof Backend, value: any) => {
    setBackends((prev) =>
      prev.map((b, i) => (i === index ? { ...b, [field]: value } : b))
    );
  };

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    if (!projectID || !name.trim()) return;

    const payload: CreateOrUpdateUpstreamRequestDTO = {
      name: name.trim(),
      load_balancing_strategy: strategy,
      backends: backends.filter((b) => b.url.trim() !== ''),
    };

    try {
      if (editingUpstreamId) {
        await api.updateUpstream(projectID, editingUpstreamId, payload);
        setSuccessMsg(`Upstream '${name}' updated successfully.`);
      } else {
        await api.createUpstream(projectID, payload);
        setSuccessMsg(`Upstream '${name}' created successfully.`);
      }
      setIsModalOpen(false);
      setTimeout(() => setSuccessMsg(null), 3000);
      fetchUpstreams();
    } catch (err: any) {
      setError(err.message || 'Failed to save upstream.');
    }
  };

  const handleDeleteUpstream = async (upstreamId: string, upstreamName: string) => {
    if (!projectID) return;
    if (!window.confirm(`Are you sure you want to delete upstream '${upstreamName}'?`)) return;

    try {
      await api.deleteUpstream(projectID, upstreamId);
      setUpstreams((prev) => prev.filter((u) => u.upstream_id !== upstreamId));
      setSuccessMsg(`Upstream '${upstreamName}' deleted.`);
      setTimeout(() => setSuccessMsg(null), 3000);
    } catch (err: any) {
      setError(err.message || 'Failed to delete upstream.');
    }
  };

  return (
    <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-8 space-y-8">
      {/* Top Header */}
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
            <Server className="w-7 h-7 text-orange-600" />
            <span>Upstream Clusters</span>
          </h1>
          <p className="font-body text-sm text-stone-500 mt-1">
            Configure target backend servers and load balancing rules.
          </p>
        </div>
        <button onClick={openCreateModal} className="btn-orange-primary shadow-sm">
          <Plus className="w-4 h-4" />
          <span>New Upstream</span>
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

      {/* Upstreams Grid */}
      {upstreams.length === 0 ? (
        <div className="card-white p-12 text-center space-y-4">
          <div className="w-12 h-12 rounded-full bg-orange-50 flex items-center justify-center text-orange-500 mx-auto">
            <Server className="w-6 h-6" />
          </div>
          <div className="space-y-1">
            <h3 className="font-heading text-base font-semibold text-stone-800">
              No upstreams configured
            </h3>
            <p className="font-body text-xs text-stone-500 max-w-sm mx-auto">
              Create an upstream cluster pointing to your backend microservice URLs.
            </p>
          </div>
          <button onClick={openCreateModal} className="btn-orange-primary text-xs">
            <Plus className="w-3.5 h-3.5" />
            <span>Create Upstream</span>
          </button>
        </div>
      ) : (
        <div className="grid grid-cols-1 md:grid-cols-2 gap-6">
          {upstreams.map((upstream) => (
            <div key={upstream.upstream_id} className="card-white p-6 space-y-4 flex flex-col justify-between">
              <div className="space-y-3">
                <div className="flex items-start justify-between">
                  <div className="space-y-1">
                    <h3 className="font-heading text-lg font-bold text-stone-900">
                      {upstream.name}
                    </h3>
                    <span className="badge-orange font-mono-url text-[10px]">
                      {upstream.load_balancing_strategy}
                    </span>
                  </div>
                  <div className="flex items-center gap-1">
                    <button
                      onClick={() => openEditModal(upstream)}
                      className="p-1.5 rounded-lg text-stone-400 hover:text-stone-700 hover:bg-stone-100 transition-colors"
                      title="Edit upstream"
                    >
                      <Edit2 className="w-4 h-4" />
                    </button>
                    <button
                      onClick={() => handleDeleteUpstream(upstream.upstream_id, upstream.name)}
                      className="p-1.5 rounded-lg text-stone-400 hover:text-red-600 hover:bg-red-50 transition-colors"
                      title="Delete upstream"
                    >
                      <Trash2 className="w-4 h-4" />
                    </button>
                  </div>
                </div>

                <div className="space-y-2 pt-2">
                  <span className="text-[10px] font-semibold text-stone-400 uppercase tracking-wider block">
                    Target Backends ({upstream.backends?.length || 0})
                  </span>
                  <div className="space-y-1.5">
                    {upstream.backends?.map((backend, idx) => (
                      <div
                        key={idx}
                        className="flex items-center justify-between px-3 py-1.5 rounded-lg bg-stone-50 border border-stone-200 text-xs font-mono-url"
                      >
                        <span className="text-stone-800 font-medium truncate">{backend.url}</span>
                        {backend.weight !== undefined && backend.weight !== null && (
                          <span className="text-orange-600 text-[10px] font-semibold bg-orange-100 px-1.5 py-0.5 rounded">
                            Weight: {backend.weight}
                          </span>
                        )}
                      </div>
                    ))}
                  </div>
                </div>
              </div>

              <div className="pt-4 border-t border-stone-100 flex items-center justify-between text-[11px] text-stone-400">
                <span>ID: <code className="font-mono-url">{upstream.upstream_id}</code></span>
              </div>
            </div>
          ))}
        </div>
      )}

      {/* Modal: Create / Edit Upstream */}
      {isModalOpen && (
        <div className="fixed inset-0 z-50 bg-stone-900/40 backdrop-blur-xs flex items-center justify-center p-4">
          <div className="bg-white rounded-2xl border border-stone-200 shadow-xl max-w-lg w-full p-6 space-y-6 max-h-[90vh] overflow-y-auto">
            <div className="flex items-center justify-between pb-4 border-b border-stone-100">
              <h3 className="font-heading text-lg font-bold text-stone-900">
                {editingUpstreamId ? 'Edit Upstream' : 'Create Upstream Cluster'}
              </h3>
              <button onClick={() => setIsModalOpen(false)} className="text-stone-400 hover:text-stone-600">
                <X className="w-4 h-4" />
              </button>
            </div>

            <form onSubmit={handleSubmit} className="space-y-4">
              <div>
                <label className="block text-xs font-semibold text-stone-700 mb-1.5">
                  Upstream Name
                </label>
                <input
                  type="text"
                  required
                  value={name}
                  onChange={(e) => setName(e.target.value)}
                  placeholder="e.g. orders-backend"
                  className="input-field"
                />
              </div>

              <div>
                <label className="block text-xs font-semibold text-stone-700 mb-1.5">
                  Load Balancing Strategy
                </label>
                <select
                  value={strategy}
                  onChange={(e) => setStrategy(e.target.value as LoadBalancingStrategy)}
                  className="input-field"
                >
                  <option value="ROUND_ROBIN">Round Robin</option>
                  <option value="RANDOM">Random</option>
                  <option value="IP_HASH">IP Hash</option>
                  <option value="WEIGHTED_ROUND_ROBIN">Weighted Round Robin</option>
                  <option value="LEAST_CONNECTIONS">Least Connections</option>
                </select>
              </div>

              {/* Backends List */}
              <div className="space-y-2 pt-2">
                <div className="flex items-center justify-between">
                  <label className="block text-xs font-semibold text-stone-700">
                    Backend Target URLs
                  </label>
                  <button
                    type="button"
                    onClick={handleAddBackend}
                    className="text-xs text-orange-600 font-semibold hover:underline flex items-center gap-1"
                  >
                    <Plus className="w-3.5 h-3.5" />
                    <span>Add Target</span>
                  </button>
                </div>

                {backends.map((backend, index) => (
                  <div key={index} className="flex items-center gap-2">
                    <input
                      type="text"
                      required
                      value={backend.url}
                      onChange={(e) => handleBackendChange(index, 'url', e.target.value)}
                      placeholder="http://orders-1:8080"
                      className="input-field font-mono-url flex-1"
                    />

                    {strategy === 'WEIGHTED_ROUND_ROBIN' && (
                      <input
                        type="number"
                        placeholder="Weight"
                        value={backend.weight ?? 1}
                        onChange={(e) =>
                          handleBackendChange(index, 'weight', parseInt(e.target.value) || 1)
                        }
                        className="input-field w-20 font-mono-url"
                      />
                    )}

                    {backends.length > 1 && (
                      <button
                        type="button"
                        onClick={() => handleRemoveBackend(index)}
                        className="text-stone-400 hover:text-red-600 p-2"
                      >
                        <Trash2 className="w-4 h-4" />
                      </button>
                    )}
                  </div>
                ))}
              </div>

              <div className="pt-4 flex items-center justify-end gap-3 border-t border-stone-100">
                <button
                  type="button"
                  onClick={() => setIsModalOpen(false)}
                  className="btn-white text-xs"
                >
                  Cancel
                </button>
                <button type="submit" className="btn-orange-primary text-xs">
                  {editingUpstreamId ? 'Save Changes' : 'Create Upstream'}
                </button>
              </div>
            </form>
          </div>
        </div>
      )}
    </div>
  );
};
