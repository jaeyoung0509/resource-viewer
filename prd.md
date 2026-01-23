PRD: K8s 기반 분산 리소스 모니터링 시스템 (Web-htop)

1. 제품 개요
1.1 목적
- Kubernetes 클러스터 내 모든 노드의 리소스(CPU, Memory, Disk)를 실시간으로 수집하고 웹 대시보드에서 시각화한다.
- 학습 목표: Linux /proc 이해, 컨테이너 네임스페이스/볼륨 마운트, WebSocket, Redis Pub/Sub, K8s DaemonSet 및 Service Discovery.

1.2 범위
- 포함: Node Agent(Go), Redis, Hub Server(Go), Web Dashboard(HTML/JS).
- 제외: 알림/이벤트 룰 엔진, 장기 저장(TSDB), 사용자 인증/권한 관리.

1.3 설계 원칙
- 단순성 우선: 로컬에서 빠르게 구동 가능.
- 학습/유지보수 용이: 표준 라이브러리 중심, 익숙한 패턴.
- 생산성: K8s 매니페스트는 Kustomize로 관리.

2. 사용자 시나리오
2.1 1차 사용자
- 로컬 K3d 환경에서 분산 모니터링 시스템을 학습/구현하려는 개발자.

2.2 주요 시나리오
- 사용자는 K3d 클러스터(1 master + 2 worker)를 띄운다.
- 각 노드에 배포된 Agent가 리소스 데이터를 Redis로 전송한다.
- Hub Server가 Redis를 구독하고 브라우저에 실시간으로 전달한다.
- 사용자는 대시보드에서 노드별 CPU/Memory/Disk 사용량을 확인한다.

3. 시스템 아키텍처
3.1 구성요소
- Node Agent (DaemonSet): 각 노드에서 /proc 데이터를 읽어 Redis로 전송.
- Redis (StatefulSet/Deployment): Pub/Sub 브로커.
- Hub Server (Deployment): Redis 구독 후 WebSocket으로 브라우저에 브로드캐스트.
- Frontend (Static): Hub Server에서 제공되는 정적 HTML/JS 대시보드.

3.2 데이터 흐름
1) Agent가 노드의 리소스 스냅샷을 JSON으로 직렬화.
2) Redis 채널 `metrics:<node-name>`로 PUBLISH.
3) Hub Server가 `metrics:*` 패턴으로 SUBSCRIBE.
4) Hub가 받은 메시지를 모든 WebSocket 클라이언트에 WriteJSON.

4. 기능 요구사항
4.1 Node Agent (Go)
- 배포 형태: Kubernetes DaemonSet.
- OS 시그널 처리: SIGINT/SIGTERM 수신 시 Graceful 종료.
- 데이터 수집:
  - gopsutil 사용.
  - 환경변수 `HOST_PROC`를 통해 /proc 경로 설정.
  - 로컬 개발: /proc, K8s 배포: /host/proc.
- 데이터 전송:
  - JSON 직렬화.
  - Redis Pub/Sub 채널: `metrics:<node-name>`.
- 노드 식별: `os.Hostname()`으로 node-name 확보.

4.2 Redis
- 배포 형태: StatefulSet 또는 Deployment.
- 설정: 기본 설정(비밀번호 없음).
- 서비스: `redis-service:6379`로 접근 가능.

4.3 Hub Server (Go)
- 배포 형태: Kubernetes Deployment (Replicas: 1).
- Redis 연결:
  - `metrics:*` 패턴 SUBSCRIBE.
  - 연결 실패 시 재시도(지수 백오프 또는 고정 간격).
- WebSocket:
  - `gorilla/websocket` 사용.
  - 엔드포인트: `/ws`.
  - 다중 클라이언트 동시 연결 지원.
- TLS:
  - Hub Server는 HTTPS/WSS로 직접 통신(단일 종료 지점).
  - 로컬용 self-signed 인증서 사용.
