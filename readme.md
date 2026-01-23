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
- Frontend: Static HTML/JS dashboard served by Hub.

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
- Connects to wss://<host>/ws.
- Renders node metrics in real time.

5. Quick Start (Local K8s)
1) Generate TLS for Hub:
   ./scripts/gen-tls-secret.sh
2) Build images:
   make docker-agent
   make docker-hub
3) Deploy:
   make k8s-apply
4) Open dashboard:
   https://localhost:30443

6. Phases
Phase 1: Docker build for agent/hub.
Phase 2: K8s deploy on Orbstack with Kustomize.

7. Troubleshooting
- Agent reads wrong /proc: ensure HOST_PROC=/host/proc and hostPath mount.
- Hub can't reach Redis: verify redis-service:6379 and namespace.
- TLS error: check Secret mount and local cert generation.
