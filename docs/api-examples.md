# API 使用示例

本文档提供了 Go API Scaffold 项目中各个 API 端点的详细使用示例。

## 基础信息

- **Base URL**: `http://localhost:8080`
- **Content-Type**: `application/json`
- **响应格式**: 统一的 JSON 响应格式

## 响应格式说明

所有 API 响应都遵循统一的格式：

```json
{
  "code": 200,
  "message": "success",
  "data": {},
  "timestamp": "2024-01-15T10:30:00Z"
}
```

- `code`: 响应状态码
- `message`: 响应消息
- `data`: 响应数据（可选）
- `timestamp`: 响应时间戳

## 健康检查 API

### 1. 基础健康检查

**请求**:
```bash
curl -X GET http://localhost:8080/health
```

**响应**:
```json
{
  "code": 200,
  "message": "success",
  "data": {
    "status": "healthy",
    "timestamp": "2024-01-15T10:30:00Z",
    "version": "1.0.0",
    "uptime": "2h30m15s"
  },
  "timestamp": "2024-01-15T10:30:00Z"
}
```

### 2. 就绪检查

**请求**:
```bash
curl -X GET http://localhost:8080/health/ready
```

**响应**:
```json
{
  "code": 200,
  "message": "success",
  "data": {
    "status": "ready"
  },
  "timestamp": "2024-01-15T10:30:00Z"
}
```

### 3. 存活检查

**请求**:
```bash
curl -X GET http://localhost:8080/health/live
```

**响应**:
```json
{
  "code": 200,
  "message": "success",
  "data": {
    "status": "alive"
  },
  "timestamp": "2024-01-15T10:30:00Z"
}
```

## 用户管理 API

### 1. 创建用户

**请求**:
```bash
curl -X POST http://localhost:8080/api/v1/users \
  -H "Content-Type: application/json" \
  -d '{
    "username": "john_doe",
    "email": "john@example.com",
    "password": "password123",
    "full_name": "John Doe"
  }'
```

**响应**:
```json
{
  "code": 201,
  "message": "User created successfully",
  "data": {
    "id": 1,
    "username": "john_doe",
    "email": "john@example.com",
    "full_name": "John Doe",
    "is_active": true,
    "created_at": "2024-01-15T10:30:00Z",
    "updated_at": "2024-01-15T10:30:00Z"
  },
  "timestamp": "2024-01-15T10:30:00Z"
}
```

**错误响应示例**:
```json
{
  "code": 400,
  "message": "Username already exists",
  "timestamp": "2024-01-15T10:30:00Z"
}
```

### 2. 获取用户详情

**请求**:
```bash
curl -X GET http://localhost:8080/api/v1/users/1
```

**响应**:
```json
{
  "code": 200,
  "message": "success",
  "data": {
    "id": 1,
    "username": "john_doe",
    "email": "john@example.com",
    "full_name": "John Doe",
    "is_active": true,
    "last_login_at": "2024-01-15T09:00:00Z",
    "created_at": "2024-01-15T08:00:00Z",
    "updated_at": "2024-01-15T09:00:00Z"
  },
  "timestamp": "2024-01-15T10:30:00Z"
}
```

**用户不存在响应**:
```json
{
  "code": 404,
  "message": "User not found",
  "timestamp": "2024-01-15T10:30:00Z"
}
```

### 3. 更新用户

**请求**:
```bash
curl -X PUT http://localhost:8080/api/v1/users/1 \
  -H "Content-Type: application/json" \
  -d '{
    "full_name": "John Smith",
    "email": "john.smith@example.com"
  }'
```

**响应**:
```json
{
  "code": 200,
  "message": "User updated successfully",
  "data": {
    "id": 1,
    "username": "john_doe",
    "email": "john.smith@example.com",
    "full_name": "John Smith",
    "is_active": true,
    "last_login_at": "2024-01-15T09:00:00Z",
    "created_at": "2024-01-15T08:00:00Z",
    "updated_at": "2024-01-15T10:30:00Z"
  },
  "timestamp": "2024-01-15T10:30:00Z"
}
```

### 4. 删除用户

**请求**:
```bash
curl -X DELETE http://localhost:8080/api/v1/users/1
```

**响应**:
```json
{
  "code": 200,
  "message": "User deleted successfully",
  "timestamp": "2024-01-15T10:30:00Z"
}
```

### 5. 获取用户列表

