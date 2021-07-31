BIN := "./bin/banners-rotation"
DOCKER_IMG="banners-rotation:main"

GIT_HASH := $(shell git log --format="%h" -n 1)
LDFLAGS := -X main.release="develop" -X main.buildDate=$(shell date -u +%Y-%m-%dT%H:%M:%S) -X main.gitHash=$(GIT_HASH)

build:
	go build -v -o $(BIN) -ldflags "$(LDFLAGS)" ./cmd/banners-rotation

build-img:
	docker build \
		--build-arg=LDFLAGS="$(LDFLAGS)" \
		-t $(DOCKER_IMG) \
		-f build/Dockerfile .

run-img: build-img
	docker run $(DOCKER_IMG)

version: build
	$(BIN) version

test:
	go test -race -v ./internal/... -count 100
	
t:
	go test -race -v ./internal/...

install-lint-deps:
	(which golangci-lint > /dev/null) || curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(shell go env GOPATH)/bin v1.37.0

lint: install-lint-deps
	golangci-lint run --config ./configs/golangci.yml

generate-deps:
	mkdir -p internal/server/pb

generate: generate-deps
	protoc -I ./proto \
    --go_out ./internal/server/pb/ --go_opt paths=source_relative \
    --go-grpc_out ./internal/server/pb/ --go-grpc_opt paths=source_relative \
    proto/BannersRotation/*.proto

generate-gateway: generate
	protoc -I ./proto --grpc-gateway_out ./internal/server/pb \
			--grpc-gateway_opt logtostderr=true \
			--grpc-gateway_opt paths=source_relative \
			--grpc-gateway_opt generate_unbound_methods=true \
			proto/BannersRotation/*.proto

dev-build-container:
	docker rm --force banners-rotation-br
	docker rm --force postgres-br
	docker rm --force rabbit-br
	docker-compose build --no-cache

run:
	docker-compose up -d

stop:
	docker-compose down

dev: dev-build-container
	docker-compose up

up:
	docker-compose up --build

integration-tests-cleanup:
	docker-compose -f docker-compose.test.yaml down \
    --rmi local \
		--volumes \
		--remove-orphans \
		--timeout 60; \
  	docker-compose rm -f

integration-tests:
	set -e ;\
	docker-compose -f docker-compose.test.yaml up --build -d ;\
	test_status_code=0 ;\
	docker-compose -f docker-compose.test.yaml run integration_tests go test -v || test_status_code=$$? ;\
	docker-compose -f docker-compose.test.yaml down ;\
	exit $$test_status_code ;

.PHONY: build build-img run-img version test lint install-lint-deps generate-deps generate generate-gateway upl run stop dev integration-tests-cleanup integration-tests
