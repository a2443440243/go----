# Go API Scaffold

一个基于Go语言的RESTful API脚手架框架，使用Gin、GORM、Viper等现代化技术栈构建。

## 功能特性

- ✅ 标准RESTful API路由结构
- ✅ 基础中间件（日志记录、请求验证、错误处理、CORS）
- ✅ 数据库连接配置模板（支持PostgreSQL、MySQL、SQLite）
- ✅ 统一响应格式规范
- ✅ 环境变量配置管理
- ✅ 健康检查端点
- ✅ 模块化代码组织结构
- ✅ 优雅关闭机制

## 技术栈

- **Web框架**: Gin
- **ORM**: GORM
- **配置管理**: Viper
- **日志**: Logrus
- **数据库**: PostgreSQL/MySQL/SQLite

## 项目结构

```
.
├── cmd/
│   └── server/
│       └── main.go          # 应用程序入口
├── internal/
│   ├── app/                 # 应用程序核心
│   ├── config/              # 配置管理
│   ├── handler/             # HTTP处理器
│   ├── middleware/          # 中间件
│   ├── model/               # 数据模型
│   ├── repository/          # 数据访问层
│   ├── router/              # 路由配置
│   └── service/             # 业务逻辑层
├── pkg/
│   ├── database/            # 数据库连接
│   ├── logger/              # 日志工具
│   └── response/            # 响应工具
├── data/                    # 数据文件目录
├── bin/                     # 编译输出目录
├── app.yaml                 # 配置文件
├── go.mod                   # Go模块文件
└── README.md                # 项目文档
```

## 快速开始

### 1. 克隆项目

```bash
git clone <repository-url>
cd go-api-scaffold
```

### 2. 安装依赖

```bash
go mod tidy
```

### 3. 配置环境

复制并编辑配置文件：

```bash
cp app.yaml.example app.yaml
```

### 4. 构建项目

```bash
go build -o bin/go-api-scaffold.exe cmd/server/main.go
```

### 5. 运行应用

```bash
./bin/go-api-scaffold.exe
```

应用将在 `http://localhost:8080` 启动。

## API 端点

### 健康检查

- `GET /api/v1/health/check` - 基础健康检查
- `GET /api/v1/health/ready` - 就绪检查
- `GET /api/v1/health/live` - 存活检查

### 示例API

- `GET /api/v1/ping` - 简单的ping测试
- `GET /api/v1/version` - 获取应用版本信息

## 配置说明

配置文件 `app.yaml` 包含以下主要配置项：

```yaml
app:
  name: "go-api-scaffold"
  version: "1.0.0"
  environment: "development"

server:
  port: 8080
  read_timeout: "30s"
  write_timeout: "30s"
  idle_timeout: "60s"

database:
  driver: "sqlite"  # postgres, mysql, sqlite
  host: "localhost"
  port: 5432
  dbname: "scaffold"
  username: "postgres"
  password: "password"
  sslmode: "disable"
  sqlite:
    path: "./data/app.db"

log:
  level: "info"  # debug, info, warn, error, fatal
```

## 开发指南

### 添加新的API端点

1. 在 `internal/handler/` 中创建处理器
2. 在 `internal/service/` 中实现业务逻辑
3. 在 `internal/repository/` 中实现数据访问
4. 在路由中注册新端点

### 数据库迁移

项目支持多种数据库，可以通过修改配置文件中的 `database.driver` 来切换：

- `postgres` - PostgreSQL
- `mysql` - MySQL
- `sqlite` - SQLite（默认）

### 中间件

项目包含以下中间件：

- **Logger**: 请求日志记录
- **Recovery**: 异常恢复
- **CORS**: 跨域资源共享

## 部署

### Docker部署

```dockerfile
FROM golang:1.21-alpine AS builder
WORKDIR /app
COPY . .
RUN go mod tidy
RUN go build -o bin/go-api-scaffold cmd/server/main.go

FROM alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /root/
COPY --from=builder /app/bin/go-api-scaffold .
COPY --from=builder /app/app.yaml .
CMD ["./go-api-scaffold"]
```

### 生产环境配置

生产环境建议：

1. 使用环境变量覆盖敏感配置
2. 启用HTTPS
3. 配置反向代理（Nginx）
4. 设置适当的日志级别
5. 配置监控和告警

## 许可证

MIT License

## 贡献

欢迎提交Issue和Pull Request！