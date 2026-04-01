.PHONY: build run test clean docker help swag

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
	@echo "  make docker   - Build Docker image"
	@echo "  make up       - Start all services with Docker Compose"
	@echo "  make down     - Stop all services"
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

# Docker构建
docker:
	docker build -t $(DOCKER_IMAGE) -f deployments/docker/Dockerfile .

# Docker Compose启动
docker-up:
	docker-compose up -d

# Docker Compose停止
docker-down:
	docker-compose down

# Docker Compose停止并删除数据卷
docker-down-v:
	docker-compose down -v

# 查看日志
logs:
	docker-compose logs -f app

# 进入容器
shell:
	docker-compose exec app sh

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
