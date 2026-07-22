# Mini API Gateway API Guide

This guide covers the control-layer API used to configure projects, upstreams, routes, and project middleware.

## Before You Start

The control layer is served at:

```text
http://localhost:{CONTROL_LAYER_PORT}/api
```

Use the JWT returned through the authentication cookie, or send it as a bearer token:

```text
Cookie: auth-token=<JWT_TOKEN>
```

```text
Authorization: Bearer <JWT_TOKEN>
```

Project, upstream, route, and middleware endpoints require authentication. Replace the example IDs below with MongoDB ObjectIDs returned by the API.

## Configure a Project

The usual configuration flow is:

1. Register and log in.
2. Create a project and save its `project_id` and `api_gw_key`.
3. Create one or more upstreams with backend URLs.
4. Create routes that point to an upstream by name.
5. Add project middleware with the plugin configuration shown below.

### Register

`POST /api/auth/register`

Request body:

```json
{
  "username": "john_doe",
  "password": "secure_password123"
}
```

The request uses this DTO:

```go
type RequestDTO struct {
    Username string `json:"username"`
    Password string `json:"password"`
}
```

The API responds with `204 No Content` and sets the `x-auth-token` cookie.

```bash
curl -i -X POST http://localhost:8080/api/auth/register \
  -H "Content-Type: application/json" \
  -d '{"username":"john_doe","password":"secure_password123"}'
```

### Login

`POST /api/auth/login`

Request body uses the same `RequestDTO` as registration:

```json
{
  "username": "john_doe",
  "password": "secure_password123"
}
```

The API responds with `204 No Content` and sets the `auth-token` cookie.

```bash
curl -i -c cookies.txt -X POST http://localhost:8080/api/auth/login \
  -H "Content-Type: application/json" \
  -d '{"username":"john_doe","password":"secure_password123"}'
```

### Logout

`GET /api/auth/logout`

The API responds with `204 No Content` and removes the `auth-token` cookie.

```bash
curl -i -b cookies.txt -X GET http://localhost:8080/api/auth/logout
```

## Project Endpoints

### Create a Project

`POST /api/projects`

Request body:

```json
{
  "name": "Orders API",
  "access_list": [
    {
      "username": "teammate",
      "permission": "editing"
    }
  ]
}
```

The request uses these DTOs. `permission` must be `editing` or `viewing`.

```go
type CreatProjectRequest struct {
    Name       string   `json:"name"`
    AccessList []Access `json:"access_list"`
}

type Access struct {
    Username   string            `json:"username"`
    Permission models.Permission `json:"permission"`
}

type Permission string

const (
    PermissionEditing Permission = "editing"
    PermissionViewing Permission = "viewing"
)
```

Response body:

```json
{
  "project_id": "507f1f77bcf86cd799439011",
  "api_gw_key": "generated-api-gateway-key"
}
```

```go
type CreatProjectResponse struct {
    Id       bson.ObjectID `json:"project_id"`
    ApiGwKey string        `json:"api_gw_key"`
}
```

```bash
curl -i -b cookies.txt -X POST http://localhost:8080/api/projects \
  -H "Content-Type: application/json" \
  -d '{"name":"Orders API","access_list":[]}'
```

### List Projects

`GET /api/projects`

Response body:

```json
{
  "projects": [
    {
      "project_id": "507f1f77bcf86cd799439011",
      "name": "Orders API"
    }
  ]
}
```

```go
type GetAllProjectResponse struct {
    Projects []ListProjectResponse `json:"projects"`
}

type ListProjectResponse struct {
    ID   bson.ObjectID `json:"project_id"`
    Name string        `json:"name"`
}
```

```bash
curl -i -b cookies.txt http://localhost:8080/api/projects
```

### Get a Project

`GET /api/projects/{projectID}`

Path parameter: `projectID` is the project ObjectID.

Response body:

```json
{
  "project_id": "507f1f77bcf86cd799439011",
  "name": "Orders API",
  "middlewares": [],
  "permission": "editing"
}
```

```go
type GetProjectResponse struct {
    ID          bson.ObjectID       `json:"project_id"`
    Name        string              `json:"name"`
    Middlewares []models.Middleware `json:"middlewares"`
    Permission  string              `json:"permission"`
}
```

```bash
curl -i -b cookies.txt http://localhost:8080/api/projects/507f1f77bcf86cd799439011
```

