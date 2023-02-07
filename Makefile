IMAGE_NAME=drl-exporter
IMAGE_VERSION=2.1.0

.PHONY: docker
docker:
	go mod vendor
	docker build -t $(IMAGE_NAME) .
	docker tag $(IMAGE_NAME):latest $(IMAGE_NAME):$(IMAGE_VERSION)

.PHONY: docker-multi
docker-multi:
	docker buildx build \
	--push \
	--platform linux/arm/v7,linux/arm64/v8,linux/amd64 \
	 --tag khaliq/drl-exporter:$(IMAGE_VERSION) .

.PHONY: update-go-deps
update-go-deps:
	@echo ">> updating Go dependencies"
	@for m in $$(go list -mod=readonly -m -f '{{ if and (not .Indirect) (not .Main)}}{{.Path}}{{end}}' all); do \
		go get $$m; \
	done
	go mod tidy
ifneq (,$(wildcard vendor))
	go mod vendor
endif
