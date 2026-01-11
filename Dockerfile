# 使用官方Go镜像作为构建环境
FROM golang:1.21-alpine AS builder

# 设置工作目录
WORKDIR /app

# 复制go.mod和go.sum文件
COPY go.mod ./

# 下载依赖
RUN go mod download

# 复制源代码
COPY . .

# 构建Go程序
RUN go build -o main .

# 使用轻量级Alpine镜像作为运行环境
FROM alpine:latest

# 安装必要的运行时依赖
RUN apk --no-cache add ca-certificates

# 设置工作目录
WORKDIR /root/

# 从构建阶段复制可执行文件
COPY --from=builder /app/main .

# 复制配置文件
COPY --from=builder /app/config.yaml .

# 暴露端口（根据配置文件中的server.port设置）
EXPOSE 8080

# 设置环境变量
ENV LOCAL_CONFIG_PATH=/root/config.yaml

# 运行程序
CMD ["./main"]