### Delete a Project

`DELETE /api/projects/{projectID}`

Path parameter: `projectID` is the project ObjectID. The API responds with `204 No Content`.

```bash
curl -i -b cookies.txt -X DELETE \
  http://localhost:8080/api/projects/507f1f77bcf86cd799439011
```

## Middleware Configuration

### Update Project Middlewares

`PATCH /api/projects/{projectID}/middlewares`

The request uses this DTO. Each `name` must be one of the registered plugin names: `cors`, `ip-filter`, `rate-limit`, or `jwt-auth`.

```go
type UpdateMiddlewaresRequest struct {
    Middlewares []models.Middleware `json:"middlewares"`
}

type Middleware struct {
    Name   string                 `bson:"name"`
    Config map[string]interface{} `json:"config"`
}
```

The complete request format is:

```json
{
  "middlewares": [
    {
      "name": "cors",
      "config": {
        "enabled": true,
        "allowed_origins": ["https://app.example.com"],
        "allowed_methods": ["GET", "POST", "PUT", "DELETE"],
        "allowed_headers": ["Content-Type", "Authorization"]
      }
    }
  ]
}
```

Send the complete middleware list in the `middlewares` field when updating it.

### CORS Plugin

Plugin name: `cors`

Config fields:

```go
type CorsConfig struct {
    Enabled        bool     `json:"enabled"`
    AllowedOrigins []string `json:"allowed_origins"`
    AllowedMethods []string `json:"allowed_methods"`
    AllowedHeaders []string `json:"allowed_headers"`
}
```

`allowed_origins` is required. Use `"*"` to allow every origin, or list exact origins. When `enabled` is false, the middleware passes requests through without adding CORS headers.

```json
{
  "name": "cors",
  "config": {
    "enabled": true,
    "allowed_origins": ["https://app.example.com", "https://admin.example.com"],
    "allowed_methods": ["GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"],
    "allowed_headers": ["Content-Type", "Authorization", "X-Request-ID"]
  }
}
```

### IP Filter Plugin

Plugin name: `ip-filter`

Config fields:

```go
type IpFilterConfig struct {
    Enabled        bool     `json:"enabled"`
    BlackListedIps []string `json:"black_listed_ips"`
    WhiteListedIps []string `json:"white_listed_ips"`
}
```

Use valid IPv4 or IPv6 addresses. If `white_listed_ips` is non-empty, only those addresses are allowed. Blacklisted addresses are rejected first.

```json
{
  "name": "ip-filter",
  "config": {
    "enabled": true,
    "black_listed_ips": ["203.0.113.10", "2001:db8::10"],
    "white_listed_ips": ["203.0.113.20", "2001:db8::20"]
  }
}
```

### Rate Limit Plugin

Plugin name: `rate-limit`

The `strategy` enum is `token_bucket` or `fixed_window`. The `key_by` enum is `ip` or `api_key`.

```go
type RateLimitStrategy string

const (
    TokenBucket RateLimitStrategy = "token_bucket"
    FixedWindow RateLimitStrategy = "fixed_window"
)

type KeyBy string

const (
    KeyByIP     KeyBy = "ip"
    KeyByAPIKey KeyBy = "api_key"
)

type RateLimitConfig struct {
    Enabled     bool               `json:"enabled"`
    Strategy    RateLimitStrategy  `json:"strategy"`
    KeyBy       KeyBy              `json:"key_by,omitempty"`
    TokenBucket *TokenBucketConfig `json:"token_bucket,omitempty"`
    FixedWindow *FixedWindowConfig `json:"fixed_window,omitempty"`
}

type TokenBucketConfig struct {
    Capacity   int `json:"capacity"`
    RefillRate int `json:"refill_rate"`
}

type FixedWindowConfig struct {
    Limit         int `json:"limit"`
    WindowSeconds int `json:"window_seconds"`
}
```

Token bucket example:

```json
{
  "name": "rate-limit",
  "config": {
    "enabled": true,
    "strategy": "token_bucket",
    "key_by": "ip",
    "token_bucket": {
      "capacity": 100,
      "refill_rate": 10
    }
  }
}
```

Fixed window example:

```json
{
  "name": "rate-limit",
  "config": {
    "enabled": true,
    "strategy": "fixed_window",
    "key_by": "api_key",
    "fixed_window": {
      "limit": 1000,
      "window_seconds": 60
    }
  }
}
```

