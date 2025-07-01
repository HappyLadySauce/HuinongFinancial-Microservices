#!/bin/bash

# Docker构建脚本 - HuinongFinancial微服务
# 使用方法：./scripts/docker.sh [service_name] [action] [options]
# 示例：./scripts/docker.sh appuser build v1.0.0
#       ./scripts/docker.sh appuser push v1.0.0
#       ./scripts/docker.sh all build v1.0.0

set -e

# 颜色定义
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
PURPLE='\033[0;35m'
CYAN='\033[0;36m'
NC='\033[0m' # No Color

# 默认配置
DEFAULT_REGISTRY="registry.huinong.internal/huinong"
DEFAULT_VERSION="latest"
BUILD_MODE="local"  # local 或 docker

# 帮助信息
show_help() {
    echo -e "${BLUE}Docker构建脚本 - HuinongFinancial微服务${NC}"
    echo ""
    echo "功能特性:"
    echo "  - 集成build.sh构建二进制文件"
    echo "  - 支持本地构建和容器内构建两种模式"
    echo "  - 支持单服务和批量构建Docker镜像"
    echo "  - 支持镜像推送到私有仓库"
    echo "  - 支持一键部署到Kubernetes"
    echo ""
    echo "使用方法:"
    echo "  $0 [service_name] [action] [version] [options]"
    echo ""
    echo "参数说明:"
    echo "  service_name: 服务名称 (appuser|oauser|loan|loanproduct|leaseproduct|lease|all)"
    echo "  action:       操作类型 (build|push|deploy|all)"
    echo "  version:      镜像版本 (默认: latest)"
    echo "  options:      构建选项"
    echo ""
    echo "操作类型:"
    echo "  build    构建Docker镜像"
    echo "  push     推送镜像到仓库"
    echo "  deploy   部署到Kubernetes"
    echo "  all      构建、推送、部署一条龙"
    echo ""
    echo "示例:"
    echo "  $0 appuser build v1.0.0              # 构建appuser镜像"
    echo "  $0 appuser push v1.0.0               # 推送appuser镜像"
    echo "  $0 appuser deploy v1.0.0             # 部署appuser到k8s"
    echo "  $0 appuser all v1.0.0                # 一键构建推送部署"
    echo "  $0 all build v1.0.0                  # 构建所有服务镜像"
    echo "  $0 all all v1.0.0                    # 所有服务一键构建推送部署"
    echo ""
    echo "环境变量:"
    echo "  DOCKER_REGISTRY: 镜像仓库地址 (默认: $DEFAULT_REGISTRY)"
    echo "  BUILD_MODE:      构建模式 local|docker (默认: local)"
    echo "  K8S_NAMESPACE:   K8s命名空间 (默认: huinong)"
    echo ""
    echo "构建模式:"
    echo "  local   本地构建二进制文件，然后打包到镜像 (推荐，速度快)"
    echo "  docker  容器内构建，环境隔离 (CI/CD推荐)"
    echo ""
}

