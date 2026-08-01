import React, { useState, useEffect } from 'react';
import { Link, useNavigate } from 'react-router-dom';
import {
  FolderPlus,
  Search,
  Trash2,
  ExternalLink,
  Key,
  Layers,
  X,
  Plus,
  AlertCircle,
  Copy,
  Check,
  Lock,
  LogIn,
} from 'lucide-react';
import { api } from '../services/api';
import { useAuth } from '../context/AuthContext';
import type { ListProjectResponse, CreateProjectResponse } from '../types/api';

export const DashboardPage: React.FC = () => {
  const navigate = useNavigate();
  const { isLoggedIn, setIsLoggedIn } = useAuth();

  const [projects, setProjects] = useState<ListProjectResponse[]>([]);
  const [searchQuery, setSearchQuery] = useState('');
  const [error, setError] = useState<string | null>(null);

  // Create Project Modal state
  const [isModalOpen, setIsModalOpen] = useState(false);
  const [newProjectName, setNewProjectName] = useState('');
  const [isCreating, setIsCreating] = useState(false);

  // New Project Created Success Banner/Modal
  const [createdProject, setCreatedProject] = useState<CreateProjectResponse | null>(null);
  const [copiedKey, setCopiedKey] = useState(false);

  const fetchProjects = async () => {
    try {
      const res = await api.listProjects();
      setProjects(res.projects || []);
      setError(null);
    } catch (err: any) {
      if (err.message && (err.message.includes('401') || err.message.toLowerCase().includes('unauthorized'))) {
        setIsLoggedIn(false);
      }
      setError(err.message || 'Failed to load projects. Authentication may be required.');
    }
  };

  useEffect(() => {
    fetchProjects();
  }, []);

  const handleOpenCreateModal = () => {
    if (!isLoggedIn) {
      // If not logged in, redirect directly to login page with notice
      navigate('/login');
      return;
    }
    setIsModalOpen(true);
  };

  const handleCreateProject = async (e: React.FormEvent) => {
    e.preventDefault();
    if (!isLoggedIn) {
      setError('You must be logged in to create a project.');
      navigate('/login');
      return;
    }

    if (!newProjectName.trim()) return;

    setIsCreating(true);
    try {
      const res = await api.createProject({
        name: newProjectName.trim(),
        access_list: [],
      });
      setCreatedProject(res);
      setNewProjectName('');
      setIsModalOpen(false);
      fetchProjects();
    } catch (err: any) {
      setError(err.message || 'Failed to create project.');
    } finally {
      setIsCreating(false);
    }
  };

  const handleDeleteProject = async (projectId: string, e: React.MouseEvent) => {
    e.preventDefault();
    e.stopPropagation();

    if (!isLoggedIn) {
      setError('You must be logged in to delete a project.');
      navigate('/login');
      return;
    }

    if (!window.confirm('Are you sure you want to delete this project?')) return;

    try {
      await api.deleteProject(projectId);
      setProjects((prev) => prev.filter((p) => p.project_id !== projectId));
    } catch (err: any) {
      setError(err.message || 'Failed to delete project.');
    }
  };

  const copyKeyToClipboard = (key: string) => {
    navigator.clipboard.writeText(key);
    setCopiedKey(true);
    setTimeout(() => setCopiedKey(false), 2000);
  };

  const filteredProjects = projects.filter((p) =>
    p.name.toLowerCase().includes(searchQuery.toLowerCase()) ||
    p.project_id.toLowerCase().includes(searchQuery.toLowerCase())
  );

  return (
    <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-8 space-y-8">
      {/* Header */}
      <div className="flex flex-col sm:flex-row sm:items-center justify-between gap-4 pb-6 border-b border-stone-200">
        <div>
          <h1 className="font-heading text-3xl font-bold text-stone-900">Projects Dashboard</h1>
          <p className="font-body text-sm text-stone-500 mt-1">
            Manage your API gateway environments, middleware plugins, upstreams, and routes.
          </p>
        </div>
        <button onClick={handleOpenCreateModal} className="btn-orange-primary shadow-sm">
          {isLoggedIn ? <FolderPlus className="w-4 h-4" /> : <Lock className="w-4 h-4" />}
          <span>New Project</span>
        </button>
      </div>

      {/* Login Warning Banner if not logged in */}
      {!isLoggedIn && (
        <div className="flex flex-col sm:flex-row items-start sm:items-center justify-between gap-4 p-4 rounded-xl bg-orange-50 border border-orange-200 text-orange-900 text-sm">
          <div className="flex items-center gap-3">
            <Lock className="w-5 h-5 text-orange-600 shrink-0" />
            <div>
              <strong className="font-heading font-bold block text-stone-900">
                Authentication Required to Create Projects
              </strong>
              <span className="text-xs text-stone-600">
                Please log in or register an account to create and configure gateway projects.
              </span>
            </div>
          </div>
          <Link to="/login" className="btn-orange-primary text-xs shrink-0 py-2 px-4">
            <LogIn className="w-4 h-4" />
            <span>Log In Now</span>
          </Link>
        </div>
      )}

      {/* Error Alert */}
      {error && (
        <div className="flex items-center gap-3 p-4 rounded-xl bg-red-50 border border-red-200 text-red-700 text-sm">
          <AlertCircle className="w-5 h-5 shrink-0" />
          <div className="flex-1">{error}</div>
          <button onClick={() => setError(null)} className="text-red-500 hover:text-red-700">
            <X className="w-4 h-4" />
          </button>
        </div>
      )}

      {/* Success Modal when project is created */}
      {createdProject && (
        <div className="bg-orange-500/10 border border-orange-300 rounded-2xl p-6 relative">
          <button
            onClick={() => setCreatedProject(null)}
            className="absolute top-4 right-4 text-stone-400 hover:text-stone-700 cursor-pointer"
          >
            <X className="w-4 h-4" />
          </button>
          <div className="flex items-start gap-4">
            <div className="w-10 h-10 rounded-xl bg-orange-100 flex items-center justify-center text-orange-600 shrink-0">
              <Key className="w-5 h-5" />
            </div>
            <div className="space-y-2 flex-1">
              <h3 className="font-heading text-lg font-bold text-stone-900">
                Project Created Successfully!
              </h3>
              <p className="font-body text-xs text-stone-600">
                Save your project ID and API Gateway Key. You will need this key when forwarding requests to the gateway.
              </p>
              <div className="grid grid-cols-1 sm:grid-cols-2 gap-3 pt-2">
                <div className="bg-white p-3 rounded-lg border border-orange-200">
                  <span className="text-[10px] font-semibold text-stone-500 uppercase tracking-wider block mb-1">
                    Project ID
                  </span>
                  <code className="font-mono-url text-xs text-stone-900 select-all">
                    {createdProject.project_id}
                  </code>
                </div>
                <div className="bg-white p-3 rounded-lg border border-orange-200 relative group">
                  <span className="text-[10px] font-semibold text-stone-500 uppercase tracking-wider block mb-1">
                    API Gateway Key
                  </span>
                  <div className="flex items-center justify-between">
                    <code className="font-mono-url text-xs text-orange-600 font-semibold truncate pr-2 select-all">
                      {createdProject.api_gw_key}
                    </code>
                    <button
                      onClick={() => copyKeyToClipboard(createdProject.api_gw_key)}
                      className="text-stone-400 hover:text-orange-600"
                    >
                      {copiedKey ? <Check className="w-4 h-4 text-emerald-500" /> : <Copy className="w-4 h-4" />}
                    </button>
                  </div>
                </div>
              </div>
            </div>
          </div>
        </div>
      )}

      {/* Controls / Search Bar */}
      <div className="flex items-center gap-4 bg-white p-3 rounded-xl border border-stone-200 shadow-xs">
        <div className="relative flex-1">
          <Search className="w-4 h-4 text-stone-400 absolute left-3.5 top-1/2 -translate-y-1/2" />
          <input
            type="text"
            value={searchQuery}
            onChange={(e) => setSearchQuery(e.target.value)}
            placeholder="Search projects by name or ID..."
            className="input-field pl-10 border-0 bg-transparent focus:ring-0"
          />
        </div>
        <div className="text-xs text-stone-400 font-medium px-3 border-l border-stone-200 hidden sm:block">
          {filteredProjects.length} {filteredProjects.length === 1 ? 'project' : 'projects'}
        </div>
      </div>

      {/* Projects Grid */}
      {filteredProjects.length === 0 ? (
        <div className="card-white p-12 text-center space-y-4">
          <div className="w-12 h-12 rounded-full bg-stone-100 flex items-center justify-center text-stone-400 mx-auto">
            <Layers className="w-6 h-6" />
          </div>
          <div className="space-y-1">
            <h3 className="font-heading text-base font-semibold text-stone-800">
              No projects found
            </h3>
            <p className="font-body text-xs text-stone-500 max-w-sm mx-auto">
              {searchQuery
                ? 'Try matching your search query with another name or ID.'
                : 'Get started by creating your first API Gateway project.'}
            </p>
          </div>
          {!searchQuery && (
            <button onClick={handleOpenCreateModal} className="btn-orange-primary text-xs">
              <Plus className="w-3.5 h-3.5" />
              <span>Create Project</span>
            </button>
          )}
        </div>
      ) : (
        <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
          {filteredProjects.map((project) => (
            <Link
              key={project.project_id}
              to={`/dashboard/project/${project.project_id}`}
              className="card-white p-6 flex flex-col justify-between group hover:shadow-md hover:border-orange-300 transition-all"
            >
              <div className="space-y-3">
                <div className="flex items-start justify-between gap-3">
                  <div className="flex items-center gap-2.5">
                    <div className="w-9 h-9 rounded-lg bg-orange-50 border border-orange-200/60 flex items-center justify-center text-orange-600 group-hover:bg-orange-600 group-hover:text-white transition-colors">
                      <Layers className="w-4 h-4" />
                    </div>
                    <h3 className="font-heading text-lg font-bold text-stone-900 group-hover:text-orange-600 transition-colors">
                      {project.name}
                    </h3>
                  </div>
                  <button
                    onClick={(e) => handleDeleteProject(project.project_id, e)}
                    className="p-1.5 rounded-md text-stone-400 hover:text-red-600 hover:bg-red-50 transition-colors cursor-pointer"
                    title="Delete project"
                  >
                    <Trash2 className="w-4 h-4" />
                  </button>
                </div>

                <div className="space-y-1.5 pt-2">
                  <span className="text-[10px] font-semibold text-stone-400 uppercase tracking-wider block">
                    Project ID
                  </span>
                  <div className="font-mono-url text-xs text-stone-600 bg-stone-50 px-2.5 py-1 rounded border border-stone-200/60 truncate">
                    {project.project_id}
                  </div>
                </div>
              </div>

              <div className="pt-6 mt-6 border-t border-stone-100 flex items-center justify-between text-xs">
                <span className="text-stone-500 font-medium flex items-center gap-1 group-hover:text-orange-600">
                  <span>Configure Gateway</span>
                  <ExternalLink className="w-3.5 h-3.5" />
                </span>
                <span className="badge-orange">Active</span>
              </div>
            </Link>
          ))}
        </div>
      )}

      {/* Modal: Create New Project */}
      {isModalOpen && (
        <div className="fixed inset-0 z-50 bg-stone-900/40 backdrop-blur-xs flex items-center justify-center p-4">
          <div className="bg-white rounded-2xl border border-stone-200 shadow-xl max-w-md w-full p-6 space-y-6 animate-in fade-in zoom-in-95 duration-150">
            <div className="flex items-center justify-between pb-4 border-b border-stone-100">
              <h3 className="font-heading text-lg font-bold text-stone-900">Create New Project</h3>
              <button
                onClick={() => setIsModalOpen(false)}
                className="text-stone-400 hover:text-stone-600 cursor-pointer"
              >
                <X className="w-4 h-4" />
              </button>
            </div>

            {!isLoggedIn ? (
              <div className="space-y-4 text-center py-4">
                <div className="w-12 h-12 rounded-full bg-orange-100 text-orange-600 flex items-center justify-center mx-auto">
                  <Lock className="w-6 h-6" />
                </div>
                <h4 className="font-heading text-base font-bold text-stone-900">
                  Authentication Required
                </h4>
                <p className="font-body text-xs text-stone-500">
                  You must be logged in to create a project. Please sign in or create an account to proceed.
                </p>
                <div className="pt-4 flex items-center justify-center gap-3">
                  <button
                    type="button"
                    onClick={() => setIsModalOpen(false)}
                    className="btn-white text-xs"
                  >
                    Cancel
                  </button>
                  <Link to="/login" className="btn-orange-primary text-xs">
                    <LogIn className="w-3.5 h-3.5" />
                    <span>Go to Login</span>
                  </Link>
                </div>
              </div>
            ) : (
              <form onSubmit={handleCreateProject} className="space-y-4">
                <div>
                  <label className="block text-xs font-semibold text-stone-700 mb-1.5">
                    Project Name
                  </label>
                  <input
                    type="text"
                    required
                    value={newProjectName}
                    onChange={(e) => setNewProjectName(e.target.value)}
                    placeholder="e.g. Orders Microservice API"
                    className="input-field"
                  />
                </div>

                <div className="pt-4 flex items-center justify-end gap-3">
                  <button
                    type="button"
                    onClick={() => setIsModalOpen(false)}
                    className="btn-white text-xs"
                  >
                    Cancel
                  </button>
                  <button
                    type="submit"
                    disabled={isCreating || !newProjectName.trim()}
                    className="btn-orange-primary text-xs"
                  >
                    Create Project
                  </button>
                </div>
              </form>
            )}
          </div>
        </div>
      )}
    </div>
  );
};
