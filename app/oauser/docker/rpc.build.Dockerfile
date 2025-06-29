# OAUser RPC 服务 Dockerfile (容器内构建版本)
FROM golang:1.24-alpine AS builder

# 国内环境优化
ENV GOPROXY=https://goproxy.cn,direct

WORKDIR /app

# 1. 复制所有 go.mod/go.sum, 利用依赖缓存
COPY app/oauser/cmd/model/go.* ./app/oauser/cmd/model/
COPY app/oauser/cmd/rpc/go.* ./app/oauser/cmd/rpc/
COPY app/oauser/cmd/api/go.* ./app/oauser/cmd/api/

# 2. 在 RPC 模块目录中下载依赖
WORKDIR /app/app/oauser/cmd/rpc
RUN go mod tidy

# 3. 复制服务源码
WORKDIR /app
COPY app/oauser/ ./app/oauser/

# 4. 构建二进制文件
WORKDIR /app/app/oauser/cmd/rpc
RUN CGO_ENABLED=0 GOOS=linux go build -o /app/oauser-rpc .

# --- 运行阶段 ---
FROM alpine:latest

# 安装运行时依赖和设置时区
RUN apk add --no-cache ca-certificates tzdata \
    && cp /usr/share/zoneinfo/Asia/Shanghai /etc/localtime \
    && echo "Asia/Shanghai" > /etc/timezone \
    && apk del tzdata

# 创建非 root 用户，提升安全性
RUN addgroup -g 1000 oauser && \
    adduser -D -s /bin/sh -u 1000 -G oauser oauser

WORKDIR /app

# 复制构建的二进制文件和配置文件
COPY --from=builder /app/oauser-rpc .
COPY app/oauser/cmd/rpc/etc ./etc

# 设置权限
RUN chmod +x ./oauser-rpc && \
    chown -R oauser:oauser /app

USER oauser
EXPOSE 20002
CMD ["./oauser-rpc", "-f", "etc/oauserrpc.yaml"] 