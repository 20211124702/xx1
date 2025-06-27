# XX-Admin 通用后台管理系统

一个基于 Vue 3 + Go 的现代化后台管理系统，包含完整的前后端解决方案和 Docker 容器化部署。

## 🚀 技术栈

### 前端 (xx-admin)
- **框架**: Vue 3 + TypeScript
- **构建工具**: Vite
- **UI组件库**: Element Plus
- **状态管理**: Pinia
- **路由**: Vue Router 4
- **HTTP客户端**: Axios
- **接口管理**: Apifox

### 后端 (xx-backend)
- **语言**: Go 1.21+
- **Web框架**: Gin
- **数据库**: MySQL + GORM
- **缓存**: Redis
- **认证**: JWT
- **微服务**: gRPC
- **并发**: Go协程 + 多线程

### 容器化
- **容器**: Docker
- **编排**: Docker Compose
- **Web服务器**: Nginx
- **热重载**: Air (开发环境)

## 📁 项目结构

```
xx1/
├── xx-admin/          # 前端项目
│   ├── src/
│   │   ├── api/       # API接口
│   │   ├── components/# 组件
│   │   ├── router/    # 路由配置
│   │   ├── store/     # 状态管理
│   │   ├── views/     # 页面
│   │   └── main.ts    # 入口文件
│   ├── Dockerfile     # 前端Docker配置
│   ├── nginx.conf     # Nginx配置
│   └── package.json
├── xx-backend/        # 后端项目
│   ├── config/        # 配置管理
│   ├── internal/      # 内部包
│   ├── pkg/          # 公共包
│   ├── Dockerfile     # 生产环境Docker配置
│   ├── Dockerfile.dev # 开发环境Docker配置
│   ├── .air.toml     # 热重载配置
│   ├── main.go       # 主程序
│   └── go.mod
├── mysql/             # 数据库初始化
│   └── init.sql      # 初始数据
├── docker-compose.yml # 生产环境编排
├── docker-compose.dev.yml # 开发环境编排
├── Makefile          # 便捷命令
└── README.md
```

## ✨ 功能特性

### 前端功能
- ✅ 现代化UI设计
- ✅ 响应式布局
- ✅ 路由守卫
- ✅ 权限控制
- ✅ 用户管理
- ✅ 角色管理
- ✅ 菜单管理
- ✅ 数据表格
- ✅ 登录认证

### 后端功能
- ✅ RESTful API
- ✅ JWT认证
- ✅ 用户管理CRUD
- ✅ 角色权限管理
- ✅ 菜单管理
- ✅ 数据库自动迁移
- ✅ Redis缓存
- ✅ gRPC微服务
- ✅ 多线程处理
- ✅ 优雅关闭
- ✅ 跨域支持
- ✅ 日志记录

### Docker功能
- ✅ 多阶段构建
- ✅ 生产环境优化
- ✅ 开发环境热重载
- ✅ 健康检查
- ✅ 数据持久化
- ✅ 网络隔离
- ✅ 一键部署

## 🛠️ 快速开始

### 方式一：Docker 部署（推荐）

#### 1. 环境要求
- Docker 20.10+
- Docker Compose 2.0+

#### 2. 一键启动
```bash
# 启动生产环境
make up

# 或使用docker-compose
docker-compose up -d
```

#### 3. 访问应用
- 前端: http://localhost
- 后端API: http://localhost:8080
- 数据库: localhost:3306
- Redis: localhost:6379

#### 4. 默认账号
- 用户名: `admin`
- 密码: `admin123`

### 方式二：开发环境

#### 1. 启动开发环境
```bash
# 启动开发环境（支持热重载）
make dev

# 或使用docker-compose
docker-compose -f docker-compose.dev.yml up -d
```

#### 2. 本地开发
```bash
# 前端开发
cd xx-admin
npm install
npm run dev

# 后端开发
cd xx-backend
go mod tidy
go run main.go
```