# 检查参数
if [[ $# -eq 0 ]] || [[ "$1" == "-h" ]] || [[ "$1" == "--help" ]]; then
    show_help
    exit 0
fi

SERVICE_NAME=$1
ACTION=${2:-"build"}
VERSION=${3:-$DEFAULT_VERSION}

# 环境变量
DOCKER_REGISTRY=${DOCKER_REGISTRY:-$DEFAULT_REGISTRY}
BUILD_MODE=${BUILD_MODE:-"local"}
K8S_NAMESPACE=${K8S_NAMESPACE:-"huinong"}

# 支持的服务列表
SERVICES=("appuser" "oauser" "loan" "loanproduct" "leaseproduct" "lease")
ACTIONS=("build" "push" "deploy" "all")

# 获取项目根目录
PROJECT_ROOT=$(cd "$(dirname "$0")/.." && pwd)
DOCKER_DIR="$PROJECT_ROOT/docker"
BIN_DIR="$PROJECT_ROOT/bin"

echo -e "${BLUE}========================================${NC}"
echo -e "${BLUE}    HuinongFinancial Docker构建工具     ${NC}"
echo -e "${BLUE}========================================${NC}"
echo -e "${BLUE}项目根目录: $PROJECT_ROOT${NC}"
echo -e "${BLUE}镜像仓库: $DOCKER_REGISTRY${NC}"
echo -e "${BLUE}构建模式: $BUILD_MODE${NC}"
echo -e "${BLUE}K8s命名空间: $K8S_NAMESPACE${NC}"

# 验证服务名称
if [[ "$SERVICE_NAME" != "all" ]] && [[ ! " ${SERVICES[@]} " =~ " ${SERVICE_NAME} " ]]; then
    echo -e "${RED}错误: 不支持的服务名称 '$SERVICE_NAME'${NC}"
    echo -e "支持的服务: ${SERVICES[*]} all"
    exit 1
fi

# 验证操作类型
if [[ ! " ${ACTIONS[@]} " =~ " ${ACTION} " ]]; then
    echo -e "${RED}错误: 不支持的操作类型 '$ACTION'${NC}"
    echo -e "支持的操作: ${ACTIONS[*]}"
    exit 1
fi

# 检查必要工具
check_tools() {
    # 检查Docker
    if ! command -v docker &> /dev/null; then
        echo -e "${RED}错误: Docker 未安装${NC}"
        exit 1
    fi
    
    # 检查kubectl (如果需要部署)
    if [[ "$ACTION" == "deploy" ]] || [[ "$ACTION" == "all" ]]; then
        if ! command -v kubectl &> /dev/null; then
            echo -e "${YELLOW}警告: kubectl 未安装，无法部署到K8s${NC}"
        fi
    fi
    
    echo -e "${GREEN}✓ 工具检查完成${NC}"
}

# 构建二进制文件
build_binaries() {
    local service=$1
    
    echo -e "${GREEN}构建 $service 二进制文件...${NC}"
    
    # 调用build.sh构建二进制文件
    local build_script="$PROJECT_ROOT/scripts/build.sh"
    if [[ ! -f "$build_script" ]]; then
        echo -e "${RED}错误: build.sh 脚本不存在${NC}"
        return 1
    fi
    
    # 在子shell中构建API和RPC，避免build.sh的exit影响本脚本
    local build_result
    (
        exec "$build_script" "$service" "all" "--release"
    )
    build_result=$?
    
    if [[ $build_result -ne 0 ]]; then
        echo -e "${RED}错误: $service 二进制构建失败${NC}"
        return 1
    fi
    
    # 验证构建结果
    local api_bin="$BIN_DIR/$service-api"
    local rpc_bin="$BIN_DIR/$service-rpc"
    
    if [[ ! -f "$api_bin" ]] || [[ ! -f "$rpc_bin" ]]; then
        echo -e "${RED}错误: $service 二进制文件不存在${NC}"
        return 1
    fi
    
    echo -e "${GREEN}✓ $service 二进制构建完成${NC}"
    return 0
}

# 创建Dockerfile
create_dockerfile() {
    local service=$1
    local service_type=$2  # api 或 rpc
    local dockerfile_path="$DOCKER_DIR/$service-$service_type.Dockerfile"
    
    # 创建docker目录
    mkdir -p "$DOCKER_DIR"
    
    # 生成Dockerfile内容
    cat > "$dockerfile_path" << EOF
# $service $service_type 服务 Dockerfile
# 构建模式: $BUILD_MODE
FROM alpine:latest

# 安装运行时依赖和设置时区
RUN apk add --no-cache ca-certificates tzdata curl \\
    && cp /usr/share/zoneinfo/Asia/Shanghai /etc/localtime \\
    && echo "Asia/Shanghai" > /etc/timezone \\
    && apk del tzdata

# 创建非root用户，提升安全性
RUN addgroup -g 1000 appuser && \\
    adduser -D -s /bin/sh -u 1000 -G appuser appuser

WORKDIR /app

# 复制二进制文件和配置文件
COPY bin/$service-$service_type ./$service-$service_type
COPY app/$service/cmd/$service_type/etc ./etc

# 设置权限
RUN chmod +x ./$service-$service_type && \\
    chown -R appuser:appuser /app

# 切换到非root用户
USER appuser

# 根据服务类型设置端口
EOF

    # 设置端口和启动命令
    if [[ "$service_type" == "api" ]]; then
        cat >> "$dockerfile_path" << EOF
EXPOSE 10001
CMD ["./$service-$service_type", "-f", "etc/$service.yaml"]
EOF
    else
        cat >> "$dockerfile_path" << EOF
EXPOSE 20001
CMD ["./$service-$service_type", "-f", "etc/${service}rpc.yaml"]
EOF
    fi
    
    echo "$dockerfile_path"
}

# 构建Docker镜像
build_docker_image() {
    local service=$1
    local service_type=$2  # api 或 rpc
    
    echo -e "${GREEN}构建 $service-$service_type Docker镜像...${NC}"
    
    local image_name="$DOCKER_REGISTRY/$service-$service_type:$VERSION"
    local dockerfile=""
    
    if [[ "$BUILD_MODE" == "local" ]]; then
        # 本地构建模式：使用预构建的二进制文件
        dockerfile=$(create_dockerfile "$service" "$service_type")
        
        echo -e "${BLUE}使用本地构建模式...${NC}"
        echo -e "${CYAN}docker build -f $dockerfile -t $image_name .${NC}"
        
        if docker build -f "$dockerfile" -t "$image_name" .; then
            echo -e "${GREEN}✓ $service-$service_type 镜像构建成功${NC}"
            echo -e "${CYAN}DEBUG: build_docker_image 清理临时文件并返回 0${NC}"
            # 清理临时Dockerfile
            rm -f "$dockerfile"
            return 0
        else
            echo -e "${RED}✗ $service-$service_type 镜像构建失败${NC}"
            echo -e "${CYAN}DEBUG: build_docker_image 清理临时文件并返回 1${NC}"
            rm -f "$dockerfile"
            return 1
        fi
    else
        # 容器构建模式：使用多阶段构建
        local dockerfile_path="app/$service/docker/$service_type.build.Dockerfile"
        
        if [[ ! -f "$dockerfile_path" ]]; then
            echo -e "${RED}错误: 多阶段构建Dockerfile不存在: $dockerfile_path${NC}"
            echo -e "${YELLOW}提示: 请先创建多阶段构建Dockerfile或使用本地构建模式${NC}"
            return 1
        fi
        
        echo -e "${BLUE}使用容器构建模式...${NC}"
        echo -e "${CYAN}docker build -f $dockerfile_path -t $image_name .${NC}"
        
        if docker build -f "$dockerfile_path" -t "$image_name" .; then
            echo -e "${GREEN}✓ $service-$service_type 镜像构建成功${NC}"
            return 0
        else
            echo -e "${RED}✗ $service-$service_type 镜像构建失败${NC}"
            return 1
        fi
    fi
}

# 推送Docker镜像
push_docker_image() {
    local service=$1
    local service_type=$2
    
    local image_name="$DOCKER_REGISTRY/$service-$service_type:$VERSION"
    
    echo -e "${GREEN}推送 $image_name...${NC}"
    
    if docker push "$image_name"; then
        echo -e "${GREEN}✓ $service-$service_type 镜像推送成功${NC}"
        return 0
    else
        echo -e "${RED}✗ $service-$service_type 镜像推送失败${NC}"
        return 1
    fi
}

# 部署到K8s
deploy_to_k8s() {
    local service=$1
    
    echo -e "${GREEN}部署 $service 到Kubernetes...${NC}"
    
    # 使用现有的K8s部署文件或创建新的
    local deployment_file="k8s/$service-deployment.yaml"
    
    if [[ ! -f "$deployment_file" ]]; then
        echo -e "${YELLOW}未找到现有K8s部署文件，创建新的...${NC}"
        deployment_file=$(create_k8s_deployment "$service")
    fi
    
    # 应用部署
    if kubectl apply -f "$deployment_file"; then
        echo -e "${GREEN}✓ $service 部署成功${NC}"
        
        # 等待部署就绪
        echo -e "${BLUE}等待部署就绪...${NC}"
        if kubectl wait --for=condition=available --timeout=300s deployment/$service-rpc deployment/$service-api -n "$K8S_NAMESPACE" 2>/dev/null; then
            echo -e "${GREEN}✓ $service 服务已就绪${NC}"
        else
            echo -e "${YELLOW}⚠ $service 部署超时，请检查状态${NC}"
        fi
        
        # 显示服务状态
        echo -e "${BLUE}服务状态:${NC}"
        kubectl get pods -n "$K8S_NAMESPACE" -l app=$service-api -o wide 2>/dev/null || true
        kubectl get pods -n "$K8S_NAMESPACE" -l app=$service-rpc -o wide 2>/dev/null || true
        
        return 0
    else
        echo -e "${RED}✗ $service 部署失败${NC}"
        return 1
    fi
}

# 处理单个服务
process_service() {
    local service=$1
    local action=$2
    
    echo -e "${PURPLE}========== 处理 $service 服务 ==========${NC}"
    echo -e "${CYAN}DEBUG: process_service 开始执行，服务=$service，操作=$action${NC}"
    
    local success=true
    
    # 临时禁用严格模式，防止单个操作失败导致函数退出
    set +e
    echo -e "${CYAN}DEBUG: 已禁用严格模式${NC}"
    
    case $action in
        "build")
            # 先构建二进制文件
            if [[ "$BUILD_MODE" == "local" ]]; then
                build_binaries "$service"
                [[ $? -ne 0 ]] && success=false
            fi
            
            # 构建Docker镜像
            if [[ "$success" == true ]]; then
                build_docker_image "$service" "api"
                [[ $? -ne 0 ]] && success=false
                build_docker_image "$service" "rpc"
                [[ $? -ne 0 ]] && success=false
            fi
            ;;
        "push")
            push_docker_image "$service" "api"
            [[ $? -ne 0 ]] && success=false
            push_docker_image "$service" "rpc"
            [[ $? -ne 0 ]] && success=false
            ;;
        "deploy")
            deploy_to_k8s "$service"
            [[ $? -ne 0 ]] && success=false
            ;;
        "all")
            # 先构建二进制文件
            if [[ "$BUILD_MODE" == "local" ]]; then
                build_binaries "$service"
                [[ $? -ne 0 ]] && success=false
            fi
            
            # 构建Docker镜像
            if [[ "$success" == true ]]; then
                build_docker_image "$service" "api"
                [[ $? -ne 0 ]] && success=false
                build_docker_image "$service" "rpc"
                [[ $? -ne 0 ]] && success=false
            fi
            
            # 推送镜像
            if [[ "$success" == true ]]; then
                push_docker_image "$service" "api"
                [[ $? -ne 0 ]] && success=false
                push_docker_image "$service" "rpc"
                [[ $? -ne 0 ]] && success=false
            fi
            
            # 部署到K8s
            if [[ "$success" == true ]]; then
                deploy_to_k8s "$service"
                [[ $? -ne 0 ]] && success=false
            fi
            ;;
    esac
    
    # 恢复严格模式
    set -e
    echo -e "${CYAN}DEBUG: 已恢复严格模式${NC}"
    
    echo -e "${CYAN}DEBUG: process_service 即将结束，success=$success${NC}"
    
    if [[ "$success" == true ]]; then
        echo -e "${GREEN}✓ $service $action 操作成功${NC}"
        echo -e "${CYAN}DEBUG: process_service 返回 0${NC}"
        return 0
    else
        echo -e "${RED}✗ $service $action 操作失败${NC}"
        echo -e "${CYAN}DEBUG: process_service 返回 1${NC}"
        return 1
    fi
}