- 브로드캐스팅:
  - Redis 메시지를 수신 즉시 모든 WebSocket 클라이언트에 WriteJSON.

4.4 Frontend (Dashboard)
- 형태: Hub Server가 제공하는 정적 HTML(`index.html`).
- 기능:
  - `ws://localhost:8080/ws` 연결.
  - 노드 이름별로 데이터 분기하여 화면 표시.
  - (선택) Chart.js로 실시간 라인 그래프 표시.

5. 비기능 요구사항
5.1 성능
- 실시간성: 데이터 수집 및 대시보드 반영이 1초 내 지연.
- 확장성: 최소 3노드(1 master + 2 worker) 환경에서 정상 동작.

5.2 안정성
- Agent는 Redis 연결 문제 시 재시도 로직 필요.
- Hub Server는 Redis 연결 장애 시 복구 가능해야 함.
- WebSocket 클라이언트 연결/해제 안전 처리.

5.3 보안
- 브라우저 <-> Hub 구간은 TLS 적용.
- Redis 및 내부 통신은 클러스터 내부 트러스트 가정으로 평문.
- 실제 운영 환경에서는 mTLS 또는 네트워크 정책 고려.

5.4 운영성
- 구성은 Kustomize로 분리(dev/local/prod).
- 설정은 환경변수 기반으로 단순화.

6. 데이터 스키마
6.1 메시지 예시 (JSON)
{
  "node": "worker-1",
  "cpu": 23.4,
  "memory": 61.2,
  "disk": 45.8,
  "timestamp": 1710000000
}
- node: 노드명
- cpu/memory/disk: 사용률(%) 또는 절대값(사전 정의 필요)
- timestamp: Unix epoch (초)

7. 로컬 개발 및 배포
7.1 필수 도구
- Orbstack (Kubernetes 활성화)
- Docker (이미지 빌드용)
- Go 1.20+
- kubectl
- kustomize
- wscat

7.2 Orbstack Kubernetes 사용
- Orbstack에서 Kubernetes 활성화 후 kubectl context 확인.
- NodePort 또는 Ingress 중 하나로 Hub Server를 노출.

7.3 로컬 개발 단계
Phase 1 (Dockerizing)
- 멀티 스테이지 Dockerfile 작성.
- `my-agent:v1`, `my-hub:v1` 이미지 빌드.

Phase 2 (K8s 배포)
- `kubectl apply -k deploy/overlays/local`로 배포.
- DaemonSet YAML 작성:
  - hostPath로 `/proc` -> `/host/proc` 마운트.
- Hub Server는 TLS 인증서 Secret 마운트.
- 배포 후 `kubectl logs`로 확인.

8. 기술적 리스크 및 대응
- /proc 경로 오류:
  - gopsutil이 컨테이너 내부 /proc만 읽지 않도록 HOST_PROC 설정 확인.
- Redis Service 발견 실패:
  - Service 이름(`redis-service`) 및 네임스페이스 일치 확인.
- WebSocket 연결 실패:
  - Hub Server 서비스/로드밸런서 포트 매핑 확인.
- TLS 설정 오류:
  - self-signed 인증서 생성/Secret 마운트 경로 확인.

9. 프로덕션급 패키지 구조
9.1 Go 디렉터리 구조 (단일 모듈)
- cmd/agent (main)
- cmd/hub (main)
- internal/agent (수집/직렬화/전송)
- internal/hub (redis/subscriber/ws/broadcast)
- internal/config (env/flag 파서)
- internal/metrics (공통 모델)
- pkg/health (헬스체크 유틸)

9.2 K8s 매니페스트 구조
- deploy/base (공통 리소스)
- deploy/overlays/local (Orbstack 로컬 환경)
- deploy/overlays/prod (추후 확장)

10. 범위 외(향후 확장)
- 장기 저장(TSDB) 연동.
- 알림/슬랙 통합.
- 사용자 인증/권한 관리.