### 方式三：传统部署

#### 1. 环境要求
- Node.js 16+
- Go 1.21+
- MySQL 8.0+
- Redis 6.0+

#### 2. 启动前端
```bash
cd xx-admin
npm install
npm run dev
```

#### 3. 启动后端
```bash
cd xx-backend
go mod tidy
go run main.go
```

## 🐳 Docker 命令

### 基础命令
```bash
# 查看帮助
make help

# 构建镜像
make build

# 启动生产环境
make up

# 启动开发环境
make dev

# 停止服务
make down
make dev-down

# 查看日志
make logs

# 清理环境
make clean
```

### 高级命令
```bash
# 进入容器
make shell-backend
make shell-frontend
make shell-mysql

# 数据库备份/恢复
make backup
make restore

# 健康检查
make health
```

### 手动Docker命令
```bash
# 构建单个服务
docker-compose build frontend
docker-compose build backend

# 重启服务
docker-compose restart backend

# 查看服务状态
docker-compose ps

# 查看服务日志
docker-compose logs -f backend
```

## 📖 API文档

### 认证接口
- `POST /api/auth/login` - 用户登录
- `POST /api/auth/logout` - 用户登出
- `GET /api/auth/profile` - 获取用户资料

### 用户管理
- `GET /api/users` - 获取用户列表
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

## 🔧 开发指南

### 前端开发

1. **添加新页面**
   ```bash
   # 在 src/views/ 下创建新页面
   # 在 src/router/index.ts 中添加路由
   ```

2. **添加新API**
   ```bash
   # 在 src/api/ 下创建新接口文件
   # 在页面中导入并使用
   ```

3. **添加新组件**
   ```bash
   # 在 src/components/ 下创建组件
   # 在页面中导入并使用
   ```

### 后端开发

1. **添加新模型**
   ```bash
   # 在 internal/model/ 下创建模型
   # 在 main.go 中添加自动迁移
   ```

2. **添加新服务**
   ```bash
   # 在 internal/service/ 下创建服务
   # 在 internal/handler/ 下创建处理器
   # 在 main.go 中注册路由
   ```

### Docker开发

1. **修改Docker配置**
   ```bash
   # 修改 Dockerfile 或 docker-compose.yml
   # 重新构建镜像
   make build
   ```

2. **开发环境热重载**
   ```bash
   # 启动开发环境
   make dev
   # 修改代码后自动重载
   ```

## 🚀 部署

### 生产环境部署

```bash
# 1. 构建镜像
make build

# 2. 启动服务
make up

# 3. 检查状态
make health
```

### 自定义部署

```bash
# 修改环境变量
vim docker-compose.yml

# 重新部署
docker-compose down
docker-compose up -d
```

### 集群部署

```bash
# 使用 Docker Swarm
docker swarm init
docker stack deploy -c docker-compose.yml xx-admin

# 使用 Kubernetes
kubectl apply -f k8s/
```

## 🔒 安全配置

### 生产环境安全
1. 修改默认密码
2. 配置SSL证书
3. 设置防火墙
4. 定期备份数据
5. 监控服务状态

### 环境变量配置
```bash
# 数据库配置
MYSQL_ROOT_PASSWORD=your_secure_password
MYSQL_USER=your_user
MYSQL_PASSWORD=your_password

# Redis配置
REDIS_PASSWORD=your_redis_password

# JWT密钥
JWT_SECRET=your_jwt_secret
```

## 📝 开发计划

- [ ] 添加单元测试
- [ ] 集成WebSocket
- [ ] 添加文件上传
- [ ] 添加数据导出
- [ ] 添加系统监控
- [ ] 添加操作日志
- [ ] 添加数据备份
- [ ] 添加多租户支持
- [ ] Kubernetes部署
- [ ] CI/CD流水线

## 🤝 贡献

欢迎提交 Issue 和 Pull Request！

## �� 许可证

MIT License 