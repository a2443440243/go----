# 开发指南

本文档提供了 Go API Scaffold 项目的详细开发指南，帮助开发者快速上手和贡献代码。

## 目录

- [开发环境设置](#开发环境设置)
- [项目结构](#项目结构)
- [开发工作流](#开发工作流)
- [代码规范](#代码规范)
- [测试指南](#测试指南)
- [调试技巧](#调试技巧)
- [性能优化](#性能优化)
- [常见问题](#常见问题)

## 开发环境设置

### 必需工具

1. **Go 1.21+**
   ```bash
   # 检查 Go 版本
   go version
   ```

2. **Git**
   ```bash
   git --version
   ```

3. **Make**（可选，但推荐）
   ```bash
   make --version
   ```

4. **Docker**（用于容器化开发）
   ```bash
   docker --version
   docker-compose --version
   ```

### 推荐工具

1. **Air**（热重载）
   ```bash
   go install github.com/cosmtrek/air@latest
   ```

2. **golangci-lint**（代码检查）
   ```bash
   # Linux/macOS
   curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(go env GOPATH)/bin v1.54.2
   
   # Windows
   go install github.com/golangci/golangci-lint/cmd/golangci-lint@v1.54.2
   ```

3. **gosec**（安全检查）
   ```bash
   go install github.com/securecodewarrior/gosec/v2/cmd/gosec@latest
   ```

### IDE 配置

#### VS Code

推荐扩展：
- Go (Google)
- Go Test Explorer
- REST Client
- Docker
- GitLens

配置文件 `.vscode/settings.json`：
```json
{
  "go.useLanguageServer": true,
  "go.lintTool": "golangci-lint",
  "go.lintFlags": ["--fast"],
  "go.formatTool": "goimports",
  "go.testFlags": ["-v", "-race"],
  "go.coverOnSave": true,
  "go.coverageDecorator": {
    "type": "gutter",
    "coveredHighlightColor": "rgba(64,128,128,0.5)",
    "uncoveredHighlightColor": "rgba(128,64,64,0.25)"
  }
}
```

#### GoLand/IntelliJ IDEA

1. 安装 Go 插件
2. 配置 Go SDK
3. 启用 Go Modules
4. 配置代码格式化工具

## 项目结构

```
go-api-scaffold/
├── cmd/                    # 应用程序入口
│   └── server/
│       └── main.go
├── internal/               # 私有应用程序代码
│   ├── config/            # 配置管理
│   ├── handler/           # HTTP 处理器
│   ├── middleware/        # 中间件
│   ├── model/             # 数据模型
│   ├── repository/        # 数据访问层
│   ├── service/           # 业务逻辑层
│   └── utils/             # 工具函数
├── pkg/                   # 可被外部应用程序使用的库代码
│   ├── database/          # 数据库连接
│   ├── logger/            # 日志工具
│   └── response/          # 响应工具
├── configs/               # 配置文件
├── docs/                  # 文档
├── scripts/               # 脚本文件
├── tests/                 # 测试文件
└── deployments/           # 部署配置
```

### 架构原则

1. **分层架构**：Handler → Service → Repository → Model
2. **依赖注入**：使用接口和依赖注入
3. **单一职责**：每个包和文件都有明确的职责
4. **接口隔离**：定义小而专注的接口

## 开发工作流

### 1. 克隆项目

```bash
git clone <repository-url>
cd go-api-scaffold
```

### 2. 安装依赖

```bash
# 使用 Make
make deps

# 或者直接使用 Go
go mod download
go mod tidy
```

### 3. 配置环境

```bash
# 复制配置文件
cp configs/app.yaml.example configs/app.yaml

# 编辑配置
vim configs/app.yaml
```

### 4. 启动开发服务器

```bash
# 使用 Make（推荐）
make dev

# 或者使用 Air
air

# 或者直接运行
go run cmd/server/main.go
```

### 5. 运行测试

```bash
# 运行所有测试
make test

# 运行测试并生成覆盖率报告
make test-cover

# 运行特定包的测试
go test -v ./internal/service/...
```

### 6. 代码检查

```bash
# 格式化代码
make fmt

# 运行 linter
make lint

# 运行安全检查
make security
```

## 代码规范

### Go 代码规范

1. **命名规范**
   - 包名：小写，简短，有意义
   - 变量名：驼峰命名法
   - 常量名：大写字母和下划线
   - 接口名：以 "er" 结尾（如 Reader, Writer）

2. **注释规范**
   ```go
   // Package service 提供业务逻辑处理
   package service
   
   // UserService 定义用户服务接口
   type UserService interface {
       // Create 创建新用户
       Create(ctx context.Context, req *CreateUserRequest) (*User, error)
   }
   ```

3. **错误处理**
   ```go
   // 好的做法
   if err != nil {
       return nil, fmt.Errorf("failed to create user: %w", err)
   }
   
   // 避免
   if err != nil {
       panic(err)
   }
   ```

4. **上下文使用**
   ```go
   // 所有公共方法都应该接受 context.Context
   func (s *userService) GetByID(ctx context.Context, id uint) (*User, error) {
       // 实现
   }
   ```

### 项目特定规范

1. **文件组织**
   - 每个包一个目录
   - 接口定义在单独的文件中
   - 测试文件以 `_test.go` 结尾

2. **依赖注入**
   ```go
   // 定义接口
   type UserRepository interface {
       Create(ctx context.Context, user *User) error
   }
   
   // 实现结构体
   type userService struct {
       repo UserRepository
   }
   
   // 构造函数
   func NewUserService(repo UserRepository) UserService {
       return &userService{repo: repo}
   }
   ```

3. **配置管理**
   ```go
   // 使用 Viper 读取配置
   type Config struct {
       Server ServerConfig `mapstructure:"server"`
       DB     DBConfig     `mapstructure:"database"`
   }
   ```

## 测试指南

### 测试类型

1. **单元测试**
   ```go
   func TestUserService_Create(t *testing.T) {
       // 准备
       mockRepo := &MockUserRepository{}
       service := NewUserService(mockRepo)
       
       // 执行
       user, err := service.Create(context.Background(), &CreateUserRequest{
           Username: "testuser",
           Email:    "test@example.com",
       })
       
       // 断言
       assert.NoError(t, err)
       assert.NotNil(t, user)
       assert.Equal(t, "testuser", user.Username)
   }
   ```

2. **集成测试**
   ```go
   func TestUserAPI_Integration(t *testing.T) {
       // 设置测试数据库
       db := setupTestDB(t)
       defer cleanupTestDB(t, db)
       
       // 创建测试服务器
       server := setupTestServer(t, db)
       defer server.Close()
       
       // 执行 HTTP 请求测试
       resp, err := http.Post(server.URL+"/api/users", "application/json", 
           strings.NewReader(`{"username":"test","email":"test@example.com"}`))
       
       assert.NoError(t, err)
       assert.Equal(t, http.StatusCreated, resp.StatusCode)
   }
   ```

3. **基准测试**
   ```go
   func BenchmarkUserService_Create(b *testing.B) {
       service := setupBenchmarkService()
       req := &CreateUserRequest{
           Username: "benchuser",
           Email:    "bench@example.com",
       }
       
       b.ResetTimer()
       for i := 0; i < b.N; i++ {
           _, _ = service.Create(context.Background(), req)
       }
   }
   ```

### 测试工具

1. **Testify**（断言和模拟）
   ```go
   import (
       "github.com/stretchr/testify/assert"
       "github.com/stretchr/testify/mock"
       "github.com/stretchr/testify/suite"
   )
   ```

2. **GoMock**（生成模拟对象）
   ```bash
   go install github.com/golang/mock/mockgen@latest
   mockgen -source=internal/repository/user.go -destination=mocks/user_repository.go
   ```

3. **Testcontainers**（集成测试）
   ```go
   import "github.com/testcontainers/testcontainers-go"
   ```

## 调试技巧

### 1. 日志调试

```go
// 使用结构化日志
logger.WithFields(logrus.Fields{
    "user_id": userID,
    "action":  "create_user",
}).Info("Creating new user")

// 调试级别日志
logger.Debug("Processing request", "request", req)
```

### 2. Delve 调试器

```bash
# 安装 Delve
go install github.com/go-delve/delve/cmd/dlv@latest

# 调试应用
dlv debug cmd/server/main.go

# 在 VS Code 中调试
# 使用 F5 或配置 launch.json
```

### 3. 性能分析

```go
// 添加 pprof 端点
import _ "net/http/pprof"

// 在单独的 goroutine 中启动 pprof 服务器
go func() {
    log.Println(http.ListenAndServe("localhost:6060", nil))
}()
```

```bash
# 分析 CPU 使用
go tool pprof http://localhost:6060/debug/pprof/profile

# 分析内存使用
go tool pprof http://localhost:6060/debug/pprof/heap
```

## 性能优化

### 1. 数据库优化

```go
// 使用连接池
db.SetMaxOpenConns(25)
db.SetMaxIdleConns(5)
db.SetConnMaxLifetime(5 * time.Minute)

// 使用预编译语句
stmt, err := db.Prepare("SELECT * FROM users WHERE id = ?")
defer stmt.Close()

// 批量操作
db.CreateInBatches(&users, 100)
```

### 2. 缓存策略

```go
// Redis 缓存
func (s *userService) GetByID(ctx context.Context, id uint) (*User, error) {
    // 先从缓存获取
    if user := s.cache.Get(fmt.Sprintf("user:%d", id)); user != nil {
        return user.(*User), nil
    }
    
    // 从数据库获取
    user, err := s.repo.GetByID(ctx, id)
    if err != nil {
        return nil, err
    }
    
    // 写入缓存
    s.cache.Set(fmt.Sprintf("user:%d", id), user, 10*time.Minute)
    return user, nil
}
```

### 3. 并发处理

```go
// 使用 worker pool
type WorkerPool struct {
    workers int
    jobs    chan Job
    results chan Result
}

func (wp *WorkerPool) Start() {
    for i := 0; i < wp.workers; i++ {
        go wp.worker()
    }
}

func (wp *WorkerPool) worker() {
    for job := range wp.jobs {
        result := job.Process()
        wp.results <- result
    }
}
```

## 常见问题

### Q: 如何添加新的 API 端点？

A: 按照以下步骤：

1. 在 `internal/model` 中定义数据模型
2. 在 `internal/repository` 中添加数据访问方法
3. 在 `internal/service` 中实现业务逻辑
4. 在 `internal/handler` 中添加 HTTP 处理器
5. 在路由中注册新端点
6. 编写测试

### Q: 如何处理数据库迁移？

A: 使用 GORM 的自动迁移功能：

```go
// 在应用启动时
db.AutoMigrate(&User{}, &Post{})

// 或者使用专门的迁移命令
func migrate() {
    db := database.GetDB()
    db.AutoMigrate(&User{})
}
```

### Q: 如何添加中间件？

A: 在 `internal/middleware` 中创建中间件：

```go
func AuthMiddleware() gin.HandlerFunc {
    return func(c *gin.Context) {
        token := c.GetHeader("Authorization")
        if token == "" {
            c.JSON(401, gin.H{"error": "Unauthorized"})
            c.Abort()
            return
        }
        c.Next()
    }
}

// 在路由中使用
r.Use(AuthMiddleware())
```

### Q: 如何配置不同环境？

A: 使用不同的配置文件：

```bash
# 开发环境
APP_ENV=development go run cmd/server/main.go

# 生产环境
APP_ENV=production go run cmd/server/main.go
```

配置文件：
- `configs/app.yaml` - 默认配置
- `configs/app.development.yaml` - 开发环境
- `configs/app.production.yaml` - 生产环境

## 贡献指南

1. Fork 项目
2. 创建功能分支：`git checkout -b feature/new-feature`
3. 提交更改：`git commit -am 'Add new feature'`
4. 推送分支：`git push origin feature/new-feature`
5. 创建 Pull Request

### 提交消息规范

```
type(scope): description

[optional body]

[optional footer]
```

类型：
- `feat`: 新功能
- `fix`: 修复 bug
- `docs`: 文档更新
- `style`: 代码格式化
- `refactor`: 重构
- `test`: 测试相关
- `chore`: 构建过程或辅助工具的变动

示例：
```
feat(user): add user registration endpoint

Implement user registration with email validation
and password hashing using bcrypt.

Closes #123
```

## 资源链接

- [Go 官方文档](https://golang.org/doc/)
- [Gin 框架文档](https://gin-gonic.com/docs/)
- [GORM 文档](https://gorm.io/docs/)
- [Viper 配置管理](https://github.com/spf13/viper)
- [Testify 测试框架](https://github.com/stretchr/testify)
- [Go 代码审查指南](https://github.com/golang/go/wiki/CodeReviewComments)