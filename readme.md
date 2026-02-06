PRD: K8s-Based Distributed Resource Monitoring (Web-htop)

1. Project Overview
Project Name: K8s-Web-Htop (working title)
Goal: Collect real-time CPU, memory, and disk usage from all nodes in a Kubernetes cluster and visualize them in a web dashboard.
Learning Goals:
- OS: Linux /proc filesystem, container isolation, host path mounts.
- Network: WebSocket, Redis Pub/Sub, K8s service discovery.
- Infra: DaemonSet pattern, multi-node architecture.

2. Architecture
- Node Agent (DaemonSet): Go app per node. Reads host /proc and publishes to Redis.
- Redis (Broker): Pub/Sub relay.
- Hub Server (Deployment): Subscribes to Redis and pushes to WebSocket clients.
- Frontend: React (TypeScript) dashboard served by Hub.

3. Local Environment (Orbstack)
Required:
- Orbstack (Kubernetes enabled)
- Docker
- Go 1.25.4
- kubectl
- kustomize
- wscat

4. Core Requirements
Agent:
- Runs as DaemonSet.
- Reads /proc via HOST_PROC env (default /host/proc).
- Publishes JSON to Redis channel metrics:<node-name>.

Redis:
- Simple Deployment + Service (no auth for learning).

Hub:
- Subscribes to metrics:*.
- WebSocket endpoint /ws.
- HTTPS/WSS with self-signed TLS (local).

Frontend:
- React + TypeScript UI connects to wss://<host>/ws.
- Shows node metrics and allows scaling selected deployments.

5. Quick Start (Local K8s)
1) Generate TLS for Hub (mkcert preferred):
   ./scripts/gen-tls-secret.sh
2) Create local InfluxDB secret file:
   cp deploy/overlays/local/influxdb-secret.env.example deploy/overlays/local/influxdb-secret.env
3) Build images:
   make docker-agent
   make docker-hub
4) Deploy:
   make k8s-apply
5) Open dashboard:
   https://localhost:30443

Scaling targets (default):
- resource-hub
- redis

Notes:
- Hub image build runs Vite to bundle the React UI.
- Pod/Deployment metrics require metrics-server in the cluster.

OTel + InfluxDB (optional):
1) Build trace service:
   make docker-trace
2) Deploy base resources (includes InfluxDB + OTel Collector):
   make k8s-apply
3) Generate a trace:
   kubectl port-forward svc/trace-svc 8081:8081
   curl http://localhost:8081/api/trace

6. Phases
Phase 1: Docker build for agent/hub.
Phase 2: K8s deploy on Orbstack with Kustomize.

7. Troubleshooting
- Agent reads wrong /proc: ensure HOST_PROC=/host/proc and hostPath mount.
- Hub can't reach Redis: verify redis-service:6379 and namespace.
- TLS error: check Secret mount and local cert generation.
