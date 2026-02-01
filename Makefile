# ============================================
# Docker Image Build
# ============================================
IMAGE_NAME ?= nhandd/stream-server
IMAGE_TAG ?= latest

image-build:
	docker build -t $(IMAGE_NAME):$(IMAGE_TAG) .

image-build-no-cache:
	docker build --no-cache -t $(IMAGE_NAME):$(IMAGE_TAG) .

image-push:
	docker push $(IMAGE_NAME):$(IMAGE_TAG)

image-run:
	@mkdir -p ./logs/supervisor ./logs/app ./data
	docker run -d --name stream-server \
		-p 1935:1935 \
		-p 1985:1985 \
		-p 18080:8080 \
		-p 10080:10080/udp \
		-p 10081:10081 \
		-e NDD_LOG_PATH="/app/logs" \
		-v $(CURDIR)/logs/supervisor:/var/log/supervisor \
		-v $(CURDIR)/logs/app:/app/logs \
		-v $(CURDIR)/data:/app/data \
		$(IMAGE_NAME):$(IMAGE_TAG)

image-run-custom-env:
	@mkdir -p ./logs/supervisor ./logs/app ./data
	docker run -d --name stream-server \
		-p 1935:1935 \
		-p 1985:1985 \
		-p 18080:8080 \
		-p 10080:10080/udp \
		-p 10081:10081 \
		-e NDD_DBT="$(NDD_DBT)" \
		-e NDD_PROXY_ADDR="$(NDD_PROXY_ADDR)" \
		-e NDD_LOG_PATH="/app/logs" \
		-v $(CURDIR)/logs/supervisor:/var/log/supervisor \
		-v $(CURDIR)/logs/app:/app/logs \
		-v $(CURDIR)/data:/app/data \
		$(IMAGE_NAME):$(IMAGE_TAG)

image-stop:
	docker stop stream-server && docker rm stream-server

image-logs:
	docker logs -f stream-server

# ============================================
# Docker Compose (for development)
# ============================================
compose-up:
	sudo docker compose up -d

compose-down:
	sudo docker compose down

compose-restart:
	sudo docker compose down
	sudo docker compose up -d

compose-logs:
	sudo docker compose logs -f