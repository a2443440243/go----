# Go API Scaffold - 优化版本

一个精简、高效的Go后端API脚手架，专门针对MySQL 8数据库进行了优化。

## 🚀 主要特性

- **轻量级设计**：移除了不必要的依赖和复杂功能
- **MySQL 8支持**：专门优化了MySQL 8数据库连接和配置
- **完整的用户管理**：包含用户CRUD操作和登录认证
- **RESTful API**：标准的REST API设计
- **中间件支持**：日志、恢复、CORS等中间件
- **Docker支持**：提供完整的Docker和Docker Compose配置
- **健康检查**：内置应用健康检查端点

## 📁 项目结构

```
.
├── cmd/server/           # 应用入口
├── internal/             # 内部包
│   ├── config/          # 配置管理
│   ├── handler/         # HTTP处理器
│   ├── middleware/      # 中间件
│   ├── model/          # 数据模型
│   ├── repository/     # 数据访问层
│   ├── router/         # 路由配置
│   └── service/        # 业务逻辑层
├── pkg/                 # 可复用包
│   ├── database/       # 数据库连接
│   ├── logger/         # 日志工具
│   └── response/       # 响应工具
├── configs/            # 配置文件
├── Dockerfile          # Docker镜像构建
├── docker-compose.yml  # Docker编排
└── Makefile           # 构建脚本
```

## 🔧 环境要求

- Go 1.21+
- MySQL 8.0+
- Docker & Docker Compose (可选)

## 🚀 快速开始

### 1. 克隆项目

```bash
git clone <repository-url>
cd go-api-scaffold
```

### 2. 安装依赖

```bash
go mod tidy
```

### 3. 配置数据库

编辑 `configs/app.yaml` 文件，更新数据库连接信息：

```yaml
database:
  driver: "mysql"
  host: "127.0.0.1"
  port: 3306
  username: "root"
  password: "your-password"
  database: "go_admin"
```

### 4. 运行应用

```bash
# 直接运行
go run cmd/server/main.go

# 或使用Makefile
make run
```

### 5. 使用Docker运行

```bash
# 使用Docker Compose（推荐）
make compose-up

# 或手动构建运行
make docker-build
make docker-run
```

## 📚 API文档

### 健康检查

- `GET /api/v1/health/check` - 健康检查
- `GET /api/v1/health/ready` - 就绪检查
- `GET /api/v1/health/live` - 存活检查

### 用户管理

- `POST /api/v1/users` - 创建用户
- `GET /api/v1/users/:id` - 获取用户信息
- `PUT /api/v1/users/:id` - 更新用户信息
- `DELETE /api/v1/users/:id` - 删除用户
- `GET /api/v1/users` - 获取用户列表

### 认证

- `POST /api/v1/login` - 用户登录

### 示例API

- `GET /api/v1/ping` - Ping测试
- `GET /api/v1/version` - 版本信息

## 🛠 开发命令

```bash
# 构建应用
make build

# 运行测试
make test

# 代码格式化
make fmt

# 代码检查
make vet

# 清理构建产物
make clean

# 下载依赖
make deps

# Docker相关
make docker-build    # 构建Docker镜像
make docker-run      # 运行Docker容器
make compose-up      # 启动所有服务
make compose-down    # 停止所有服务

# 健康检查
make health
```

## 🏗 架构设计

项目采用经典的分层架构：

- **Handler层**：处理HTTP请求和响应
- **Service层**：业务逻辑处理
- **Repository层**：数据访问抽象
- **Model层**：数据模型定义

## 🔒 安全特性

- 密码使用bcrypt加密
- 参数验证和数据校验
- 中间件支持（CORS、日志、恢复）
- 优雅关闭和错误处理

## ⚡ 性能优化

- 数据库连接池配置
- 精简的依赖管理
- 轻量级Docker镜像
- 高效的JSON处理

## 🚧 开发计划

- [ ] JWT认证实现
- [ ] API限流中间件
- [ ] 更多数据库操作示例
- [ ] 单元测试补充
- [ ] 性能监控集成

## 📄 许可证

MIT License

## 🤝 贡献

欢迎提交Issue和Pull Request来帮助改进这个项目！

---

## 🔄 优化历史

### v1.0.0 优化内容

1. **依赖精简**：移除了Swagger、PostgreSQL、SQLite等不必要的依赖
2. **数据库优化**：专门适配MySQL 8，简化配置
3. **功能简化**：移除了复杂的监控、日志收集等功能
4. **配置统一**：统一使用app.yaml配置文件
5. **Docker简化**：移除ELK、Prometheus等复杂服务
6. **Makefile优化**：保留核心构建命令，移除复杂选项

### 核心改进

- ✅ Go版本从1.24.5降级到1.21（稳定版本）
- ✅ 移除54个不必要的间接依赖
- ✅ 数据库配置从PostgreSQL改为MySQL 8
- ✅ 移除Redis依赖和配置
- ✅ 简化Docker镜像构建
- ✅ 优化项目启动流程
- ✅ 实际连接数据库并自动迁移表结构