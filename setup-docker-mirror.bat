@echo off
echo 🚀 配置 Docker 镜像加速器...

REM 检查 Docker 是否安装
docker --version >nul 2>&1
if errorlevel 1 (
    echo ❌ Docker 未安装，请先安装 Docker Desktop
    pause
    exit /b 1
)

echo ✅ Docker 已安装

REM 创建配置目录
if not exist "%USERPROFILE%\.docker" mkdir "%USERPROFILE%\.docker"

REM 创建 daemon.json 配置文件
echo 正在创建 Docker 配置文件...
(
echo {
echo   "registry-mirrors": [
echo     "https://registry.cn-hangzhou.aliyuncs.com",
echo     "https://docker.mirrors.ustc.edu.cn",
echo     "https://hub-mirror.c.163.com"
echo   ],
echo   "insecure-registries": [],
echo   "debug": false,
echo   "experimental": false
echo }
) > "%USERPROFILE%\.docker\daemon.json"

echo ✅ 配置文件已创建: %USERPROFILE%\.docker\daemon.json

echo.
echo 📋 请按以下步骤操作：
echo 1. 打开 Docker Desktop
echo 2. 进入 Settings ^> Docker Engine
echo 3. 将配置文件内容复制到编辑器中
echo 4. 点击 "Apply ^& Restart"
echo 5. 等待 Docker 重启完成
echo.
echo 🎯 配置完成后，运行以下命令启动服务：
echo docker-compose up -d
echo.
pause 