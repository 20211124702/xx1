#!/bin/bash

echo "🔍 检查 Docker 配置..."

# 检查 Docker 是否运行
if ! docker info > /dev/null 2>&1; then
    echo "❌ Docker 未运行，请先启动 Docker"
    exit 1
fi

echo "✅ Docker 正在运行"

# 检查镜像加速器配置
echo "📋 当前镜像加速器配置："
docker info | grep -A 10 "Registry Mirrors"

# 测试拉取镜像
echo "🧪 测试镜像拉取..."
docker pull hello-world:latest

if [ $? -eq 0 ]; then
    echo "✅ 镜像拉取成功，加速器配置正常"
else
    echo "❌ 镜像拉取失败，请检查网络或加速器配置"
fi

echo "🎯 配置完成！现在可以启动服务了："
echo "docker-compose up -d" 