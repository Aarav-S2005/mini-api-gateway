import type {
  AuthRequestDTO,
  CreateProjectRequest,
  CreateProjectResponse,
  GetAllProjectResponse,
  GetProjectResponse,
  UpdateAccessListRequest,
  Middleware,
  CreateOrUpdateUpstreamRequestDTO,
  GetUpstreamResponseDTO,
  GetAllUpstreamResponseDTO,
  CreateOrUpdateRouteRequestDTO,
  GetRouteResponseDTO,
  GetAllRoutesResponseDTO,
} from '../types/api';

const API_BASE = (import.meta.env.VITE_API_BASE_URL as string) || '/api';
const cleanBaseUrl = API_BASE.replace(/\/+$/, '');

let activeLoaderTrigger: ((promise: Promise<any>) => Promise<any>) | null = null;

export const setApiLoaderTrigger = (trigger: <T>(promise: Promise<T>) => Promise<T>) => {
  activeLoaderTrigger = trigger;
};

async function customFetch<T>(endpoint: string, options: RequestInit = {}): Promise<T> {
  const fetchPromise = (async () => {
    const defaultHeaders: Record<string, string> = {
      'Content-Type': 'application/json',
    };

    const config: RequestInit = {
      ...options,
      credentials: 'include', // Ensure cookies are sent and received
      headers: {
        ...defaultHeaders,
        ...(options.headers || {}),
      },
    };

    const url = `${cleanBaseUrl}${endpoint.startsWith('/') ? endpoint : `/${endpoint}`}`;
    const response = await fetch(url, config);

    if (response.status === 204) {
      return null as T;
    }

    const contentType = response.headers.get('content-type');
    let data: any = null;

    if (contentType && contentType.includes('application/json')) {
      data = await response.json();
    } else {
      const text = await response.text();
      data = text ? { message: text } : null;
    }

    if (!response.ok) {
      const errorMessage = data?.message || data?.error || `Request failed with status ${response.status}`;
      throw new Error(errorMessage);
    }

    return data as T;
  })();

  if (activeLoaderTrigger) {
    return activeLoaderTrigger(fetchPromise);
  }
  return fetchPromise;
}

export const api = {
  // Auth
  register: (data: AuthRequestDTO) =>
    customFetch<null>('/auth/register', {
      method: 'POST',
      body: JSON.stringify(data),
    }),

  login: (data: AuthRequestDTO) =>
    customFetch<null>('/auth/login', {
      method: 'POST',
      body: JSON.stringify(data),
    }),

  logout: () =>
    customFetch<null>('/auth/logout', {
      method: 'GET',
    }),

  // Projects
  createProject: (data: CreateProjectRequest) =>
    customFetch<CreateProjectResponse>('/projects', {
      method: 'POST',
      body: JSON.stringify(data),
    }),

  listProjects: () =>
    customFetch<GetAllProjectResponse>('/projects', {
      method: 'GET',
    }),

  getProject: (projectId: string) =>
    customFetch<GetProjectResponse>(`/projects/${projectId}`, {
      method: 'GET',
    }),

  deleteProject: (projectId: string) =>
    customFetch<null>(`/projects/${projectId}`, {
      method: 'DELETE',
    }),

  updateAccessList: (projectId: string, data: UpdateAccessListRequest) =>
    customFetch<null>(`/projects/${projectId}/accesslist`, {
      method: 'PATCH',
      body: JSON.stringify(data),
    }),

  // Middlewares
  updateMiddlewares: (projectId: string, middlewares: Middleware[]) =>
    customFetch<null>(`/projects/${projectId}/middlewares`, {
      method: 'PATCH',
      body: JSON.stringify({ middlewares }),
    }),

  deleteMiddleware: (projectId: string, name: string) =>
    customFetch<null>(`/projects/${projectId}/middlewares/${name}`, {
      method: 'DELETE',
    }),

  // Upstreams
  createUpstream: (projectId: string, data: CreateOrUpdateUpstreamRequestDTO) =>
    customFetch<{ upstream_id: string }>(`/projects/${projectId}/upstreams`, {
      method: 'POST',
      body: JSON.stringify(data),
    }),

  listUpstreams: (projectId: string) =>
    customFetch<GetAllUpstreamResponseDTO>(`/projects/${projectId}/upstreams`, {
      method: 'GET',
    }),

  getUpstream: (projectId: string, upstreamId: string) =>
    customFetch<GetUpstreamResponseDTO>(`/projects/${projectId}/upstreams/${upstreamId}`, {
      method: 'GET',
    }),

  updateUpstream: (projectId: string, upstreamId: string, data: CreateOrUpdateUpstreamRequestDTO) =>
    customFetch<null>(`/projects/${projectId}/upstreams/${upstreamId}`, {
      method: 'PUT',
      body: JSON.stringify(data),
    }),

  deleteUpstream: (projectId: string, upstreamId: string) =>
    customFetch<null>(`/projects/${projectId}/upstreams/${upstreamId}`, {
      method: 'DELETE',
    }),

  // Routes
  createRoute: (projectId: string, data: CreateOrUpdateRouteRequestDTO) =>
    customFetch<{ route_id: string }>(`/projects/${projectId}/routes`, {
      method: 'POST',
      body: JSON.stringify(data),
    }),

  listRoutes: (projectId: string) =>
    customFetch<GetAllRoutesResponseDTO>(`/projects/${projectId}/routes`, {
      method: 'GET',
    }),

  getRoute: (projectId: string, routeId: string) =>
    customFetch<GetRouteResponseDTO>(`/projects/${projectId}/routes/${routeId}`, {
      method: 'GET',
    }),

  updateRoute: (projectId: string, routeId: string, data: CreateOrUpdateRouteRequestDTO) =>
    customFetch<null>(`/projects/${projectId}/routes/${routeId}`, {
      method: 'PUT',
      body: JSON.stringify(data),
    }),

  deleteRoute: (projectId: string, routeId: string) =>
    customFetch<null>(`/projects/${projectId}/routes/${routeId}`, {
      method: 'DELETE',
    }),
};
