APP_NAME = "konterfai"
OLLAMA_MODEL ?= "qwen2:0.5b"
DOCKER_TAG ?= "dev"
BROWSER ?= "firefox"

.PHONY: all
all: build

.PHONY: all-arch-build
all-arch-build: build-amd64 build-arm64 build-armv7

.PHONY: build
build:
	@echo "Building..."
	@go build -o bin/$(APP_NAME) cmd/$(APP_NAME)/main.go
	@echo "Done."

.PHONY: test
test:
	@echo "Testing..."
	@go test -v ./...
	@echo "Done."

.PHONY: coverage
coverage:
	@echo "Coverage..."
	@go test -v ./... -coverprofile=coverage.out
	@go tool cover -html=coverage.out -o coverage.html
	@$(BROWSER) coverage.html
	@echo "Done."

.PHONY: coverage-ci
coverage-ci:
	@echo "Coverage..."
	@go test -v ./... -coverprofile=coverage.out
	@go tool cover -html=coverage.out -o coverage.html
	@echo "Done."

.PHONY: golangci-lint
golangci-lint:
	@echo "Linting..."
	@golangci-lint run
	@echo "Done."

.PHONY: clean
clean:
	@rm -rf bin

.PHONY: start-ollama
start-ollama: stop-ollama
	@docker run -d -v ollama:/root/.ollama -p 11434:11434 --name ollama ollama/ollama:latest
	@docker exec ollama ollama run $(OLLAMA_MODEL)

.PHONY: start-ollama-gpu
start-ollama-gpu: stop-ollama
	@docker run -d -v ollama:/root/.ollama -p 11434:11434 --gpus all --name ollama ollama/ollama:latest
	@docker exec ollama ollama run $(OLLAMA_MODEL)

.PHONY : stop-ollama
stop-ollama:
	-docker stop ollama
	-docker rm ollama

.PHONY: run
run:
	@echo "Running..."
	@go run -race cmd/$(APP_NAME)/main.go --ollama-model $(OLLAMA_MODEL)

.PHONY: docker-build
docker-build:
	@docker buildx build -t $(APP_NAME)/$(APP_NAME):$(DOCKER_TAG) .

.PHONY: docker-run
docker-run:
	@docker run --net=host -e OLLAMA_MODEL=$(OLLAMA_MODEL) -p 8080:8080 --name ${APP_NAME} $(APP_NAME)/$(APP_NAME):$(DOCKER_TAG)

.PHONY: docker-stop
docker-stop:
	-docker stop $(APP_NAME)
	-docker rm $(APP_NAME)

.PHONY: docker-compose-up
docker-compose-up: docker-build
	@docker-compose -f docker-compose-dev.yml up -d

.PHONY: docker-compose-down
docker-compose-down:
	@docker-compose -f docker-compose-dev.yml down

.PHONY: build-amd64
build-amd64:
	@echo "Building for amd64..."
	@GOOS=linux GOARCH=amd64 go build -o bin/$(APP_NAME)_amd64 cmd/$(APP_NAME)/main.go
	@echo "Done."

.PHONY: build-arm64
build-arm64:
	@echo "Building for arm64..."
	@GOOS=linux GOARCH=arm64 go build -o bin/$(APP_NAME)_arm64 cmd/$(APP_NAME)/main.go
	@echo "Done."

build-armv7:
	@echo "Building for armv7..."
	@GOOS=linux GOARCH=arm go build -o bin/$(APP_NAME)_armv7 cmd/$(APP_NAME)/main.go
	@echo "Done."
