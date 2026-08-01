export type Permission = 'editing' | 'viewing';

export interface Access {
  username: string;
  permission: Permission;
}

export interface AuthRequestDTO {
  username: string;
  password: string;
}

export interface CreateProjectRequest {
  name: string;
  access_list: Access[];
}

export interface CreateProjectResponse {
  project_id: string;
  api_gw_key: string;
}

export interface ListProjectResponse {
  project_id: string;
  name: string;
}

export interface GetAllProjectResponse {
  projects: ListProjectResponse[];
}

export interface GetProjectResponse {
  project_id: string;
  name: string;
  middlewares: Middleware[];
  permission: string;
}

export interface UpdateAccessListRequest {
  access_list: Access[];
}

export interface Middleware {
  name: string;
  config: Record<string, any>;
}

export interface CorsConfig {
  enabled: boolean;
  allowed_origins: string[];
  allowed_methods: string[];
  allowed_headers: string[];
}

export interface IpFilterConfig {
  enabled: boolean;
  black_listed_ips: string[];
  white_listed_ips: string[];
}

export type RateLimitStrategy = 'token_bucket' | 'fixed_window';
export type KeyBy = 'ip' | 'api_key';

export interface TokenBucketConfig {
  capacity: number;
  refill_rate: number;
}

export interface FixedWindowConfig {
  limit: number;
  window_seconds: number;
}

export interface RateLimitConfig {
  enabled: boolean;
  strategy: RateLimitStrategy;
  key_by?: KeyBy;
  token_bucket?: TokenBucketConfig;
  fixed_window?: FixedWindowConfig;
}

export type TokenSource = 'cookie' | 'header';

export interface JwtConfig {
  enabled: boolean;
  algorithm: string;
  public_key: string;
  user_id_claim: string;
  token_source: TokenSource;
  header_name?: string;
  prefix?: string;
  cookie_name?: string;
}

export type LoadBalancingStrategy =
  | 'ROUND_ROBIN'
  | 'RANDOM'
  | 'IP_HASH'
  | 'WEIGHTED_ROUND_ROBIN'
  | 'LEAST_CONNECTIONS';

export interface Backend {
  url: string;
  weight?: number;
}

export interface CreateOrUpdateUpstreamRequestDTO {
  name: string;
  load_balancing_strategy: LoadBalancingStrategy;
  backends: Backend[];
}

export interface GetUpstreamResponseDTO {
  upstream_id: string;
  name: string;
  load_balancing_strategy: LoadBalancingStrategy;
  backends: Backend[];
}

export interface GetAllUpstreamResponseDTO {
  upstreams: GetUpstreamResponseDTO[];
}

export type PathType = 'exact' | 'prefix' | 'regex';
export type AuthMode = 'none' | 'required';

export interface CreateOrUpdateRouteRequestDTO {
  path: string;
  path_type: PathType;
  method: string;
  upstream_name: string;
  auth_mode: AuthMode;
  enabled: boolean;
}

export interface GetRouteResponseDTO {
  route_id: string;
  path: string;
  path_type: PathType;
  method: string;
  upstream_name: string;
  auth_mode: AuthMode;
  enabled: boolean;
}

export interface GetAllRoutesResponseDTO {
  routes: GetRouteResponseDTO[];
}
