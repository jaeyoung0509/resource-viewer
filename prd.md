PRD: K8s-Based Distributed Resource Monitoring (Web-htop)

1. Product Overview
1.1 Purpose
- Collect real-time CPU, memory, and disk usage from all nodes in a Kubernetes cluster.
- Visualize metrics in a web dashboard served by the Hub.
- Learning goals: Linux /proc, container mounts, WebSocket, Redis Pub/Sub, K8s DaemonSet, service discovery.

1.2 Scope
- In scope: Node Agent (Go), Redis, Hub Server (Go), React dashboard (TypeScript).
- Out of scope: alerting, long-term storage (TSDB), auth/ACL.

1.3 Design Principles
- Simplicity first: fast local setup and minimal moving parts.
- Maintainability: clear package boundaries and config via env.
- Productivity: Kustomize overlays for local/prod.

2. User Scenarios
2.1 Primary user
- Developer learning distributed monitoring on a local K8s cluster (Orbstack).

2.2 Key scenarios
- User deploys to a local multi-node cluster.
- Agent on each node publishes metrics to Redis.
- Hub subscribes and broadcasts to browsers via WebSocket.
- User views node metrics and scales selected deployments from the UI.

3. System Architecture
3.1 Components
- Node Agent (DaemonSet): reads host /proc and publishes metrics to Redis.
- Redis (Deployment/StatefulSet): Pub/Sub broker.
- Hub Server (Deployment): subscribes to Redis and streams WebSocket updates.
- Frontend (React + TS): static build served by Hub.

3.2 Data Flow
1) Agent collects resource snapshot and serializes to JSON.
2) Redis PUBLISH on `metrics:<node-name>`.
3) Hub SUBSCRIBE on `metrics:*`.
4) Hub broadcasts to all WebSocket clients.

4. Functional Requirements
4.1 Node Agent (Go)
- Deployment: DaemonSet.
- Graceful shutdown on SIGINT/SIGTERM.
- Data collection:
  - Use gopsutil.
  - HOST_PROC env to select /proc path (default /host/proc).
- Data publish:
  - JSON payload.
  - Redis channel `metrics:<node-name>`.
- Node identity: os.Hostname().

4.2 Redis
- Deployment: simple Deployment or StatefulSet.
- No password for local learning.
- Service DNS: redis-service:6379.

4.3 Hub Server (Go)
- Deployment: Deployment (replicas = 1).
- Redis subscribe: `metrics:*` pattern, retry on reconnect.
- WebSocket:
  - gorilla/websocket.
  - Endpoint `/ws`.
  - Multi-client broadcast.
- TLS:
  - HTTPS/WSS endpoint with local TLS.
- Scale control API:
  - GET `/api/targets` returns allowed deployment targets and current replicas.
  - POST `/api/scale` scales a target (0..SCALE_MAX).
  - RBAC scoped to allowed deployments in namespace.

4.4 Frontend (React + TypeScript)
- Served as static build by Hub.
- Features:
  - Connects to `wss://<host>/ws`.
  - Displays node metrics by node name.
  - Scaling controls for selected deployments.

5. Non-Functional Requirements
5.1 Performance
- Metrics reflected on dashboard within 1s.
- Works on 3-node local cluster (1 master + 2 workers).

5.2 Reliability
- Retry on Redis connection failures.
- WebSocket client connections handled safely.

5.3 Security
- Browser <-> Hub uses TLS (mkcert preferred).
- In-cluster traffic trusted for local dev.

5.4 Operability
- Kustomize overlays for local/prod.
- Env-based configuration.

6. Data Schema
6.1 Metric payload (JSON)
{
  "node": "worker-1",
  "cpu": 23.4,
  "memory": 61.2,
  "disk": 45.8,
  "timestamp": 1710000000
}
- node: node name
- cpu/memory/disk: percent usage
- timestamp: Unix epoch (seconds)

7. Local Development and Deployment
7.1 Tools
- Orbstack (Kubernetes enabled)
- Docker
- Go 1.25.4
- kubectl
- kustomize
- wscat

7.2 Local steps
- Generate TLS: ./scripts/gen-tls-secret.sh
- Build images: make docker-agent / make docker-hub
- Deploy: make k8s-apply
- Open: https://localhost:30443

8. Risks and Mitigations
- /proc path error: ensure HOST_PROC is set and hostPath mounted.
- Redis DNS error: verify redis-service and namespace.
- TLS errors: check cert Secret and local trust.
- Scaling API safety: restrict to allowlist and namespace-scoped RBAC.

9. Packaging Structure
9.1 Go layout
- cmd/agent
- cmd/hub
- internal/agent
- internal/hub
- internal/config
- internal/metrics
- pkg/health

9.2 K8s manifests
- deploy/base
- deploy/overlays/local
- deploy/overlays/prod

10. Future Extensions (Out of Scope)
- TSDB storage
- Alerting/Slack integration
- AuthZ/AuthN
