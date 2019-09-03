GIT_COMMIT := $(shell git rev-list -1 HEAD)

include .env

.EXPORT_ALL_VARIABLES:
GO111MODULE=on
CGO_ENABLED=0

.PHONY: build healthcheck server

all: build-img build-int-img

build-src-img:
	docker build \
		-f build/src.Dockerfile \
		-t $(PROJECT)-src \
		.
		
build-img: build-src-img
	docker build \
		--build-arg PROJECT=$(PROJECT) \
		--build-arg VERSION=$(GIT_COMMIT) \
		-f build/build.Dockerfile \
		-t $(PROJECT) \
		.

build-int-img: build-src-img
	docker build \
		--build-arg PROJECT=$(PROJECT) \
		-f ./build/int.Dockerfile \
		-t $(PROJECT)-integration \
		.

run: build-img stop down
	docker-compose \
		--project-directory ./build/local \
		-f ./build/local/docker-compose.yml \
		-f ./build/local/docker-compose.dependencies.yml \
		-p $(PROJECT) \
		up --force-recreate --remove-orphans

run-integration: build-img build-int-img stop down
	docker-compose \
		--project-directory ./build/local \
		-f ./build/local/docker-compose.yml \
		-f ./build/local/docker-compose.dependencies.yml \
		-f ./build/local/docker-compose.integration.yml \
		-p $(PROJECT)-integration \
		up --force-recreate --remove-orphans --build --abort-on-container-exit --exit-code-from test

run-integration-only: build-int-img stop
	docker run -it --rm \
		--network="$(PROJECT)_default" \
		-e SERVER_URL=http://application:1323 \
		$(PROJECT)-integration

down:
	docker-compose \
		--project-directory ./build/local \
		-f ./build/local/docker-compose.yml \
		-f ./build/local/docker-compose.dependencies.yml \
		-p $(PROJECT) \
		down --remove-orphans

	docker-compose \
		--project-directory ./build/local \
		-f ./build/local/docker-compose.yml \
		-f ./build/local/docker-compose.dependencies.yml \
		-f ./build/local/docker-compose.integration.yml \
		-p $(PROJECT)-integration \
		down --remove-orphans
	
config:
	docker-compose \
		--project-directory ./build/local \
		-f ./build/local/docker-compose.yml \
		-f ./build/local/docker-compose.dependencies.yml \
		-p $(PROJECT) \
		config

logs:
	docker-compose \
		--project-directory ./build/local \
		-f ./build/local/docker-compose.yml \
		-f ./build/local/docker-compose.dependencies.yml \
		-p $(PROJECT) \
		logs

name:
	@echo $(PROJECT)

version:
	@echo $(GIT_COMMIT)

test:
	go test ./... -cover
	
test-all:
	go test ./... -count=1 -cover
	
cover:
	go test ./... -coverprofile=cover.out --covermode=atomic && go tool cover -html=cover.out

integration:
	go test -c -tags "integration debug" ./cmd/server/integration && ./integration.test -test.v

build: lint
	rm -f ./healthcheck ./healthcheck.exe && go build -a -ldflags="-s -w -X main.version=$(GIT_COMMIT)"	./cmd/healthcheck
	rm -f ./server ./server.exe 		  && go build -a -ldflags="-s -w -X main.version=$(GIT_COMMIT)" ./cmd/server

dev: start
	rm -f ./server ./server.exe && go build ./cmd/server && ./server

d:
	rm -f ./server ./server.exe && go build ./cmd/server && ./server

dev-hot: start
	fresh

healthcheck:
	go run ./cmd/healthcheck

db:
	docker start db || \
	docker run -itd --rm --name db -p 5432:5432 \
	    --env-file=.env \
		postgres:alpine

start: down db

stop:
	docker stop db || (exit 0)

docker-reset: docker-prune docker-rm docker-rmi docker-rmv
	docker system prune --force -a

docker-prune:
	docker system prune --force -a

docker-rm:
	docker rm $$(docker ps -a -q) || (exit 0)

docker-rmi:
	docker rmi $$(docker images -q) --force || (exit 0)

docker-rmv:
	docker volume rm $$(docker volume ls -q) || (exit 0)

dive:
	docker run --rm -it \
		-v /var/run/docker.sock:/var/run/docker.sock \
		wagoodman/dive:latest $(PROJECT)

dive-src:
	docker run --rm -it \
		-v /var/run/docker.sock:/var/run/docker.sock \
		wagoodman/dive:latest $(PROJECT)-src

dive-integration:
	docker run --rm -it \
		-v /var/run/docker.sock:/var/run/docker.sock \
		wagoodman/dive:latest $(PROJECT)-integration

weight:
	goweight ./cmd/server

deps: export GO111MODULE=off
deps: 
	go get github.com/codegangsta/gin
	go get github.com/jondot/goweight
	go get github.com/psampaz/go-mod-outdated

	GO111MODULE=on go mod tidy

tidy:
	go mod tidy

outdated:
	go list -u -m -json all | go-mod-outdated -update -direct 

lint:
	golangci-lint run

ctop:
	docker run -it --rm \
		--name ctop \
		-v /var/run/docker.sock:/var/run/docker.sock:ro \
		quay.io/vektorlab/ctop:latest

dry:
	docker run -it --rm \
		--name dry \
		-v /var/run/docker.sock:/var/run/docker.sock \
		moncho/dry

lazy:
	docker run -it --rm \
		--name lazy \
		-v /var/run/docker.sock:/var/run/docker.sock \
		lazyteam/lazydocker