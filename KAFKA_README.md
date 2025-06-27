# Kafka 集成说明

本项目已成功集成 Apache Kafka，用于事件驱动架构和日志记录。

## 功能特性

### 1. 事件记录
- **用户事件**: 登录、登出、注册、更新、删除
- **系统日志**: 错误日志、信息日志
- **实时处理**: 异步事件处理，不影响主业务流程

### 2. 消息主题
- `user_events`: 用户相关事件
- `system_logs`: 系统日志和错误信息

### 3. 消费者服务
- 自动启动用户事件消费者
- 自动启动系统日志消费者
- 支持优雅关闭

## 环境配置

### Docker Compose 配置
项目已配置完整的 Kafka 环境：

```yaml
# Zookeeper 服务
zookeeper:
  image: confluentinc/cp-zookeeper:7.4.0
  container_name: xx-zookeeper
  ports:
    - "2181:2181"

# Kafka 服务  
kafka:
  image: confluentinc/cp-kafka:7.4.0
  container_name: xx-kafka
  ports:
    - "9092:9092"
    - "9101:9101"
  depends_on:
    zookeeper:
      condition: service_healthy
```

### 环境变量
```bash
KAFKA_BROKERS=kafka:29092
KAFKA_TOPIC_USER_EVENTS=user_events
KAFKA_TOPIC_SYSTEM_LOGS=system_logs
```

## API 接口

### 1. 获取 Kafka 状态
```http
GET /api/kafka/status
```

响应示例：
```json
{
  "code": 200,
  "message": "获取成功",
  "data": {
    "status": "enabled"
  }
}
```

### 2. 发送测试消息
```http
POST /api/kafka/test
Authorization: Bearer <token>
Content-Type: application/json

{
  "message": "测试消息内容"
}
```

响应示例：
```json
{
  "code": 200,
  "message": "消息发送成功"
}
```

## 事件类型

### 用户事件 (user_events)
1. **user_login**: 用户登录
   ```json
   {
     "type": "user_login",
     "timestamp": "2024-01-01T12:00:00Z",
     "data": {
       "user_id": 1,
       "username": "admin",
       "ip": "192.168.1.100",
       "action": "login"
     }
   }
   ```

2. **user_logout**: 用户登出
3. **user_register**: 用户注册
4. **user_update**: 用户信息更新
5. **user_delete**: 用户删除

### 系统日志 (system_logs)
1. **system_error**: 系统错误
   ```json
   {
     "type": "system_error",
     "timestamp": "2024-01-01T12:00:00Z",
     "data": {
       "service": "auth",
       "error": "数据库连接失败",
       "details": {...},
       "level": "error"
     }
   }
   ```

2. **system_info**: 系统信息

## 代码结构

```
xx-backend/
├── pkg/kafka/
│   └── kafka.go          # Kafka 客户端实现
├── internal/service/
│   └── kafka_service.go  # Kafka 业务服务
├── internal/handler/
│   └── kafka_handler.go  # Kafka API 处理器
└── config/
    └── config.go         # Kafka 配置
```

## 启动服务

1. **启动所有服务**:
   ```bash
   docker-compose up -d
   ```

2. **查看服务状态**:
   ```bash
   docker-compose ps
   ```

3. **查看 Kafka 日志**:
   ```bash
   docker-compose logs kafka
   ```

## 监控和调试

### 1. 查看 Kafka 主题
```bash
# 进入 Kafka 容器
docker exec -it xx-kafka bash

# 列出所有主题
kafka-topics --bootstrap-server localhost:9092 --list

# 查看主题详情
kafka-topics --bootstrap-server localhost:9092 --describe --topic user_events
```

### 2. 消费消息
```bash
# 消费用户事件
kafka-console-consumer --bootstrap-server localhost:9092 --topic user_events --from-beginning

# 消费系统日志
kafka-console-consumer --bootstrap-server localhost:9092 --topic system_logs --from-beginning
```

### 3. 生产消息
```bash
# 发送测试消息
kafka-console-producer --bootstrap-server localhost:9092 --topic user_events
```

## 故障排除

### 1. Kafka 连接失败
- 检查 Zookeeper 是否正常运行
- 确认网络连接和端口配置
- 查看容器日志: `docker-compose logs kafka`

### 2. 主题创建失败
- 检查 Kafka 服务状态
- 确认权限配置
- 查看错误日志

### 3. 消息发送失败
- 检查生产者配置
- 确认主题存在
- 验证网络连接

## 扩展功能

### 1. 添加新事件类型
1. 在 `kafka_service.go` 中添加新方法
2. 在相应的业务服务中调用
3. 在消费者中添加处理逻辑

### 2. 添加新主题
1. 更新配置文件
2. 在 `CreateTopics()` 方法中添加新主题
3. 创建对应的消费者

### 3. 消息持久化
- 可以将重要事件存储到数据库
- 实现消息重试机制
- 添加死信队列处理

## 性能优化

1. **批量处理**: 支持批量发送消息
2. **连接池**: 复用 Kafka 连接
3. **异步处理**: 非阻塞消息发送
4. **错误重试**: 自动重试失败的消息

## 安全考虑

1. **认证**: 配置 SASL 认证
2. **加密**: 启用 SSL/TLS 加密
3. **权限**: 设置适当的访问权限
4. **审计**: 记录所有操作日志 