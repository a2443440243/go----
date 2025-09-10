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
GOLINT := golangci-lint

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
	@echo "  build      - Build the application"
	@echo "  run        - Run the application"
	@echo "  test       - Run tests"
	@echo "  test-cover - Run tests with coverage"
	@echo "  clean      - Clean build artifacts"
	@echo "  deps       - Download dependencies"
	@echo "  fmt        - Format code"
	@echo "  vet        - Run go vet"
	@echo "  lint       - Run linter"
	@echo "  docker-build - Build Docker image"
	@echo "  docker-run   - Run Docker container"
	@echo "  docker-push  - Push Docker image"
	@echo "  dev        - Run in development mode"
	@echo "  prod       - Build for production"
	@echo "  install    - Install the application"
	@echo "  uninstall  - Uninstall the application"

# 构建应用
.PHONY: build
build:
	@echo "Building $(APP_NAME)..."
	@$(MKDIR) $(BUILD_DIR) 2>/dev/null || true
	@$(GOBUILD) -o $(BINARY_PATH) -v $(MAIN_FILE)
	@echo "Build completed: $(BINARY_PATH)"

# 生产环境构建
.PHONY: prod
prod:
	@echo "Building $(APP_NAME) for production..."
	@$(MKDIR) $(BUILD_DIR) 2>/dev/null || true
	@CGO_ENABLED=0 GOOS=linux $(GOBUILD) -a -installsuffix cgo -ldflags '-extldflags "-static"' -o $(BINARY_PATH) $(MAIN_FILE)
	@echo "Production build completed: $(BINARY_PATH)"

# 交叉编译
.PHONY: build-all
build-all:
	@echo "Building for multiple platforms..."
	@$(MKDIR) $(BUILD_DIR) 2>/dev/null || true
	@GOOS=linux GOARCH=amd64 $(GOBUILD) -o $(BUILD_DIR)/$(APP_NAME)-linux-amd64 $(MAIN_FILE)
	@GOOS=windows GOARCH=amd64 $(GOBUILD) -o $(BUILD_DIR)/$(APP_NAME)-windows-amd64.exe $(MAIN_FILE)
	@GOOS=darwin GOARCH=amd64 $(GOBUILD) -o $(BUILD_DIR)/$(APP_NAME)-darwin-amd64 $(MAIN_FILE)
	@GOOS=darwin GOARCH=arm64 $(GOBUILD) -o $(BUILD_DIR)/$(APP_NAME)-darwin-arm64 $(MAIN_FILE)
	@echo "Cross-compilation completed"

# 运行应用
.PHONY: run
run:
	@echo "Running $(APP_NAME)..."
	@$(GOCMD) run $(MAIN_FILE)

# 开发模式运行（带热重载）
.PHONY: dev
dev:
	@echo "Running $(APP_NAME) in development mode..."
	@if command -v air > /dev/null; then \
		air; \
	else \
		echo "Air not found. Installing..."; \
		$(GOGET) -u github.com/cosmtrek/air; \
		air; \
	fi

# 运行测试
.PHONY: test
test:
	@echo "Running tests..."
	@$(GOTEST) -v ./...

# 运行测试并生成覆盖率报告
.PHONY: test-cover
test-cover:
	@echo "Running tests with coverage..."
	@$(GOTEST) -v -coverprofile=coverage.out ./...
	@$(GOCMD) tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report generated: coverage.html"

# 基准测试
.PHONY: bench
bench:
	@echo "Running benchmarks..."
	@$(GOTEST) -bench=. -benchmem ./...

