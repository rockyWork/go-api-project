# Go API Project

一个基于 Go + Gin + GORM + MySQL + JWT 的现代化 Web API 服务模板，包含用户认证和管理功能。

## 功能特性

- **JWT 认证**: Access Token + Refresh Token 双令牌机制
- **用户管理**: 用户注册、登录、CRUD、权限控制
- **分层架构**: Controller → Service → Repository 三层架构
- **API 文档**: 自动生成 Swagger 文档
- **容器化部署**: Podman + Podman Compose（推荐） 或 Docker 一键部署
- **日志系统**: 结构化日志（Zap）
- **配置管理**: 多环境配置文件
- **扩展预留**: 预留 Redis、RBAC、操作日志等扩展接口

## 技术栈

| 组件 | 技术 |
|------|------|
| Web 框架 | Gin |
| ORM | GORM |
| 数据库 | MySQL 8.0 |
| 缓存 | Redis（预留） |
| 认证 | JWT |
| 配置 | Viper |
| 日志 | Zap |
| 文档 | Swagger |
| 部署 | Docker + Docker Compose |

## 项目结构

```
go-api-project/
├── api/                    # API 层
│   ├── v1/                 # API v1 版本
│   │   ├── auth.go         # 认证接口
│   │   └── user.go         # 用户接口
│   └── router.go           # 路由注册
├── internal/               # 内部模块
│   ├── model/              # 数据模型
│   ├── service/            # 业务逻辑
│   ├── repository/         # 数据访问
│   └── middleware/         # 中间件
├── pkg/                    # 公共包
│   ├── jwt/                # JWT 工具
│   ├── response/           # 统一响应
│   ├── logger/             # 日志工具
│   └── utils/              # 工具函数
├── config/                 # 配置文件
├── deployments/docker/     # 容器配置（Podman/Docker）
├── docs/                   # 文档
├── scripts/                # 脚本
└── main.go                 # 入口文件
```

## 快速开始

### 环境要求

- Go 1.21+
- MySQL 8.0+
- Podman & Podman Compose（推荐）或 Docker & Docker Compose

### 1. 克隆项目

```bash
git clone <your-repo-url>
cd go-api-project
```

### 2. 安装依赖

```bash
go mod download
```

###  安装choco
第一步：安装 Chocolatey（包管理器）
请以管理员身份打开 PowerShell，然后逐条复制以下命令：

Set-ExecutionPolicy Bypass -Scope Process -Force; 

[System.Net.ServicePointManager]::SecurityProtocol = [System.Net.ServicePointManager]::SecurityProtocol -bor 3072; 

iex ((New-Object System.Net.WebClient).DownloadString('https://community.chocolatey.org/install.ps1'))

### 安装make
choco install make -y


### 3. 配置数据库

修改 `config/config.yaml` 中的数据库配置：

```yaml
database:
  host: localhost      # Docker 部署时使用 mysql
  port: 3306
  username: root
  password: root123
  dbname: go_api_project
```

### 4. 启动服务

#### 本地开发模式

```bash
# 方式 1：直接运行
go run main.go

# 方式 2：使用 Makefile
make run
```

#### Podman 部署（推荐）

```bash
# 启动所有服务（应用 + MySQL）
podman-compose up -d
# 或使用 Makefile
make podman-up

# 查看日志
podman-compose logs -f app
# 或
make logs

# 停止服务
podman-compose down
# 或
make podman-down

# 停止并删除数据卷
podman-compose down -v
```

#### Docker 兼容模式

如果没有安装 Podman，也可以使用原 Docker 命令：

```bash
make docker-up
```

### 5. 验证服务

```bash
# 健康检查
curl http://localhost:8080/health

# Swagger 文档
open http://localhost:8080/swagger/index.html
```

## API 接口

### 认证接口

| 接口 | 方法 | 描述 | 认证 |
|------|------|------|------|
| `/api/v1/auth/register` | POST | 用户注册 | 否 |
| `/api/v1/auth/login` | POST | 用户登录 | 否 |
| `/api/v1/auth/refresh` | POST | 刷新 Token | 否 |
| `/api/v1/auth/logout` | POST | 用户登出 | 是 |

### 用户接口

