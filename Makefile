IMAGE_NAME=drl-exporter
IMAGE_VERSION=2.0

.PHONY: docker
docker:
	docker build -t $(IMAGE_NAME) .
	docker tag $(IMAGE_NAME):latest $(IMAGE_NAME):$(IMAGE_VERSION)

.PHONY: docker-multi
docker-multi:
	docker buildx build \
	--push \
	--platform linux/arm/v7,linux/arm64/v8,linux/amd64 \
	 --tag khaliq/drl-exporter:$(IMAGE_VERSION) .

