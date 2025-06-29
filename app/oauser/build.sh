#!/bin/bash
set -e

# ==========================================
# OAUser 微服务自动化构建部署脚本
# ==========================================
# 支持的构建模式：
# - local：本地构建二进制文件，然后打包到容器（推荐，速度快）
# - docker：容器内构建，环境隔离（CI/CD 推荐）
#
# 使用方法：
# BUILD_MODE=local ./build.sh v1.0.0 build    # 本地构建
# BUILD_MODE=docker ./build.sh v1.0.0 build   # 容器构建
# ./build.sh v1.0.0 deploy                    # 完整部署流程
# ==========================================

# --- 配置参数 ---
SERVICE_NAME="oauser"
REGISTRY="registry.huinong.internal/huinong"
VERSION=${1:-"latest"}
ACTION=${2:-"build"}

# 构建模式: local(默认) 或 docker
BUILD_MODE=${BUILD_MODE:-"local"}

# --- Helper Functions ---
log_step() { echo -e "\n\e[34m>>> $1\e[0m"; }
log_info() { echo "    $1"; }
log_success() { echo -e "\e[32m✅ $1\e[0m"; }
log_error() { echo -e "\e[31m❌ $1\e[0m"; exit 1; }

# --- 验证参数 ---
validate_params() {
    if [[ ! "$VERSION" =~ ^v[0-9]+\.[0-9]+\.[0-9]+$ ]] && [[ "$VERSION" != "latest" ]]; then
        log_error "版本号格式错误，请使用 vx.y.z 格式，例如：v1.0.0"
    fi
    
    if [[ ! "$ACTION" =~ ^(build|push|deploy|all)$ ]]; then
        log_error "操作类型错误，支持：build, push, deploy, all"
    fi
}

# --- 本地构建二进制文件 ---
build_binaries() {
    log_step "构建 $SERVICE_NAME 二进制文件 (本地模式)"
    
    # 设置编译环境
    export CGO_ENABLED=0
    export GOOS=linux
    export GOARCH=amd64
    
    # 构建 API 服务
    log_info "构建 API 二进制文件..."
    (cd app/$SERVICE_NAME/cmd/api && go build -ldflags="-s -w" -o ../../../../$SERVICE_NAME-api .)
    
    # 构建 RPC 服务
    log_info "构建 RPC 二进制文件..."
    (cd app/$SERVICE_NAME/cmd/rpc && go build -ldflags="-s -w" -o ../../../../$SERVICE_NAME-rpc .)
    
    log_success "二进制文件构建完成"
}

# --- 构建 Docker 镜像 ---
build_images() {
    log_step "构建 Docker 镜像 - $SERVICE_NAME:$VERSION (模式: $BUILD_MODE)"
    
    if [ "$BUILD_MODE" = "local" ]; then
        # 本地构建模式：先构建二进制文件，再打包容器
        build_binaries
        
        log_info "使用本地二进制文件构建镜像..."
        docker build -f ./app/$SERVICE_NAME/docker/api.Dockerfile \
            -t $REGISTRY/$SERVICE_NAME-api:$VERSION .
        docker build -f ./app/$SERVICE_NAME/docker/rpc.Dockerfile \
            -t $REGISTRY/$SERVICE_NAME-rpc:$VERSION .
        
        # 清理临时文件
        log_info "清理临时二进制文件..."
        rm -f $SERVICE_NAME-api $SERVICE_NAME-rpc
        
    else
        # 容器构建模式：在容器内构建
        log_info "使用容器内构建模式..."
        docker build -f ./app/$SERVICE_NAME/docker/api.build.Dockerfile \
            -t $REGISTRY/$SERVICE_NAME-api:$VERSION .
        docker build -f ./app/$SERVICE_NAME/docker/rpc.build.Dockerfile \
            -t $REGISTRY/$SERVICE_NAME-rpc:$VERSION .
    fi
    
    log_success "Docker 镜像构建完成"
}

# --- 推送镜像 ---
push_images() {
    log_step "推送镜像到仓库"
    
    log_info "推送 $REGISTRY/$SERVICE_NAME-api:$VERSION"
    docker push $REGISTRY/$SERVICE_NAME-api:$VERSION
    
    log_info "推送 $REGISTRY/$SERVICE_NAME-rpc:$VERSION"
    docker push $REGISTRY/$SERVICE_NAME-rpc:$VERSION
    
    log_success "镜像推送完成"
}

# --- 部署到 Kubernetes ---
deploy_k8s() {
    log_step "部署到 Kubernetes"
    
    if [ ! -f "app/$SERVICE_NAME/k8s/$SERVICE_NAME-deployment.yaml" ]; then
        log_error "找不到 K8s 部署文件：app/$SERVICE_NAME/k8s/$SERVICE_NAME-deployment.yaml"
    fi
    
    # 替换镜像版本
    sed -i.bak "s|image: .*/$SERVICE_NAME-api:.*|image: $REGISTRY/$SERVICE_NAME-api:$VERSION|g" \
        app/$SERVICE_NAME/k8s/$SERVICE_NAME-deployment.yaml
    sed -i.bak "s|image: .*/$SERVICE_NAME-rpc:.*|image: $REGISTRY/$SERVICE_NAME-rpc:$VERSION|g" \
        app/$SERVICE_NAME/k8s/$SERVICE_NAME-deployment.yaml
    
    # 应用部署
    log_info "应用 K8s 配置..."
    kubectl apply -f app/$SERVICE_NAME/k8s/$SERVICE_NAME-deployment.yaml
    
    # 恢复原始文件
    mv app/$SERVICE_NAME/k8s/$SERVICE_NAME-deployment.yaml.bak \
       app/$SERVICE_NAME/k8s/$SERVICE_NAME-deployment.yaml
    
    log_success "Kubernetes 部署完成"
}

# --- 主要执行逻辑 ---
main() {
    validate_params
    
    echo "=========================================="
    echo "🚀 OAUser 微服务构建部署"
    echo "=========================================="
    echo "服务名称: $SERVICE_NAME"
    echo "版本号:   $VERSION"
    echo "构建模式: $BUILD_MODE"
    echo "执行操作: $ACTION"
    echo "=========================================="
    
    case $ACTION in
        "build")
            build_images
            ;;
        "push")
            push_images
            ;;
        "deploy")
            deploy_k8s
            ;;
        "all")
            build_images
            push_images
            deploy_k8s
            ;;
    esac
    
    log_success "🎉 操作完成！服务：$SERVICE_NAME:$VERSION"
}

# --- 执行入口 ---
main "$@" 