The selected strategy must have its matching nested object. `capacity`, `refill_rate`, `limit`, and `window_seconds` must be positive.

### JWT Authentication Plugin

Plugin name: `jwt-auth`

The `token_source` enum is `cookie` or `header`. The supported signing algorithms are `RS256` and `ES256`.

```go
type Source string

const (
    SourceCookie Source = "cookie"
    SourceHeader Source = "header"
)

type JwtConfig struct {
    Enabled     bool   `json:"enabled"`
    Algorithm   string `json:"algorithm"`
    PublicKey   string `json:"public_key"`
    UserIDClaim string `json:"user_id_claim"`
    TokenSource Source `json:"token_source"`
    HeaderName  string `json:"header_name,omitempty"`
    Prefix      string `json:"prefix,omitempty"`
    CookieName  string `json:"cookie_name,omitempty"`
}
```

Header token example:

```json
{
  "name": "jwt-auth",
  "config": {
    "enabled": true,
    "algorithm": "RS256",
    "public_key": "-----BEGIN PUBLIC KEY-----\n...\n-----END PUBLIC KEY-----",
    "user_id_claim": "sub",
    "token_source": "header",
    "header_name": "Authorization",
    "prefix": "Bearer"
  }
}
```

Cookie token example:

```json
{
  "name": "jwt-auth",
  "config": {
    "enabled": true,
    "algorithm": "ES256",
    "public_key": "-----BEGIN PUBLIC KEY-----\n...\n-----END PUBLIC KEY-----",
    "user_id_claim": "user_id",
    "token_source": "cookie",
    "cookie_name": "auth-token"
  }
}
```

`public_key`, `algorithm`, `user_id_claim`, and `token_source` are required. Header mode also requires `header_name` and `prefix`; cookie mode requires `cookie_name`.

### Apply Middleware Configuration

This example configures CORS, IP filtering, and rate limiting together:

```bash
curl -i -b cookies.txt -X PATCH \
  http://localhost:8080/api/projects/507f1f77bcf86cd799439011/middlewares \
  -H "Content-Type: application/json" \
  -d '{
    "middlewares": [
      {
        "name": "cors",
        "config": {
          "enabled": true,
          "allowed_origins": ["https://app.example.com"],
          "allowed_methods": ["GET", "POST", "OPTIONS"],
          "allowed_headers": ["Content-Type", "Authorization"]
        }
      },
      {
        "name": "ip-filter",
        "config": {
          "enabled": true,
          "black_listed_ips": ["203.0.113.10"],
          "white_listed_ips": []
        }
      },
      {
        "name": "rate-limit",
        "config": {
          "enabled": true,
          "strategy": "token_bucket",
          "key_by": "ip",
          "token_bucket": {"capacity": 100, "refill_rate": 10}
        }
      }
    ]
  }'
```

The API responds with `204 No Content`.

### Delete a Middleware

`DELETE /api/projects/{projectID}/middlewares/{name}`

Path parameters: `projectID` is the project ObjectID and `name` is the plugin name, such as `cors` or `rate-limit`.

```bash
curl -i -b cookies.txt -X DELETE \
  http://localhost:8080/api/projects/507f1f77bcf86cd799439011/middlewares/cors
```

## Access List

### Update Project Access

`PATCH /api/projects/{projectID}/accesslist`

Request body uses `Permission`, whose values are `editing` and `viewing`:

```json
{
  "access_list": [
    {"username": "editor", "permission": "editing"},
    {"username": "viewer", "permission": "viewing"}
  ]
}
```

```go
type UpdateAccessListRequest struct {
    AccessList []Access `json:"access_list"`
}
```

The API responds with `204 No Content`.

```bash
curl -i -b cookies.txt -X PATCH \
  http://localhost:8080/api/projects/507f1f77bcf86cd799439011/accesslist \
  -H "Content-Type: application/json" \
  -d '{"access_list":[{"username":"editor","permission":"editing"}]}'
```

## Upstream Endpoints

### Create an Upstream

`POST /api/projects/{projectID}/upstreams`

Path parameter: `projectID` is the project ObjectID.

Request body:

```json
{
  "name": "orders-backend",
  "load_balancing_strategy": "ROUND_ROBIN",
  "backends": [
    {"url": "http://orders-1:8080"},
    {"url": "http://orders-2:8080"}
  ]
}
```

