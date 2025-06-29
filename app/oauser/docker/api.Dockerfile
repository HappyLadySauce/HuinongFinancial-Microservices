# OAUser API 服务 Dockerfile (本地构建版本)
# 假设二进制文件 'oauser-api' 已在项目根目录构建好
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

# 复制预构建的二进制文件和配置文件
COPY oauser-api .
COPY app/oauser/cmd/api/etc ./etc

# 设置权限
RUN chmod +x ./oauser-api && \
    chown -R oauser:oauser /app

USER oauser
EXPOSE 10002
CMD ["./oauser-api", "-f", "etc/oauser.yaml"] 