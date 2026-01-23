IMAGE_AGENT ?= resource-checker/agent:dev
IMAGE_HUB ?= resource-checker/hub:dev
PLATFORM ?= linux/arm64

.PHONY: tls docker-agent docker-hub k8s-apply k8s-delete

tls:
	./scripts/gen-tls-secret.sh

docker-agent:
	docker build --platform $(PLATFORM) -f Dockerfile.agent -t $(IMAGE_AGENT) .

docker-hub:
	docker build --platform $(PLATFORM) -f Dockerfile.hub -t $(IMAGE_HUB) .

k8s-apply:
	kubectl apply -k deploy/overlays/local

k8s-delete:
	kubectl delete -k deploy/overlays/local
