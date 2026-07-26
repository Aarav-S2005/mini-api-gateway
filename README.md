# Mini API Gateway

A lightweight, high-performance API gateway written in Go that separates configuration management (control plane) from request processing (gateway plane). The control plane provides authenticated APIs for managing projects, routes, upstreams, and middleware. The gateway plane consumes this configuration to forward client traffic to healthy backend services using load balancing and health checks.

---

## Architecture Overview

Mini API Gateway is built on a clean control-plane/data-plane separation architecture, enabling independent scaling and management of configuration versus traffic handling.

### Core Components

- **Control Plane**: Manages gateway configuration with authenticated REST APIs for projects, routes, upstreams, and middleware
- **Gateway Plane**: High-performance reverse proxy that loads configuration into an in-memory snapshot and routes client requests to backends
- **Plugin Manager**: Provides composable middleware plugins (CORS, IP filtering, rate limiting, JWT authentication)
- **Database**: Responsible for creating a connection to mongoDB

### Design Principles

The gateway prioritizes request-path performance by building a compiled in-memory snapshot at startup and on configuration changes. This ensures that request handling is independent from database queries during normal traffic. Configuration updates flow through Redis, allowing multiple gateway instances to stay synchronized in real-time.

---

## Control Plane

The control plane is the configuration hub of the gateway system, implemented under [app/control-layer](app/control-layer). It provides a REST API for managing all gateway configurations and uses MongoDB for persistence.

### Core Responsibilities

**Authentication & User Management**
- User registration and login with JWT-based authentication
- Secure password handling and session management via HTTP cookies

**Project Management**
- Create and manage logical groupings (projects) that contain routes and upstreams
- Assign a unique `api_gw_key` to each project for gateway routing
- Associate middleware policies at the project level

**Upstream Configuration**
- Define backend services (upstreams) with multiple backend URLs
- Configure load-balancing strategies per upstream
- Set weights for backends in weighted-round-robin scenarios

**Route Configuration**
- Create HTTP routes by path patterns and methods
- Link routes to upstreams for backend selection
- Enable/disable routes without redeployment

**Middleware Management**
- Attach plugins to projects with custom configuration
- Specify JWT authentication requirements per route
- Order middleware execution (CORS → IP Filter → Rate Limit)

### Configuration Publishing

When configuration changes are saved to MongoDB, the control plane publishes Redis notifications. This triggers live synchronization in all connected gateway instances, which rebuild only the affected snapshot entries. 

See [DOCS.md](DOCS.md) for complete API documentation.

---

## Gateway Plane

The gateway plane is the request-handling core of the system, implemented under [app/gateway](app/gateway). It loads control plane configuration into a high-performance in-memory snapshot and forwards client requests to backend services.

### Load Balancing

The gateway supports multiple load-balancing strategies for intelligent traffic distribution across backend servers:

- **Round Robin** — Distributes requests equally across healthy backends in a circular pattern
- **Weighted Round Robin** — Distributes requests proportionally based on configured backend weights
- **Random** — Selects a random healthy backend for each request
- **IP Hash** — Routes based on client IP, ensuring requests from the same client reach the same backend
- **Least Connections** — Routes to the backend with the fewest active connections

Each upstream maintains atomic counters for active connections and backend health status. Strategies only consider healthy backends; unhealthy backends are temporarily excluded.

### Health Checking

A continuous health checker probes all configured backends at regular intervals:

- **Check Frequency** — Runs every 10 seconds
- **Check Methods** — Attempts HTTP health endpoint first (`/health`), falls back to TCP connection
- **State Transitions** — Requires 3 consecutive successes to mark backend healthy; 3 consecutive failures to mark unhealthy
- **Failure Detection** — 5xx backend responses automatically trigger failure counting
- **Timeout** — Health checks have a 3-second timeout

The health checker maintains separate success/failure counters per backend to prevent flapping on transient failures.

### Live Syncing

Configuration changes propagate through Redis pub/sub to all running gateway instances:

- **Change Detection** — Control plane publishes events when projects, routes, or upstreams are modified
- **Selective Snapshot Rebuild** — Gateway instances rebuild only affected snapshot sections (e.g., single project or upstream)
- **Zero Downtime** — New snapshot is atomically swapped; existing requests continue on old snapshot
- **In-Memory Registry** — Provides consistent snapshot view to all concurrent request handlers

---

## Request Handling in Gateway

When a request arrives at the gateway, the following process occurs:

### Request Path

1. **Gateway Key Validation** — Client must include `X-Gateway-Key` header; requests without it receive `401 Unauthorized`

2. **Project Resolution** — Gateway looks up the project in the in-memory snapshot using the API key; invalid keys return `401 Unauthorized`

