.PHONY: help build up down dev clean logs

# 默认目标
help:
	@echo "可用的命令:"
	@echo "  build    - 构建所有Docker镜像"
	@echo "  up       - 启动生产环境"
	@echo "  down     - 停止所有服务"
	@echo "  dev      - 启动开发环境"
	@echo "  dev-down - 停止开发环境"
	@echo "  clean    - 清理所有容器和镜像"
	@echo "  logs     - 查看服务日志"
	@echo "  frontend - 构建前端镜像"
	@echo "  backend  - 构建后端镜像"

# 构建所有镜像
build:
	docker-compose build

# 启动生产环境
up:
	docker-compose up -d

# 停止所有服务
down:
	docker-compose down

# 启动开发环境
dev:
	docker-compose -f docker-compose.dev.yml up -d

# 停止开发环境
dev-down:
	docker-compose -f docker-compose.dev.yml down

# 清理所有容器和镜像
clean:
	docker-compose down -v --rmi all
	docker system prune -f

# 查看服务日志
logs:
	docker-compose logs -f

# 构建前端镜像
frontend:
	docker build -t xx-admin ./xx-admin

# 构建后端镜像
backend:
	docker build -t xx-backend ./xx-backend

# 进入容器
shell-backend:
	docker-compose exec backend sh

shell-frontend:
	docker-compose exec frontend sh

shell-mysql:
	docker-compose exec mysql mysql -u xx_user -p xx_pass123 xx_admin

# 数据库备份
backup:
	docker-compose exec mysql mysqldump -u xx_user -p xx_pass123 xx_admin > backup_$(shell date +%Y%m%d_%H%M%S).sql

# 数据库恢复
restore:
	@read -p "请输入备份文件名: " file; \
	docker-compose exec -T mysql mysql -u xx_user -p xx_pass123 xx_admin < $$file

# 健康检查
health:
	@echo "检查服务状态..."
	@docker-compose ps
	@echo "\n检查健康状态..."
	@curl -f http://localhost/health || echo "前端服务未就绪"
	@curl -f http://localhost:8080/health || echo "后端服务未就绪" 