The request uses this DTO. `load_balancing_strategy` must be `ROUND_ROBIN`, `RANDOM`, `IP_HASH`, `WEIGHTED_ROUND_ROBIN`, or `LEAST_CONNECTIONS`.

```go
type CreateOrUpdateUpstreamRequestDTO struct {
    Name                  string                       `json:"name"`
    LoadBalancingStrategy models.LoadBalancingStrategy `json:"load_balancing_strategy"`
    Backends              []models.Backend             `json:"backends"`
}

type Backend struct {
    URL    string `json:"url"`
    Weight *int   `json:"weight,omitempty"`
}

type LoadBalancingStrategy string

const (
    RoundRobinLoadBalancing LoadBalancingStrategy = "ROUND_ROBIN"
    RandomLoadBalancing     LoadBalancingStrategy = "RANDOM"
    IPHashLoadBalancing     LoadBalancingStrategy = "IP_HASH"
    WeightedRoundRobin      LoadBalancingStrategy = "WEIGHTED_ROUND_ROBIN"
    LeastConnections        LoadBalancingStrategy = "LEAST_CONNECTIONS"
)
```

Use `weight` when using `WEIGHTED_ROUND_ROBIN`.

Response body:

```json
{"upstream_id":"507f1f77bcf86cd799439012"}
```

```go
type CreateUpstreamResponseDTO struct {
    UpstreamID bson.ObjectID `json:"upstream_id"`
}
```

```bash
curl -i -b cookies.txt -X POST \
  http://localhost:8080/api/projects/507f1f77bcf86cd799439011/upstreams \
  -H "Content-Type: application/json" \
  -d '{"name":"orders-backend","load_balancing_strategy":"ROUND_ROBIN","backends":[{"url":"http://orders-1:8080"},{"url":"http://orders-2:8080"}]}'
```

### List Upstreams

`GET /api/projects/{projectID}/upstreams`

Response body:

```json
{
  "upstreams": [
    {
      "upstream_id": "507f1f77bcf86cd799439012",
      "name": "orders-backend",
      "load_balancing_strategy": "ROUND_ROBIN",
      "backends": [{"url":"http://orders-1:8080"}]
    }
  ]
}
```

```go
type GetAllUpstreamResponseDTO struct {
    Upstreams []GetUpstreamResponseDTO `json:"upstreams"`
}

type GetUpstreamResponseDTO struct {
    ID                    bson.ObjectID                `json:"upstream_id"`
    Name                  string                       `json:"name"`
    LoadBalancingStrategy models.LoadBalancingStrategy `json:"load_balancing_strategy"`
    Backends              []models.Backend             `json:"backends"`
}
```

```bash
curl -i -b cookies.txt \
  http://localhost:8080/api/projects/507f1f77bcf86cd799439011/upstreams
```

### Get an Upstream

`GET /api/projects/{projectID}/upstreams/{upstreamID}`

Path parameters: `projectID` and `upstreamID` are ObjectIDs. The response uses `GetUpstreamResponseDTO` shown above.

```bash
curl -i -b cookies.txt \
  http://localhost:8080/api/projects/507f1f77bcf86cd799439011/upstreams/507f1f77bcf86cd799439012
```

### Update an Upstream

`PUT /api/projects/{projectID}/upstreams/{upstreamID}`

Path parameters: `projectID` and `upstreamID` are ObjectIDs. The request uses `CreateOrUpdateUpstreamRequestDTO` shown in the create endpoint.

```bash
curl -i -b cookies.txt -X PUT \
  http://localhost:8080/api/projects/507f1f77bcf86cd799439011/upstreams/507f1f77bcf86cd799439012 \
  -H "Content-Type: application/json" \
  -d '{"name":"orders-backend","load_balancing_strategy":"LEAST_CONNECTIONS","backends":[{"url":"http://orders-1:8080"},{"url":"http://orders-2:8080"}]}'
```

The API responds with `204 No Content`.

### Delete an Upstream

`DELETE /api/projects/{projectID}/upstreams/{upstreamID}`

Path parameters: `projectID` and `upstreamID` are ObjectIDs. The API responds with `204 No Content`.

```bash
curl -i -b cookies.txt -X DELETE \
  http://localhost:8080/api/projects/507f1f77bcf86cd799439011/upstreams/507f1f77bcf86cd799439012
```

## Route Endpoints

### Create a Route