3. **Middleware Chain Execution** — Request passes through project-level middleware in order:
   - CORS validation (if configured)
   - IP filtering (if configured)
   - Rate limiting (if configured)

4. **Route Matching** — Gateway matches request path and HTTP method against enabled routes in the project's Chi router

5. **Upstream Selection** — Matched route specifies which upstream to use for backend discovery

6. **Backend Selection** — Load-balancing strategy selects a healthy backend from the upstream's backend list

7. **Request Forwarding** — Reverse proxy forwards the HTTP request to the selected backend, including request body, headers, and query parameters

8. **Response Handling** — Backend response is returned to client; 5xx responses automatically trigger failure counting in health tracking

9. **Logging** — Request is logged with method, path, project ID, and total processing time

---

## Directory Structure

```text
app/
├── control-layer/
│   ├── cmd/server/
│   │   └── main.go                    Control plane entry point
│   ├── config/
│   │   ├── config.go                  Environment configuration
│   │   └── redis.go                   Redis event publisher
│   ├── internal/
│   │   ├── app_error/                 Error handling and HTTP responses
│   │   └── endpoints/
│   │       ├── auth/                  User registration, login, logout
│   │       ├── project/               Project management endpoints
│   │       ├── route/                 Route configuration endpoints
│   │       └── upstream/              Upstream backend management
│   └── Dockerfile                     Control plane container image
├── gateway/
│   ├── cmd/server/
│   │   └── main.go                    Gateway entry point
│   ├── config/
│   │   ├── config.go                  Environment configuration
│   │   └── redis.go                   Redis subscription setup
│   ├── internal/
│   │   ├── endpoint/
│   │   │   └── serve.go               HTTP request handler
│   │   ├── health/
│   │   │   └── checker.go             Backend health checking
│   │   ├── lb/
│   │   │   ├── manager.go             Load balancer state management
│   │   │   ├── strategy.go            Strategy selection dispatcher
│   │   │   ├── round_robin.go         Round robin implementation
│   │   │   ├── weighted_round_robin.go Weighted round robin implementation
│   │   │   ├── random.go              Random selection implementation
│   │   │   ├── ip_hash.go             IP hash implementation
│   │   │   └── least_connection.go    Least connections implementation
│   │   ├── proxy/
│   │   │   └── reverse-proxy.go       Reverse proxy builder and error handling
│   │   ├── store/
│   │   │   ├── snapshot.go            In-memory snapshot data structures
│   │   │   ├── loader.go              Initial snapshot building from MongoDB
│   │   │   └── registry.go            Atomic snapshot storage and retrieval
│   │   └── sync/
│   │       └── sync-loader.go         Live configuration synchronization
│   └── Dockerfile                     Gateway container image
├── db/
│   ├── mongo.go                       MongoDB connection
│   └── models/
│       ├── user.go                    User data model
│       ├── project.go                 Project data model
│       ├── upstream.go                Upstream data model
│       ├── route.go                   Route data model
│       └── middleware.go              Middleware configuration model
├── plugin-manager/
│   ├── plugin/
│   │   ├── plugin.go                  Plugin interface definition
│   │   ├── cors-plugin.go             CORS middleware plugin
│   │   ├── ip-filter-plugin.go        IP filtering middleware plugin
│   │   ├── rate-limit-plugin.go       Rate limiting middleware plugin
│   │   ├── jwt-plugin.go              JWT authentication middleware plugin
│   │   ├── fixed-window-impl.go       Fixed window rate limiter implementation
│   │   └── token-bucket-impl.go       Token bucket rate limiter implementation
│   └── registry/
│       └── registry.go                Plugin registry and loader
└── go.mod                             Go module definition and dependencies
```

---

## Reverse Proxy Implementation

### Reverse Proxy Architecture

The gateway builds reverse proxies using Go's standard `net/http/httputil.ReverseProxy` with a shared HTTP transport. This design provides several benefits:

- **Connection Pooling** — Single transport maintains connection pool shared across all reverse proxies, reducing overhead
- **Keep-Alive Support** — HTTP keep-alive is enabled by default, reducing latency for repeated backend calls
- **Concurrency Safe** — Transport is thread-safe and handles concurrent requests from multiple goroutines

Each reverse proxy is configured with:
- **URL Rewriting** — Removes the `X-Gateway-Key` header before forwarding to backend
- **Error Handler** — Catches connection failures and marks backend unhealthy
- **Response Handler** — Detects 5xx responses and triggers health state updates
- **Flush Interval** — Set to -1 for immediate response streaming


### Connection Tracking

