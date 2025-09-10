# 部署指南

本文档详细介绍了 Go API Scaffold 项目的各种部署方式和最佳实践。

## 目录

- [环境准备](#环境准备)
- [本地部署](#本地部署)
- [Docker 部署](#docker-部署)
- [Docker Compose 部署](#docker-compose-部署)
- [云平台部署](#云平台部署)
- [生产环境配置](#生产环境配置)
- [监控和日志](#监控和日志)
- [故障排除](#故障排除)

## 环境准备

### 系统要求

- **Go**: 1.19 或更高版本
- **数据库**: PostgreSQL 12+, MySQL 8.0+, 或 SQLite 3
- **内存**: 最少 512MB RAM
- **存储**: 最少 1GB 可用空间

### 依赖检查

```bash
# 检查 Go 版本
go version

# 检查数据库连接（以 PostgreSQL 为例）
psql -h localhost -U postgres -d postgres -c "SELECT version();"
```

## 本地部署

### 1. 克隆和构建

```bash
# 克隆项目
git clone <repository-url>
cd go-api-scaffold

# 安装依赖
go mod tidy

# 构建应用
go build -o bin/server cmd/server/main.go
```

### 2. 配置环境

```bash
# 复制环境变量模板
cp .env.example .env

# 编辑配置文件
vim .env
```

### 3. 数据库准备

```sql
-- PostgreSQL
CREATE DATABASE scaffold_db;
CREATE USER scaffold_user WITH PASSWORD 'your_password';
GRANT ALL PRIVILEGES ON DATABASE scaffold_db TO scaffold_user;

-- MySQL
CREATE DATABASE scaffold_db;
CREATE USER 'scaffold_user'@'%' IDENTIFIED BY 'your_password';
GRANT ALL PRIVILEGES ON scaffold_db.* TO 'scaffold_user'@'%';
FLUSH PRIVILEGES;
```

### 4. 启动应用

```bash
# 直接运行
./bin/server

# 或者使用 go run
go run cmd/server/main.go

# 后台运行
nohup ./bin/server > app.log 2>&1 &
```

## Docker 部署

### 1. 创建 Dockerfile

```dockerfile
# 多阶段构建
FROM golang:1.21-alpine AS builder

# 设置工作目录
WORKDIR /app

# 安装必要的包
RUN apk add --no-cache git ca-certificates tzdata

# 复制 go mod 文件
COPY go.mod go.sum ./

# 下载依赖
RUN go mod download

# 复制源代码
COPY . .

# 构建应用
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main cmd/server/main.go

# 运行阶段
FROM alpine:latest

# 安装 ca-certificates 和 tzdata
RUN apk --no-cache add ca-certificates tzdata

# 创建非 root 用户
RUN addgroup -g 1001 -S appgroup && \
    adduser -u 1001 -S appuser -G appgroup

# 设置工作目录
WORKDIR /app

# 从构建阶段复制文件
COPY --from=builder /app/main .
COPY --from=builder /app/configs ./configs

# 更改文件所有者
RUN chown -R appuser:appgroup /app

# 切换到非 root 用户
USER appuser

# 暴露端口
EXPOSE 8080

# 健康检查
HEALTHCHECK --interval=30s --timeout=3s --start-period=5s --retries=3 \
    CMD wget --no-verbose --tries=1 --spider http://localhost:8080/health || exit 1

# 启动命令
CMD ["./main"]
```

### 2. 构建和运行

```bash
# 构建镜像
docker build -t go-api-scaffold:latest .

# 运行容器
docker run -d \
  --name go-api-scaffold \
  -p 8080:8080 \
  -e DB_HOST=host.docker.internal \
  -e DB_PORT=5432 \
  -e DB_NAME=scaffold_db \
  -e DB_USER=scaffold_user \
  -e DB_PASSWORD=your_password \
  go-api-scaffold:latest

# 查看日志
docker logs -f go-api-scaffold

# 进入容器
docker exec -it go-api-scaffold sh
```

### 3. Docker 优化

```dockerfile
# 使用 distroless 镜像（更安全）
FROM gcr.io/distroless/static:nonroot

WORKDIR /

COPY --from=builder /app/main /main
COPY --from=builder /app/configs /configs

USER nonroot:nonroot

EXPOSE 8080

ENTRYPOINT ["/main"]
```

## Docker Compose 部署

### 1. 创建 docker-compose.yml

```yaml
version: '3.8'

services:
  app:
    build:
      context: .
      dockerfile: Dockerfile
    container_name: go-api-scaffold
    ports:
      - "8080:8080"
    environment:
      - APP_ENV=production
      - DB_HOST=postgres
      - DB_PORT=5432
      - DB_NAME=scaffold_db
      - DB_USER=scaffold_user
      - DB_PASSWORD=your_password
      - LOG_LEVEL=info
    depends_on:
      postgres:
        condition: service_healthy
    restart: unless-stopped
    networks:
      - app-network
    volumes:
      - ./logs:/app/logs

  postgres:
    image: postgres:15-alpine
    container_name: postgres-db
    environment:
      - POSTGRES_DB=scaffold_db
      - POSTGRES_USER=scaffold_user
      - POSTGRES_PASSWORD=your_password
    ports:
      - "5432:5432"
    volumes:
      - postgres_data:/var/lib/postgresql/data
      - ./init.sql:/docker-entrypoint-initdb.d/init.sql
    healthcheck:
      test: ["CMD-SHELL", "pg_isready -U scaffold_user -d scaffold_db"]
      interval: 10s
      timeout: 5s
      retries: 5
    restart: unless-stopped
    networks:
      - app-network

  redis:
    image: redis:7-alpine
    container_name: redis-cache
    ports:
      - "6379:6379"
    volumes:
      - redis_data:/data
    healthcheck:
      test: ["CMD", "redis-cli", "ping"]
      interval: 10s
      timeout: 3s
      retries: 3
    restart: unless-stopped
    networks:
      - app-network

  nginx:
    image: nginx:alpine
    container_name: nginx-proxy
    ports:
      - "80:80"
      - "443:443"
    volumes:
      - ./nginx.conf:/etc/nginx/nginx.conf
      - ./ssl:/etc/nginx/ssl
    depends_on:
      - app
    restart: unless-stopped
    networks:
      - app-network

volumes:
  postgres_data:
  redis_data:

networks:
  app-network:
    driver: bridge
```

### 2. Nginx 配置

```nginx
events {
    worker_connections 1024;
}

http {
    upstream app {
        server app:8080;
    }

    server {
        listen 80;
        server_name localhost;

        location / {
            proxy_pass http://app;
            proxy_set_header Host $host;
            proxy_set_header X-Real-IP $remote_addr;
            proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
            proxy_set_header X-Forwarded-Proto $scheme;
        }

        location /health {
            proxy_pass http://app/health;
            access_log off;
        }
    }
}
```

### 3. 启动服务

```bash
# 启动所有服务
docker-compose up -d

# 查看服务状态
docker-compose ps

# 查看日志
docker-compose logs -f app

# 停止服务
docker-compose down

# 重新构建并启动
docker-compose up -d --build
```

## 云平台部署

### AWS ECS 部署

1. **创建 ECS 任务定义**

```json
{
  "family": "go-api-scaffold",
  "networkMode": "awsvpc",
  "requiresCompatibilities": ["FARGATE"],
  "cpu": "256",
  "memory": "512",
  "executionRoleArn": "arn:aws:iam::account:role/ecsTaskExecutionRole",
  "containerDefinitions": [
    {
      "name": "go-api-scaffold",
      "image": "your-account.dkr.ecr.region.amazonaws.com/go-api-scaffold:latest",
      "portMappings": [
        {
          "containerPort": 8080,
          "protocol": "tcp"
        }
      ],
      "environment": [
        {
          "name": "APP_ENV",
          "value": "production"
        }
      ],
      "secrets": [
        {
          "name": "DB_PASSWORD",
          "valueFrom": "arn:aws:secretsmanager:region:account:secret:db-password"
        }
      ],
      "logConfiguration": {
        "logDriver": "awslogs",
        "options": {
          "awslogs-group": "/ecs/go-api-scaffold",
          "awslogs-region": "us-west-2",
          "awslogs-stream-prefix": "ecs"
        }
      }
    }
  ]
}
```

### Google Cloud Run 部署

```bash
# 构建并推送到 Google Container Registry
gcloud builds submit --tag gcr.io/PROJECT_ID/go-api-scaffold

# 部署到 Cloud Run
gcloud run deploy go-api-scaffold \
  --image gcr.io/PROJECT_ID/go-api-scaffold \
  --platform managed \
  --region us-central1 \
  --allow-unauthenticated \
  --set-env-vars APP_ENV=production
```

### Kubernetes 部署

```yaml
# deployment.yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: go-api-scaffold
  labels:
    app: go-api-scaffold
spec:
  replicas: 3
  selector:
    matchLabels:
      app: go-api-scaffold
  template:
    metadata:
      labels:
        app: go-api-scaffold
    spec:
      containers:
      - name: go-api-scaffold
        image: go-api-scaffold:latest
        ports:
        - containerPort: 8080
        env:
        - name: APP_ENV
          value: "production"
        - name: DB_HOST
          valueFrom:
            secretKeyRef:
              name: db-secret
              key: host
        resources:
          requests:
            memory: "256Mi"
            cpu: "250m"
          limits:
            memory: "512Mi"
            cpu: "500m"
        livenessProbe:
          httpGet:
            path: /health/live
            port: 8080
          initialDelaySeconds: 30
          periodSeconds: 10
        readinessProbe:
          httpGet:
            path: /health/ready
            port: 8080
          initialDelaySeconds: 5
          periodSeconds: 5

---
apiVersion: v1
kind: Service
metadata:
  name: go-api-scaffold-service
spec:
  selector:
    app: go-api-scaffold
  ports:
    - protocol: TCP
      port: 80
      targetPort: 8080
  type: LoadBalancer
```

## 生产环境配置

### 1. 环境变量配置

```bash
# 生产环境 .env
APP_ENV=production
APP_NAME=go-api-scaffold
APP_VERSION=1.0.0

# 服务器配置
SERVER_HOST=0.0.0.0
SERVER_PORT=8080
SERVER_READ_TIMEOUT=30
SERVER_WRITE_TIMEOUT=30

# 数据库配置
DB_DRIVER=postgres
DB_HOST=your-db-host
DB_PORT=5432
DB_NAME=scaffold_db
DB_USER=scaffold_user
DB_PASSWORD=your-secure-password
DB_MAX_OPEN_CONNS=100
DB_MAX_IDLE_CONNS=10
DB_CONN_MAX_LIFETIME=3600

# 日志配置
LOG_LEVEL=info
LOG_FORMAT=json
LOG_OUTPUT=stdout

# Redis 配置（如果使用）
REDIS_HOST=your-redis-host
REDIS_PORT=6379
REDIS_PASSWORD=your-redis-password
```

### 2. 安全配置

```yaml
# configs/config.yaml
app:
  env: production
  debug: false

security:
  cors:
    allowed_origins:
      - "https://yourdomain.com"
      - "https://api.yourdomain.com"
    allowed_methods:
      - "GET"
      - "POST"
      - "PUT"
      - "DELETE"
    allowed_headers:
      - "Content-Type"
      - "Authorization"
    max_age: 86400

  rate_limit:
    enabled: true
    requests_per_minute: 100
    burst: 200
```

### 3. 系统服务配置

```ini
# /etc/systemd/system/go-api-scaffold.service
[Unit]
Description=Go API Scaffold
After=network.target

[Service]
Type=simple
User=appuser
Group=appgroup
WorkingDirectory=/opt/go-api-scaffold
ExecStart=/opt/go-api-scaffold/bin/server
Restart=always
RestartSec=5
EnvironmentFile=/opt/go-api-scaffold/.env

# 安全设置
NoNewPrivileges=true
PrivateTmp=true
ProtectSystem=strict
ProtectHome=true
ReadWritePaths=/opt/go-api-scaffold/logs

[Install]
WantedBy=multi-user.target
```

```bash
# 启用和启动服务
sudo systemctl enable go-api-scaffold
sudo systemctl start go-api-scaffold
sudo systemctl status go-api-scaffold
```

## 监控和日志

### 1. 日志配置

```yaml
# 使用 ELK Stack
version: '3.8'
services:
  elasticsearch:
    image: docker.elastic.co/elasticsearch/elasticsearch:8.5.0
    environment:
      - discovery.type=single-node
      - "ES_JAVA_OPTS=-Xms512m -Xmx512m"
    ports:
      - "9200:9200"

  logstash:
    image: docker.elastic.co/logstash/logstash:8.5.0
    volumes:
      - ./logstash.conf:/usr/share/logstash/pipeline/logstash.conf
    ports:
      - "5044:5044"
    depends_on:
      - elasticsearch

  kibana:
    image: docker.elastic.co/kibana/kibana:8.5.0
    ports:
      - "5601:5601"
    depends_on:
      - elasticsearch
```

### 2. 监控配置

```yaml
# Prometheus + Grafana
prometheus:
  image: prom/prometheus
  ports:
    - "9090:9090"
  volumes:
    - ./prometheus.yml:/etc/prometheus/prometheus.yml

grafana:
  image: grafana/grafana
  ports:
    - "3000:3000"
  environment:
    - GF_SECURITY_ADMIN_PASSWORD=admin
  volumes:
    - grafana_data:/var/lib/grafana
```

## 故障排除

### 常见问题

1. **数据库连接失败**
```bash
# 检查数据库连接
telnet db-host 5432

# 检查数据库日志
docker logs postgres-container
```

2. **内存不足**
```bash
# 检查内存使用
free -h
top -p $(pgrep server)

# 调整容器内存限制
docker run -m 1g go-api-scaffold
```

3. **端口冲突**
```bash
# 检查端口占用
netstat -tulpn | grep 8080
lsof -i :8080
```

### 调试命令

```bash
# 检查应用状态
curl http://localhost:8080/health

# 查看应用日志
tail -f /var/log/go-api-scaffold/app.log

# 检查系统资源
htop
iostat -x 1

# 数据库连接测试
psql -h localhost -U scaffold_user -d scaffold_db -c "SELECT 1;"
```

### 性能优化

1. **数据库优化**
```sql
-- 添加索引
CREATE INDEX idx_users_username ON users(username);
CREATE INDEX idx_users_email ON users(email);

-- 分析查询性能
EXPLAIN ANALYZE SELECT * FROM users WHERE username = 'john';
```

2. **应用优化**
```go
// 连接池配置
db.SetMaxOpenConns(100)
db.SetMaxIdleConns(10)
db.SetConnMaxLifetime(time.Hour)

// 启用 gzip 压缩
r.Use(gin.Recovery())
r.Use(gin.Logger())
r.Use(gzip.Gzip(gzip.DefaultCompression))
```

3. **缓存配置**
```go
// Redis 缓存
rdb := redis.NewClient(&redis.Options{
    Addr:         "localhost:6379",
    PoolSize:     10,
    MinIdleConns: 5,
})
```

---

**注意事项**:
1. 生产环境部署前请进行充分测试
2. 定期备份数据库和配置文件
3. 监控应用性能和资源使用情况
4. 及时更新依赖包和安全补丁
5. 实施适当的安全措施和访问控制