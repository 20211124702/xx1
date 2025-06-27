#!/bin/bash

echo "🚀 开始启动 XX-Admin 系统..."

# 1. 先拉取基础镜像
echo "📦 拉取基础镜像..."
docker pull mysql:8.0
docker pull redis:7-alpine
docker pull nginx:alpine

# 2. 启动数据库服务
echo "🗄️ 启动 MySQL..."
docker-compose up -d mysql
sleep 10

echo "🔴 启动 Redis..."
docker-compose up -d redis
sleep 5

# 3. 构建并启动后端
echo "🔧 构建后端服务..."
docker-compose build backend
docker-compose up -d backend
sleep 10

# 4. 构建并启动前端
echo "🎨 构建前端服务..."
docker-compose build frontend
docker-compose up -d frontend

echo "✅ 所有服务启动完成！"
echo "🌐 前端地址: http://localhost"
echo "🔌 后端地址: http://localhost:8080"
echo "📊 数据库: localhost:3306"
echo "🔴 Redis: localhost:6379" 