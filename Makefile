.PHONY: build run test clean docker help swag podman-up podman-down podman-down-v docker-up docker-down docker-down-v

# 变量
APP_NAME=go-api-project
DOCKER_IMAGE=$(APP_NAME):latest

# 默认目标
help:
	@echo "Available targets:"
	@echo "  make build    - Build the application"
	@echo "  make run      - Run the application locally"
	@echo "  make test     - Run tests"
	@echo "  make clean    - Clean build artifacts"
	@echo "  make docker   - Build container image"
	@echo "  make podman-up - Start services with Podman (recommended)"
	@echo "  make docker-up - Start services with Docker/Podman"
	@echo "  make swag     - Generate Swagger documentation"
	@echo "  make deps     - Download dependencies"

# 构建应用
build:
	go build -o bin/$(APP_NAME) .

# 本地运行
run:
	go run main.go

# 运行测试
test:
	go test -v ./...

# 清理
clean:
	rm -rf bin/
	go clean

# 下载依赖
deps:
	go mod download
	go mod tidy

# 生成Swagger文档
swag:
	swag init -g main.go

# 容器构建 (Podman/Docker)
build-container:
	podman build -t $(DOCKER_IMAGE) -f deployments/docker/Dockerfile . || docker build -t $(DOCKER_IMAGE) -f deployments/docker/Dockerfile .

# Docker构建 (兼容)
docker:
	$(MAKE) build-container

# Podman Compose启动 (推荐)
podman-up:
	podman-compose up -d

# Podman Compose停止
podman-down:
	podman-compose down

# Podman Compose停止并删除数据卷
podman-down-v:
	podman-compose down -v

# 兼容旧命令（Docker）
docker-up:
	podman-compose up -d || docker-compose up -d

docker-down:
	podman-compose down || docker-compose down

docker-down-v:
	podman-compose down -v || docker-compose down -v

# 查看日志
logs:
	podman-compose logs -f app || docker-compose logs -f app

# 进入容器
shell:
	podman-compose exec app sh || docker-compose exec app sh

# 数据库迁移（使用GORM自动迁移）
migrate:
	go run main.go

# 代码格式化
fmt:
	go fmt ./...

# 代码检查
lint:
	golangci-lint run

# 开发模式（热重载）
dev:
	air
