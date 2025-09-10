# Go API Scaffold Makefile

# 变量定义
APP_NAME := go-api-scaffold
APP_VERSION := 1.0.0
BUILD_DIR := bin
MAIN_FILE := cmd/server/main.go
DOCKER_IMAGE := $(APP_NAME):$(APP_VERSION)
DOCKER_LATEST := $(APP_NAME):latest

# Go 相关变量
GOCMD := go
GOBUILD := $(GOCMD) build
GOCLEAN := $(GOCMD) clean
GOTEST := $(GOCMD) test
GOGET := $(GOCMD) get
GOMOD := $(GOCMD) mod
GOFMT := gofmt
GOVET := $(GOCMD) vet

# 操作系统检测
ifeq ($(OS),Windows_NT)
    BINARY_EXT := .exe
    RM := del /Q
    MKDIR := mkdir
else
    BINARY_EXT :=
    RM := rm -f
    MKDIR := mkdir -p
endif

BINARY_NAME := $(APP_NAME)$(BINARY_EXT)
BINARY_PATH := $(BUILD_DIR)/$(BINARY_NAME)

# 默认目标
.PHONY: all
all: clean deps fmt vet test build

# 帮助信息
.PHONY: help
help:
	@echo "Available commands:"
	@echo "  build       - Build the application"
	@echo "  run         - Run the application"
	@echo "  test        - Run tests"
	@echo "  clean       - Clean build artifacts"
	@echo "  deps        - Download dependencies"
	@echo "  fmt         - Format code"
	@echo "  vet         - Run go vet"
	@echo "  docker-build - Build Docker image"
	@echo "  docker-run   - Run Docker container"
	@echo "  compose-up   - Start services with docker-compose"
	@echo "  compose-down - Stop services with docker-compose"

# 构建应用
.PHONY: build
build:
	@echo "Building $(APP_NAME)..."
	@$(MKDIR) $(BUILD_DIR) 2>/dev/null || true
	@$(GOBUILD) -o $(BINARY_PATH) -v $(MAIN_FILE)
	@echo "Build completed: $(BINARY_PATH)"

# 运行应用
.PHONY: run
run:
	@echo "Running $(APP_NAME)..."
	@$(GOCMD) run $(MAIN_FILE)

# 运行测试
.PHONY: test
test:
	@echo "Running tests..."
	@$(GOTEST) -v ./...

# 清理构建产物
.PHONY: clean
clean:
	@echo "Cleaning..."
	@$(GOCLEAN)
	@$(RM) $(BUILD_DIR)/* 2>/dev/null || true
	@echo "Clean completed"

# 下载依赖
.PHONY: deps
deps:
	@echo "Downloading dependencies..."
	@$(GOGET) -d -v ./...
	@$(GOMOD) tidy
	@$(GOMOD) verify

# 格式化代码
.PHONY: fmt
fmt:
	@echo "Formatting code..."
	@$(GOFMT) -s -w .

# 代码检查
.PHONY: vet
vet:
	@echo "Running go vet..."
	@$(GOVET) ./...

# Docker 相关命令
.PHONY: docker-build
docker-build:
	@echo "Building Docker image..."
	@docker build -t $(DOCKER_IMAGE) -t $(DOCKER_LATEST) .
	@echo "Docker image built: $(DOCKER_IMAGE)"

.PHONY: docker-run
docker-run:
	@echo "Running Docker container..."
	@docker run -d \
		--name $(APP_NAME) \
		-p 8080:8080 \
		-e APP_ENV=development \
		$(DOCKER_LATEST)
	@echo "Container started: $(APP_NAME)"

.PHONY: docker-stop
docker-stop:
	@echo "Stopping Docker container..."
	@docker stop $(APP_NAME) || true
	@docker rm $(APP_NAME) || true

.PHONY: docker-logs
docker-logs:
	@docker logs -f $(APP_NAME)

# Docker Compose 命令
.PHONY: compose-up
compose-up:
	@echo "Starting services with Docker Compose..."
	@docker-compose up -d

.PHONY: compose-down
compose-down:
	@echo "Stopping services with Docker Compose..."
	@docker-compose down

.PHONY: compose-logs
compose-logs:
	@docker-compose logs -f

.PHONY: compose-build
compose-build:
	@echo "Building services with Docker Compose..."
	@docker-compose build

# 健康检查
.PHONY: health
health:
	@echo "Checking application health..."
	@curl -f http://localhost:8080/api/v1/health/check || echo "Application is not running or unhealthy"

# 显示版本信息
.PHONY: version
version:
	@echo "$(APP_NAME) version $(APP_VERSION)"