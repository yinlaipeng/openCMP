.PHONY: build run test clean docker-build docker-run help deps fmt lint docker-up docker-down

# 变量
BINARY_NAME=opencmp
CMD_DIR=cmd/server
GO=go
GOFLAGS=-v

# 默认目标
all: build

# 构建
build:
	@echo "Building $(BINARY_NAME)..."
	$(GO) build $(GOFLAGS) -o $(BINARY_NAME) $(CMD_DIR)/main.go

# 运行
run: build
	@echo "Running $(BINARY_NAME)..."
	./$(BINARY_NAME) -config configs/config.yaml

# 测试
test:
	@echo "Running tests..."
	$(GO) test $(GOFLAGS) -race -cover ./...

# 清理
clean:
	@echo "Cleaning..."
	$(GO) clean
	rm -f $(BINARY_NAME)

# Docker 构建
docker-build:
	@echo "Building Docker image..."
	docker build -t $(BINARY_NAME):latest .

# Docker 运行
docker-run:
	@echo "Running Docker container..."
	docker run -d -p 8080:8080 --name $(BINARY_NAME) $(BINARY_NAME):latest

# Docker Compose 启动 MySQL
docker-up:
	@echo "Starting MySQL with docker-compose..."
	docker-compose up -d

# Docker Compose 停止 MySQL
docker-down:
	@echo "Stopping MySQL..."
	docker-compose down

# 安装依赖
deps:
	@echo "Installing dependencies..."
	$(GO) mod download
	$(GO) mod tidy

# 代码格式化
fmt:
	@echo "Formatting code..."
	$(GO) fmt ./...

# 代码检查
lint:
	@echo "Running linter..."
	golangci-lint run || echo "golangci-lint not installed, skip"

# 帮助
help:
	@echo "Makefile commands:"
	@echo "  make build        - Build the binary"
	@echo "  make run          - Build and run"
	@echo "  make test         - Run tests"
	@echo "  make clean        - Clean build artifacts"
	@echo "  make deps         - Install dependencies"
	@echo "  make fmt          - Format code"
	@echo "  make lint         - Run linter"
	@echo "  make docker-build - Build Docker image"
	@echo "  make docker-run   - Run Docker container"
	@echo "  make docker-up    - Start MySQL with docker-compose"
	@echo "  make docker-down  - Stop MySQL"