Active connections are tracked atomically per backend using the load balancer manager:

- **On Request Start** — `IncConn()` increments connection counter
- **On Request End** — `DecConn()` decrements connection counter
- **Used By Least Connections** — Strategy queries current count to select least-loaded backend

---

## Plugin Manager

The plugin system provides extensible, composable middleware for request processing. Plugins are registered in a central registry and instantiated with configuration from the control plane.

### Built-in Plugins

**CORS Plugin**
- Validates Origin header against allowed domains
- Responds to preflight OPTIONS requests
- Adds appropriate CORS response headers
- Configuration: list of allowed origins, allowed methods, allowed headers

**IP Filter Plugin**
- Whitelists or blacklists client IP addresses
- Supports CIDR notation for IP ranges
- Checks X-Forwarded-For header for proxy scenarios
- Configuration: whitelist/blacklist mode, list of IP addresses or ranges

**Rate Limit Plugin**
- Implements two rate-limiting strategies: fixed window and token bucket
- Tracks rates per client IP address
- Returns `429 Too Many Requests` when limit exceeded
- Configuration: strategy type, requests per interval, interval duration

**JWT Authentication Plugin**
- Validates JWT bearer tokens in Authorization header
- Extracts JWT secret from control plane configuration
- Routes can require authentication for specific endpoints
- Configuration: JWT secret, token validation algorithm

### Plugin Interface

All plugins implement the `Plugin` interface:

```go
type Plugin interface {
    Name() string
    Validate(config map[string]interface{}) error
    CreateMiddleware(config map[string]interface{}) (func(next http.Handler) http.Handler, error)
}
```

**Validation** — Plugins validate configuration before creation to catch errors early and avoid runtime failures.

**Middleware Creation** — Plugins create HTTP middleware functions that wrap the next handler in the chain.

### Middleware Chain Execution

Project middleware is executed in order:

1. CORS (if configured)
2. IP Filter (if configured)  
3. Rate Limit (if configured)
4. Route Matching and JWT Auth (if route requires auth)

If any middleware rejects the request, subsequent middleware is skipped.

---

## Getting Started

### Requirements

- **Go** — Version 1.26.1 or compatible
- **MongoDB** — For storing projects, routes, upstreams, and middleware configuration
- **Redis** — For pub/sub notifications between control plane and gateway instances

### Option 1: Using Docker Compose

Docker Compose is the quickest way to get the full system running locally:

**Prerequisites**
- Docker and Docker Compose installed

**Steps**

1. Clone the repository:
   ```bash
   git clone https://github.com/Aarav-S2005/mini-api-gateway.git
   cd mini-api-gateway
   ```

2. Start all services:
   ```bash
   docker-compose up --build
   ```

   This starts:
   - Redis on port 6379
   - MongoDB on port 27017
   - Control plane on port 3001
   - Gateway on port 3002


### Option 2: Manual Setup

For development or debugging, you can run services manually:

**Prerequisites**
- MongoDB running locally or accessible via network
- Redis running locally or accessible via network
- Go 1.26.1+ installed

**Steps**

1. Clone the repository:
   ```bash
   git clone https://github.com/Aarav-S2005/mini-api-gateway.git
   cd mini-api-gateway/app
   ```

2. Configure environment variables:
   ```bash
   cp .env.example .env
   ```

3. Edit `.env` with your configuration:
   ```bash
   MONGO_URI=mongodb://localhost:27017/
   MONGO_DB=api-gw
   REDIS_URI=localhost:6379
   CONTROL_PLANE_PORT=3001
   GATEWAY_PORT=3002
   JWT_SECRET_CONTROL_PLANE=your-secret-key-here
   ```

4. Start MongoDB (if running locally):
   ```bash
   mongod --dbpath ./data
   ```

5. Start Redis (if running locally):
   ```bash
   redis-server
   ```

6. Build and run the control plane:
   ```bash
   go run ./control-layer/cmd/server/main.go
   ```

   Control plane listens on port specified in `CONTROL_PLANE_PORT`

7. In another terminal, build and run the gateway:
   ```bash
   go run ./gateway/cmd/server/main.go
   ```

   Gateway listens on port specified in `GATEWAY_PORT`

### First Steps After Startup

Once both services are running, follow the configuration flow in [DOCS.md](DOCS.md):

1. Register a user via the authentication API
2. Log in and obtain a JWT token
3. Create a project and save the `api_gw_key`
4. Create upstreams with backend URLs
5. Create routes pointing to upstreams
6. Add middleware configuration to the project
7. Send test requests through the gateway with `X-Gateway-Key` header

For detailed API documentation and configuration examples, see [DOCS.md](DOCS.md).