`POST /api/projects/{projectID}/routes`

Path parameter: `projectID` is the project ObjectID.

Request body:

```json
{
  "path": "/orders",
  "path_type": "prefix",
  "method": "GET",
  "upstream_name": "orders-backend",
  "auth_mode": "required",
  "enabled": true
}
```

The request uses this DTO. `path_type` must be `exact`, `prefix`, or `regex`; `auth_mode` must be `none` or `required`.

```go
type CreateOrUpdateRouteRequestDTO struct {
    Path         string          `json:"path"`
    PathType     models.PathType `json:"path_type"`
    Method       string          `json:"method"`
    UpstreamName string          `json:"upstream_name"`
    AuthMode     models.AuthMode `json:"auth_mode"`
    Enabled      bool            `json:"enabled"`
}

type PathType string

const (
    PathExact  PathType = "exact"
    PathPrefix PathType = "prefix"
    PathRegex  PathType = "regex"
)

type AuthMode string

const (
    AuthNone     AuthMode = "none"
    AuthRequired AuthMode = "required"
)
```

Response body:

```json
{"route_id":"507f1f77bcf86cd799439013"}
```

```go
type CreateRouteResponseDTO struct {
    RouteID bson.ObjectID `json:"route_id"`
}
```

```bash
curl -i -b cookies.txt -X POST \
  http://localhost:8080/api/projects/507f1f77bcf86cd799439011/routes \
  -H "Content-Type: application/json" \
  -d '{"path":"/orders","path_type":"prefix","method":"GET","upstream_name":"orders-backend","auth_mode":"required","enabled":true}'
```

### List Routes

`GET /api/projects/{projectID}/routes`

Response body:

```json
{
  "routes": [
    {
      "route_id": "507f1f77bcf86cd799439013",
      "path": "/orders",
      "path_type": "prefix",
      "method": "GET",
      "upstream_name": "orders-backend",
      "auth_mode": "required",
      "enabled": true
    }
  ]
}
```

```go
type GetAllRoutesResponseDTO struct {
    Routes []GetRouteResponseDTO `json:"routes"`
}

type GetRouteResponseDTO struct {
    ID           bson.ObjectID   `json:"route_id"`
    Path         string          `json:"path"`
    PathType     models.PathType `json:"path_type"`
    Method       string          `json:"method"`
    UpstreamName string          `json:"upstream_name"`
    AuthMode     models.AuthMode `json:"auth_mode"`
    Enabled      bool            `json:"enabled"`
}
```

```bash
curl -i -b cookies.txt \
  http://localhost:8080/api/projects/507f1f77bcf86cd799439011/routes
```

### Get a Route

`GET /api/projects/{projectID}/routes/{routeID}`

Path parameters: `projectID` and `routeID` are ObjectIDs. The response uses `GetRouteResponseDTO` shown in the list endpoint.

```bash
curl -i -b cookies.txt \
  http://localhost:8080/api/projects/507f1f77bcf86cd799439011/routes/507f1f77bcf86cd799439013
```

### Update a Route

`PUT /api/projects/{projectID}/routes/{routeID}`

Path parameters: `projectID` and `routeID` are ObjectIDs. The request uses `CreateOrUpdateRouteRequestDTO` shown in the create endpoint.

```bash
curl -i -b cookies.txt -X PUT \
  http://localhost:8080/api/projects/507f1f77bcf86cd799439011/routes/507f1f77bcf86cd799439013 \
  -H "Content-Type: application/json" \
  -d '{"path":"/orders/{id}","path_type":"prefix","method":"GET","upstream_name":"orders-backend","auth_mode":"required","enabled":true}'
```

The API responds with `204 No Content`.

### Delete a Route

`DELETE /api/projects/{projectID}/routes/{routeID}`

Path parameters: `projectID` and `routeID` are ObjectIDs. The API responds with `204 No Content`.

```bash
curl -i -b cookies.txt -X DELETE \
  http://localhost:8080/api/projects/507f1f77bcf86cd799439011/routes/507f1f77bcf86cd799439013
```

## Gateway Request

After the project, upstream, route, and middleware configuration is synchronized to the gateway plane, send application traffic through the gateway using the project gateway key as `X-Gateway-Key`.

```bash
curl -i http://localhost:{GATEWAY_PORT}/orders \
  -H "X-Gateway-Key: <API_GW_KEY>"
```
