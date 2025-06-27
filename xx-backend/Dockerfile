# 构建阶段
FROM alibaba-cloud-linux-3-registry.cn-hangzhou.cr.aliyuncs.com/alinux3/golang:1.19.4 AS builder

# 安装必要的包
# RUN apk add --no-cache git ca-certificates tzdata

WORKDIR /app

# 复制go mod文件
COPY go.mod go.sum ./
# RUN go version
# 下载依赖
# RUN go mod download

# 复制源代码
COPY . .

# 构建应用
# RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main .
RUN GOPROXY=https://goproxy.cn,direct go build -o main .
# 生产阶段
FROM alibaba-cloud-linux-3-registry.cn-hangzhou.cr.aliyuncs.com/alinux3/golang:1.19.4

# 安装必要的包
# RUN apk --no-cache add ca-certificates tzdata

WORKDIR /root/

# 从构建阶段复制二进制文件
COPY --from=builder /app/main .

# 创建非root用户
# RUN adduser -D -g '' appuser && chown -R appuser /root
# USER appuser

# 暴露端口
EXPOSE 8080 50051

# 健康检查
HEALTHCHECK --interval=30s --timeout=3s --start-period=5s --retries=3 \
  CMD wget --no-verbose --tries=1 --spider http://localhost:8080/health || exit 1

# 启动应用
CMD ["./main"] 