# 清理构建产物
.PHONY: clean
clean:
	@echo "Cleaning..."
	@$(GOCLEAN)
	@$(RM) $(BUILD_DIR)/* 2>/dev/null || true
	@$(RM) coverage.out coverage.html 2>/dev/null || true
	@echo "Clean completed"

# 下载依赖
.PHONY: deps
deps:
	@echo "Downloading dependencies..."
	@$(GOGET) -d -v ./...
	@$(GOMOD) tidy
	@$(GOMOD) verify

# 更新依赖
.PHONY: deps-update
deps-update:
	@echo "Updating dependencies..."
	@$(GOGET) -u ./...
	@$(GOMOD) tidy

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

# 代码规范检查
.PHONY: lint
lint:
	@echo "Running linter..."
	@if command -v $(GOLINT) > /dev/null; then \
		$(GOLINT) run; \
	else \
		echo "golangci-lint not found. Please install it first."; \
		echo "Visit: https://golangci-lint.run/usage/install/"; \
	fi

# 安装 linter
.PHONY: install-lint
install-lint:
	@echo "Installing golangci-lint..."
	@curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(shell go env GOPATH)/bin v1.54.2

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

.PHONY: docker-push
docker-push:
	@echo "Pushing Docker image..."
	@docker push $(DOCKER_IMAGE)
	@docker push $(DOCKER_LATEST)

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

# 数据库相关命令
.PHONY: db-migrate
db-migrate:
	@echo "Running database migrations..."
	@$(GOCMD) run $(MAIN_FILE) migrate

.PHONY: db-seed
db-seed:
	@echo "Seeding database..."
	@$(GOCMD) run $(MAIN_FILE) seed

.PHONY: db-reset
db-reset:
	@echo "Resetting database..."
	@$(GOCMD) run $(MAIN_FILE) reset

# 安装应用到系统
.PHONY: install
install: build
	@echo "Installing $(APP_NAME)..."
	@cp $(BINARY_PATH) /usr/local/bin/$(APP_NAME) 2>/dev/null || \
		echo "Please run 'sudo make install' to install to /usr/local/bin"

# 卸载应用
.PHONY: uninstall
uninstall:
	@echo "Uninstalling $(APP_NAME)..."
	@$(RM) /usr/local/bin/$(APP_NAME) 2>/dev/null || \
		echo "Please run 'sudo make uninstall' to remove from /usr/local/bin"

# 生成 API 文档
.PHONY: docs
docs:
	@echo "Generating API documentation..."
	@if command -v swag > /dev/null; then \
		swag init -g $(MAIN_FILE); \
	else \
		echo "Swagger not found. Installing..."; \
		$(GOGET) -u github.com/swaggo/swag/cmd/swag; \
		swag init -g $(MAIN_FILE); \
	fi

# 代码生成
.PHONY: generate
generate:
	@echo "Running go generate..."
	@$(GOCMD) generate ./...

# 安全检查
.PHONY: security
security:
	@echo "Running security checks..."
	@if command -v gosec > /dev/null; then \
		gosec ./...; \
	else \
		echo "gosec not found. Installing..."; \
		$(GOGET) github.com/securecodewarrior/gosec/v2/cmd/gosec; \
		gosec ./...; \
	fi

# 性能分析
.PHONY: profile
profile:
	@echo "Running performance profiling..."
	@$(GOCMD) test -cpuprofile cpu.prof -memprofile mem.prof -bench . ./...
	@echo "Profile files generated: cpu.prof, mem.prof"

# 检查更新
.PHONY: check-updates
check-updates:
	@echo "Checking for dependency updates..."
	@$(GOCMD) list -u -m all

# 完整的 CI 流程
.PHONY: ci
ci: clean deps fmt vet lint test security
	@echo "CI pipeline completed successfully"

# 发布准备
.PHONY: release
release: ci build-all
	@echo "Release artifacts prepared in $(BUILD_DIR)/"

# 健康检查
.PHONY: health
health:
	@echo "Checking application health..."
	@curl -f http://localhost:8080/health || echo "Application is not running or unhealthy"

# 显示版本信息
.PHONY: version
version:
	@echo "$(APP_NAME) version $(APP_VERSION)"

# 显示项目统计
.PHONY: stats
stats:
	@echo "Project Statistics:"
	@echo "Lines of code:"
	@find . -name '*.go' -not -path './vendor/*' | xargs wc -l | tail -1
	@echo "Number of Go files:"
	@find . -name '*.go' -not -path './vendor/*' | wc -l
	@echo "Dependencies:"
	@$(GOMOD) graph | wc -l

# 监控模式运行
.PHONY: watch
watch:
	@echo "Running in watch mode..."
	@if command -v fswatch > /dev/null; then \
		fswatch -o . | xargs -n1 -I{} make run; \
	elif command -v inotifywait > /dev/null; then \
		while inotifywait -r -e modify .; do make run; done; \
	else \
		echo "No file watcher found. Please install fswatch or inotify-tools"; \
	fi