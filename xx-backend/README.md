# XX-Backend

基于 Go + Gin + Redis + MySQL + gRPC + 多线程的后端管理系统

## 技术栈

- **Web框架**: Gin
- **数据库**: MySQL + GORM
- **缓存**: Redis
- **认证**: JWT
- **微服务**: gRPC
- **并发**: Go协程 + 多线程
- **配置**: 环境变量

## 项目结构

```
xx-backend/
├── config/          # 配置管理
├── internal/        # 内部包
│   ├── handler/     # HTTP处理器
│   ├── middleware/  # 中间件
│   ├── model/       # 数据模型
│   └── service/     # 业务逻辑
├── pkg/             # 公共包
│   ├── database/    # 数据库连接
│   └── redis/       # Redis连接
├── go.mod           # Go模块文件
├── main.go          # 主程序入口
└── README.md        # 项目说明
```

## 功能特性

- ✅ 用户认证 (JWT + Redis)
- ✅ 用户管理 (CRUD)
- ✅ 角色管理
- ✅ 菜单管理
- ✅ 权限控制
- ✅ 多线程处理
- ✅ gRPC服务
- ✅ 优雅关闭
- ✅ 跨域支持
- ✅ 日志记录

## 安装运行

### 1. 环境要求

- Go 1.21+
- MySQL 8.0+
- Redis 6.0+

### 2. 安装依赖

```bash
go mod tidy
```

### 3. 配置环境变量

```bash
# 应用配置
export APP_MODE=debug
export APP_PORT=8080

# MySQL配置
export MYSQL_HOST=localhost
export MYSQL_PORT=3306
export MYSQL_USER=root
export MYSQL_PASSWORD=your_password
export MYSQL_DATABASE=xx_admin

# Redis配置
export REDIS_HOST=localhost
export REDIS_PORT=6379
export REDIS_PASSWORD=
export REDIS_DB=0
```

### 4. 创建数据库

```sql
CREATE DATABASE xx_admin CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;
```

### 5. 运行项目

```bash
go run main.go
```

## API接口

### 认证相关

- `POST /api/auth/login` - 用户登录
- `POST /api/auth/logout` - 用户登出
- `GET /api/auth/profile` - 获取用户资料

### 用户管理

- `GET /api/users` - 获取用户列表
- `GET /api/users/:id` - 获取用户详情
- `POST /api/users` - 创建用户
- `PUT /api/users/:id` - 更新用户
- `DELETE /api/users/:id` - 删除用户

### 角色管理

- `GET /api/roles` - 获取角色列表
- `POST /api/roles` - 创建角色
- `PUT /api/roles/:id` - 更新角色
- `DELETE /api/roles/:id` - 删除角色

### 菜单管理

- `GET /api/menus` - 获取菜单列表
- `POST /api/menus` - 创建菜单
- `PUT /api/menus/:id` - 更新菜单
- `DELETE /api/menus/:id` - 删除菜单

## 多线程特性

项目使用Go协程实现多线程处理：

1. **gRPC服务**: 在独立协程中运行
2. **批量处理**: 用户批量操作使用工作池模式
3. **并发安全**: 使用互斥锁保护共享资源

## 部署

### Docker部署

```dockerfile
FROM golang:1.21-alpine AS builder
WORKDIR /app
COPY . .
RUN go mod download
RUN go build -o main .

FROM alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /root/
COPY --from=builder /app/main .
CMD ["./main"]
```

### 生产环境

1. 设置 `APP_MODE=production`
2. 配置生产环境数据库
3. 使用反向代理 (Nginx)
4. 配置SSL证书

## 开发

### 添加新功能

1. 在 `model/` 中定义数据模型
2. 在 `service/` 中实现业务逻辑
3. 在 `handler/` 中处理HTTP请求
4. 在 `main.go` 中注册路由

### 测试

```bash
go test ./...
```

## 许可证

MIT License 