| 接口 | 方法 | 描述 | 认证 | 权限 |
|------|------|------|------|------|
| `/api/v1/users/me` | GET | 当前用户信息 | 是 | 用户/管理员 |
| `/api/v1/users` | GET | 用户列表 | 是 | 管理员 |
| `/api/v1/users` | POST | 创建用户 | 是 | 管理员 |
| `/api/v1/users/:id` | GET | 用户详情 | 是 | 自己/管理员 |
| `/api/v1/users/:id` | PUT | 更新用户 | 是 | 自己/管理员 |
| `/api/v1/users/:id` | DELETE | 删除用户 | 是 | 管理员 |
| `/api/v1/users/change-password` | POST | 修改密码 | 是 | 用户 |

### 请求示例

#### 注册
```bash
curl -X POST http://localhost:8080/api/v1/auth/register \
  -H "Content-Type: application/json" \
  -d '{
    "username": "testuser",
    "email": "test@example.com",
    "password": "Test123456"
  }'
```

#### 登录
```bash
curl -X POST http://localhost:8080/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "username": "testuser",
    "password": "Test123456"
  }'
```

#### 访问需要认证的接口
```bash
curl -X GET http://localhost:8080/api/v1/users/me \
  -H "Authorization: Bearer <your-access-token>"
```

## 服务列表

| 服务名 | 端口 | 描述 | 启动命令 |
|--------|------|------|----------|
| go-api-app | 8080 | Go Web 服务 | `podman-compose up -d` 或 `go run main.go` |
| go-api-mysql | 3306 | MySQL 数据库 | `podman-compose up -d` |
| go-api-redis | 6379 | Redis 缓存（预留） | `podman-compose up -d`（需取消注释） |

## 常用命令

```bash
# 查看所有可用命令
make help

# 构建应用
make build

# 运行测试
make test

# 生成 Swagger 文档
make swag

# Docker 构建
make docker

# Docker Compose 启动
make docker-up

# Docker Compose 停止
make docker-down

# 查看日志
make logs

# 进入容器
make shell

# 代码格式化
make fmt
```

## 配置说明

配置文件位于 `config/` 目录：

| 文件 | 用途 |
|------|------|
| `config.yaml` | 默认配置（开发环境） |
| `config.docker.yaml` | Docker 环境配置 |

环境变量覆盖：
- `CONFIG_PATH`: 指定配置文件路径
- 其他配置项可通过环境变量覆盖（使用 Viper 的自动环境变量绑定）

## 扩展开发

### 添加新模块

1. **定义模型** (`internal/model/`)
   ```go
   type NewModel struct { ... }
   ```

2. **创建 Repository** (`internal/repository/`)
   ```go
   type NewRepository struct { db *gorm.DB }
   ```

3. **实现 Service** (`internal/service/`)
   ```go
   type NewService struct { repo *repository.NewRepository }
   ```

4. **添加 Handler** (`api/v1/`)
   ```go
   func (h *Handler) NewEndpoint(c *gin.Context) { ... }
   ```

5. **注册路由** (`api/router.go`)
   ```go
   apiV1.GET("/new-endpoint", handler.NewEndpoint)
   ```

### 预留扩展点

- **缓存层**: `pkg/cache/`（预留 Redis 接口）
- **消息队列**: `pkg/mq/`（预留队列接口）
- **文件存储**: `pkg/storage/`（预留 OSS/S3 接口）
- **限流**: `internal/middleware/ratelimit.go`（预留）

## 安全建议

部署到生产环境前，请修改以下配置：

1. **JWT 密钥**: 修改 `config/config.yaml` 中的 `jwt.access_secret` 和 `jwt.refresh_secret`
2. **数据库密码**: 使用强密码并限制访问权限
3. **HTTPS**: 使用反向代理（Nginx/Caddy）配置 HTTPS
4. **日志级别**: 生产环境设置为 `info` 或 `warn`
5. **CORS**: 限制允许的域名，不要使用 `*`

## 常见问题

### 1. 数据库连接失败

检查 MySQL 服务是否启动，配置文件中的数据库连接信息是否正确。

```bash
# 检查容器状态
podman-compose ps
# 或
podman-compose ps

# 查看 MySQL 日志
podman-compose logs mysql
# 或
podman-compose logs mysql
```

### 2. 端口被占用

```bash
# 查看 8080 端口占用
lsof -i :8080

# 修改 compose.yml 中的端口映射
ports:
  - "8081:8080"  # 将主机 8081 映射到容器 8080
```

### 3. Swagger 文档不显示

确保已安装 swag 工具并生成文档：

```bash
go install github.com/swaggo/swag/cmd/swag@latest
swag init -g main.go
```

## 许可证

MIT License

## 贡献

欢迎提交 Issue 和 Pull Request！