# 处理所有服务
process_all_services() {
    local action=$1
    
    echo -e "${BLUE}开始处理所有服务...${NC}"
    echo -e "${BLUE}服务列表: ${SERVICES[*]}${NC}"
    
    local total_success=true
    local success_count=0
    local total_count=${#SERVICES[@]}
    
    for service in "${SERVICES[@]}"; do
        echo -e "${CYAN}准备处理服务: $service ($(date))${NC}"
        
        # 使用 set +e 临时禁用严格模式，避免单个服务失败导致整个脚本退出
        set +e
        process_service "$service" "$action"
        local service_result=$?
        set -e
        
        echo -e "${CYAN}$service 处理完毕，返回值: $service_result${NC}"
        
        if [[ $service_result -eq 0 ]]; then
            success_count=$((success_count + 1))
            echo -e "${GREEN}✓ $service 处理成功${NC}"
        else
            total_success=false
            echo -e "${RED}✗ $service 处理失败${NC}"
        fi
        
        echo -e "${CYAN}已处理 $success_count 个服务，继续处理下一个...${NC}"
        echo ""
    done
    
    echo -e "${CYAN}循环结束，开始汇总...${NC}"
    
    echo -e "${BLUE}========== 处理汇总 ==========${NC}"
    echo -e "${BLUE}成功: $success_count/$total_count${NC}"
    
    if [[ "$total_success" == true ]]; then
        echo -e "${GREEN}✓ 所有服务处理成功！${NC}"
    else
        echo -e "${YELLOW}⚠ 部分服务处理失败${NC}"
    fi
    
    [[ "$total_success" == true ]]
}

# 显示镜像信息  
show_images() {
    echo -e "${BLUE}========== 镜像信息 ==========${NC}"
    
    for service in "${SERVICES[@]}"; do
        local api_image="$DOCKER_REGISTRY/$service-api:$VERSION"
        local rpc_image="$DOCKER_REGISTRY/$service-rpc:$VERSION"
        
        if docker images -q "$api_image" &>/dev/null; then
            local api_size=$(docker images --format "table {{.Size}}" "$api_image" | tail -n +2)
            echo -e "${GREEN}  $api_image${NC} (${CYAN}$api_size${NC})"
        fi
        
        if docker images -q "$rpc_image" &>/dev/null; then
            local rpc_size=$(docker images --format "table {{.Size}}" "$rpc_image" | tail -n +2)
            echo -e "${GREEN}  $rpc_image${NC} (${CYAN}$rpc_size${NC})"
        fi
    done
}

# 主函数
main() {
    # 进入项目根目录
    cd "$PROJECT_ROOT"
    
    # 检查工具
    check_tools
    
    # 处理服务
    if [[ "$SERVICE_NAME" == "all" ]]; then
        process_all_services "$ACTION"
    else
        process_service "$SERVICE_NAME" "$ACTION"
    fi
    
    local exit_code=$?
    
    # 显示镜像信息
    if [[ "$ACTION" == "build" ]] || [[ "$ACTION" == "all" ]]; then
        echo ""
        show_images
    fi
    
    echo ""
    echo -e "${GREEN}========================================${NC}"
    if [[ $exit_code -eq 0 ]]; then
        echo -e "${GREEN}       Docker操作完成！               ${NC}"
    else
        echo -e "${RED}       Docker操作失败！               ${NC}"
    fi
    echo -e "${GREEN}========================================${NC}"
    
    exit $exit_code
}

# 执行主函数
main "$@"