**请求**:
```bash
# 基础请求
curl -X GET http://localhost:8080/api/v1/users

# 带分页参数
curl -X GET "http://localhost:8080/api/v1/users?page=1&limit=10"

# 带搜索参数
curl -X GET "http://localhost:8080/api/v1/users?search=john&page=1&limit=5"
```

**响应**:
```json
{
  "code": 200,
  "message": "success",
  "data": {
    "users": [
      {
        "id": 1,
        "username": "john_doe",
        "email": "john@example.com",
        "full_name": "John Doe",
        "is_active": true,
        "last_login_at": "2024-01-15T09:00:00Z",
        "created_at": "2024-01-15T08:00:00Z",
        "updated_at": "2024-01-15T09:00:00Z"
      },
      {
        "id": 2,
        "username": "jane_doe",
        "email": "jane@example.com",
        "full_name": "Jane Doe",
        "is_active": true,
        "last_login_at": null,
        "created_at": "2024-01-15T08:30:00Z",
        "updated_at": "2024-01-15T08:30:00Z"
      }
    ],
    "pagination": {
      "page": 1,
      "limit": 10,
      "total": 2,
      "total_pages": 1
    }
  },
  "timestamp": "2024-01-15T10:30:00Z"
}
```

## 认证 API

### 用户登录

**请求**:
```bash
curl -X POST http://localhost:8080/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "username": "john_doe",
    "password": "password123"
  }'
```

**成功响应**:
```json
{
  "code": 200,
  "message": "Login successful",
  "data": {
    "user": {
      "id": 1,
      "username": "john_doe",
      "email": "john@example.com",
      "full_name": "John Doe",
      "is_active": true,
      "last_login_at": "2024-01-15T10:30:00Z"
    }
  },
  "timestamp": "2024-01-15T10:30:00Z"
}
```

**登录失败响应**:
```json
{
  "code": 401,
  "message": "Invalid username or password",
  "timestamp": "2024-01-15T10:30:00Z"
}
```

**使用邮箱登录**:
```bash
curl -X POST http://localhost:8080/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "username": "john@example.com",
    "password": "password123"
  }'
```

## 错误处理

### 常见错误码

- `400` - 请求参数错误
- `401` - 未授权
- `404` - 资源不存在
- `409` - 资源冲突（如用户名已存在）
- `422` - 数据验证失败
- `500` - 服务器内部错误

### 验证错误示例

**请求**:
```bash
curl -X POST http://localhost:8080/api/v1/users \
  -H "Content-Type: application/json" \
  -d '{
    "username": "",
    "email": "invalid-email",
    "password": "123"
  }'
```

**响应**:
```json
{
  "code": 422,
  "message": "Validation failed",
  "data": {
    "errors": [
      {
        "field": "username",
        "message": "Username is required"
      },
      {
        "field": "email",
        "message": "Email format is invalid"
      },
      {
        "field": "password",
        "message": "Password must be at least 6 characters"
      }
    ]
  },
  "timestamp": "2024-01-15T10:30:00Z"
}
```

## 使用 Postman

### 环境变量设置

在 Postman 中创建环境变量：

- `base_url`: `http://localhost:8080`
- `api_version`: `v1`

### 集合示例

可以创建以下请求集合：

1. **Health Checks**
   - GET `{{base_url}}/health`
   - GET `{{base_url}}/health/ready`
   - GET `{{base_url}}/health/live`

2. **User Management**
   - POST `{{base_url}}/api/{{api_version}}/users`
   - GET `{{base_url}}/api/{{api_version}}/users/{{user_id}}`
   - PUT `{{base_url}}/api/{{api_version}}/users/{{user_id}}`
   - DELETE `{{base_url}}/api/{{api_version}}/users/{{user_id}}`
   - GET `{{base_url}}/api/{{api_version}}/users`

3. **Authentication**
   - POST `{{base_url}}/api/{{api_version}}/auth/login`

## 测试建议

1. **顺序测试**: 先创建用户，再进行其他操作
2. **边界测试**: 测试空值、超长字符串等边界情况
3. **错误测试**: 故意发送错误数据，验证错误处理
4. **并发测试**: 使用工具进行并发请求测试

## 性能测试

使用 Apache Bench 进行简单的性能测试：

```bash
# 健康检查端点性能测试
ab -n 1000 -c 10 http://localhost:8080/health

# 用户列表端点性能测试
ab -n 500 -c 5 http://localhost:8080/api/v1/users
```

---

**注意**: 在生产环境中使用时，请确保：
1. 使用 HTTPS
2. 实现适当的认证和授权机制
3. 添加请求限流
4. 配置适当的 CORS 策略