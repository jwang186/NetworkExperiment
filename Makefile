OSARCH = $(shell arch)
ifeq ($(OSARCH), x86_64)
ARCH=amd64
else ifeq ($(OSARCH), aarch64)
ARCH=arm64
endif

agent:
	GOPATH=$(shell pwd)
	CGO_ENABLED=0 GOOS=linux GOARCH=$(ARCH) go build udsExperiment.go

.PHONY: docker-build
docker-build: agent
	docker build -t uds-exp --build-arg uwsgi_uid=$(id -u) --build-arg uwsgi_gid=$(id -g) .

.PHONY: run
run: docker-build
	mkdir -p tmp
	cp envoy-uds-admin.yaml $(shell pwd)/tmp
	docker run --rm \
		-v $(shell pwd)/tmp:/tmp \
		uds-exp

.PHONY: clean
clean:
	$(RM) udsExperiment
	rm -rf $(shell pwd)/tmp
	docker rmi -f uds-exp



