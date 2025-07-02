#!/bin/bash

# 微服务启动停止脚本
# 作者: Auto Generated
# 描述: 管理汇农金融微服务的启动、停止和状态查看

# 配置区域
BASE_DIR="/root/HuinongFinancial-Microservices"
BIN_DIR="$BASE_DIR/bin"
PID_DIR="/tmp/huinong_pids"
LOG_DIR="/var/log/huinong"
DOCKER_NETWORK="huinong-network"

# 启动模式 (default: local)
STARTUP_MODE="local"

# 服务配置映射 (服务名:配置文件路径:工作目录)
declare -A SERVICE_CONFIG=(
    ["appuser-rpc"]="$BASE_DIR/app/appuser/cmd/rpc/etc/appuserrpc.yaml:$BASE_DIR/app/appuser/cmd/rpc"
    ["appuser-api"]="$BASE_DIR/app/appuser/cmd/api/etc/appuser.yaml:$BASE_DIR/app/appuser/cmd/api"
    ["oauser-rpc"]="$BASE_DIR/app/oauser/cmd/rpc/etc/oauserrpc.yaml:$BASE_DIR/app/oauser/cmd/rpc"
    ["oauser-api"]="$BASE_DIR/app/oauser/cmd/api/etc/oauser.yaml:$BASE_DIR/app/oauser/cmd/api"
    ["loanproduct-rpc"]="$BASE_DIR/app/loanproduct/cmd/rpc/etc/loanproductrpc.yaml:$BASE_DIR/app/loanproduct/cmd/rpc"
    ["loanproduct-api"]="$BASE_DIR/app/loanproduct/cmd/api/etc/loanproduct.yaml:$BASE_DIR/app/loanproduct/cmd/api"
    ["leaseproduct-rpc"]="$BASE_DIR/app/leaseproduct/cmd/rpc/etc/leaseproductrpc.yaml:$BASE_DIR/app/leaseproduct/cmd/rpc"
    ["leaseproduct-api"]="$BASE_DIR/app/leaseproduct/cmd/api/etc/leaseproduct.yaml:$BASE_DIR/app/leaseproduct/cmd/api"
    ["loan-rpc"]="$BASE_DIR/app/loan/cmd/rpc/etc/loanrpc.yaml:$BASE_DIR/app/loan/cmd/rpc"
    ["loan-api"]="$BASE_DIR/app/loan/cmd/api/etc/loan.yaml:$BASE_DIR/app/loan/cmd/api"
    ["lease-rpc"]="$BASE_DIR/app/lease/cmd/rpc/etc/leaserpc.yaml:$BASE_DIR/app/lease/cmd/rpc"
    ["lease-api"]="$BASE_DIR/app/lease/cmd/api/etc/lease.yaml:$BASE_DIR/app/lease/cmd/api"
)

# Docker 镜像配置映射 (服务名:镜像名:端口)
declare -A DOCKER_CONFIG=(
    ["appuser-rpc"]="registry.huinong.internal/huinong/appuser-rpc:latest:20001"
    ["appuser-api"]="registry.huinong.internal/huinong/appuser-api:latest:10001"
    ["oauser-rpc"]="registry.huinong.internal/huinong/oauser-rpc:latest:20002"
    ["oauser-api"]="registry.huinong.internal/huinong/oauser-api:latest:10002"
    ["loanproduct-rpc"]="registry.huinong.internal/huinong/loanproduct-rpc:latest:20003"
    ["loanproduct-api"]="registry.huinong.internal/huinong/loanproduct-api:latest:10003"
    ["leaseproduct-rpc"]="registry.huinong.internal/huinong/leaseproduct-rpc:latest:20004"
    ["leaseproduct-api"]="registry.huinong.internal/huinong/leaseproduct-api:latest:10004"
    ["loan-rpc"]="registry.huinong.internal/huinong/loan-rpc:latest:20005"
    ["loan-api"]="registry.huinong.internal/huinong/loan-api:latest:10005"
    ["lease-rpc"]="registry.huinong.internal/huinong/lease-rpc:latest:20006"
    ["lease-api"]="registry.huinong.internal/huinong/lease-api:latest:10006"
)

# 新的启动顺序配置 - 先启动RPC，再启动API
declare -a FIRST_RPC_GROUP=("appuser-rpc" "oauser-rpc")
declare -a SECOND_RPC_GROUP=("loanproduct-rpc" "leaseproduct-rpc")
declare -a THIRD_RPC_GROUP=("loan-rpc" "lease-rpc")
declare -a FIRST_API_GROUP=("appuser-api" "oauser-api")
declare -a SECOND_API_GROUP=("loanproduct-api" "leaseproduct-api")
declare -a THIRD_API_GROUP=("loan-api" "lease-api")

# 颜色定义
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# 创建必要目录
create_dirs() {
    mkdir -p "$PID_DIR"
    mkdir -p "$LOG_DIR"
}

# 打印带颜色的消息
print_message() {
    local color=$1
    local message=$2
    echo -e "${color}${message}${NC}"
}

# 检查服务是否存在
check_service_exists() {
    local service=$1
    if [[ ! -f "$BIN_DIR/$service" ]]; then
        print_message $RED "错误: 服务 $service 不存在于 $BIN_DIR"
        return 1
    fi
    
    # 检查配置文件是否存在
    local config_info="${SERVICE_CONFIG[$service]}"
    if [[ -z "$config_info" ]]; then
        print_message $RED "错误: 服务 $service 的配置信息未找到"
        return 1
    fi
    
    local config_file="${config_info%%:*}"
    if [[ ! -f "$config_file" ]]; then
        print_message $RED "错误: 配置文件 $config_file 不存在"
        return 1
    fi
    
    return 0
}

# 启动单个服务
start_service() {
    local service=$1
    local pid_file="$PID_DIR/$service.pid"
    
    # 检查服务是否已经运行
    if [[ -f "$pid_file" ]]; then
        local pid=$(cat "$pid_file")
        if kill -0 "$pid" 2>/dev/null; then
            print_message $YELLOW "服务 $service 已经在运行中 (PID: $pid)"
            return 0
        else
            # PID文件存在但进程不存在，删除旧的PID文件
            rm -f "$pid_file"
        fi
    fi
    
    # 检查服务文件是否存在
    if ! check_service_exists "$service"; then
        return 1
    fi
    
    print_message $BLUE "正在启动服务: $service"
    
    # 获取配置信息
    local config_info="${SERVICE_CONFIG[$service]}"
    local config_file="${config_info%%:*}"
    local work_dir="${config_info##*:}"
    
    # 启动服务 - 从工作目录运行，使用-f参数指定配置文件
    cd "$work_dir" || {
        print_message $RED "错误: 无法切换到工作目录 $work_dir"
        return 1
    }
    
    nohup "$BIN_DIR/$service" -f "$config_file" > "$LOG_DIR/$service.log" 2>&1 &
    local pid=$!
    
    # 保存PID
    echo "$pid" > "$pid_file"
    
    # 等待服务启动
    sleep 1
    
    # 检查服务是否成功启动
    if kill -0 "$pid" 2>/dev/null; then
        print_message $GREEN "✓ 服务 $service 启动成功 (PID: $pid)"
        return 0
    else
        print_message $RED "✗ 服务 $service 启动失败"
        print_message $YELLOW "  错误日志: $LOG_DIR/$service.log"
        print_message $YELLOW "  最后几行日志："
        tail -5 "$LOG_DIR/$service.log" | sed 's/^/    /'
        rm -f "$pid_file"
        return 1
    fi
}

# 停止单个服务
stop_service() {
    local service=$1
    local pid_file="$PID_DIR/$service.pid"
    
    if [[ ! -f "$pid_file" ]]; then
        print_message $YELLOW "服务 $service 没有运行"
        return 0
    fi
    
    local pid=$(cat "$pid_file")
    if ! kill -0 "$pid" 2>/dev/null; then
        print_message $YELLOW "服务 $service 进程不存在，清理PID文件"
        rm -f "$pid_file"
        return 0
    fi
    
    print_message $BLUE "正在停止服务: $service (PID: $pid)"
    
    # 优雅停止
    kill -TERM "$pid" 2>/dev/null
    
    # 等待进程停止
    local count=0
    while kill -0 "$pid" 2>/dev/null && [[ $count -lt 10 ]]; do
        sleep 1
        ((count++))
    done
    
    # 如果进程仍在运行，强制杀死
    if kill -0 "$pid" 2>/dev/null; then
        print_message $YELLOW "优雅停止失败，强制杀死进程"
        kill -KILL "$pid" 2>/dev/null
        sleep 1
    fi
    
    # 清理PID文件
    rm -f "$pid_file"
    
    if ! kill -0 "$pid" 2>/dev/null; then
        print_message $GREEN "✓ 服务 $service 已停止"
        return 0
    else
        print_message $RED "✗ 服务 $service 停止失败"
        return 1
    fi
}

# 启动服务组
start_group() {
    local group_name=$1
    shift
    local services=("$@")
    
    print_message $BLUE "\n========== 启动 $group_name ==========\n"
    
    for service in "${services[@]}"; do
        start_service "$service"
        sleep 2
    done
    
    # 等待服务组稳定
    print_message $YELLOW "等待服务组稳定..."
    sleep 5
}

# 检查服务状态
check_status() {
    local service=$1
    local pid_file="$PID_DIR/$service.pid"
    
    if [[ ! -f "$pid_file" ]]; then
        echo -e "$service: ${RED}未运行${NC}"
        return 1
    fi
    
    local pid=$(cat "$pid_file")
    if kill -0 "$pid" 2>/dev/null; then
        echo -e "$service: ${GREEN}运行中${NC} (PID: $pid)"
        return 0
    else
        echo -e "$service: ${RED}已停止${NC} (PID文件存在但进程不存在)"
        rm -f "$pid_file"
        return 1
    fi
}

# 清理指定服务的日志文件
clean_service_logs() {
    local service=$1
    print_message $YELLOW "清理服务 $service 的日志文件..."
    
    # 清理脚本生成的该服务日志文件
    if [[ -d "$LOG_DIR" ]]; then
        rm -f "$LOG_DIR/$service.log"
        print_message $GREEN "✓ 脚本日志文件清理完成: $LOG_DIR/$service.log"
    fi
    
    # 清理该服务目录下的logs文件夹
    local config_info="${SERVICE_CONFIG[$service]}"
    local work_dir="${config_info##*:}"
    local logs_dir="$work_dir/logs"
    
    if [[ -d "$logs_dir" ]]; then
        rm -rf "$logs_dir"
        print_message $GREEN "✓ 服务日志清理完成: $service ($logs_dir)"
    fi
}

# 清理所有日志文件
clean_logs() {
    print_message $YELLOW "清理旧的日志文件..."
    
    # 清理脚本生成的日志文件
    if [[ -d "$LOG_DIR" ]]; then
        rm -f "$LOG_DIR"/*.log
        print_message $GREEN "✓ 脚本日志文件清理完成: $LOG_DIR"
    fi
    
    # 清理各个微服务目录下的logs文件夹
    for service in "${!SERVICE_CONFIG[@]}"; do
        local config_info="${SERVICE_CONFIG[$service]}"
        local work_dir="${config_info##*:}"
        local logs_dir="$work_dir/logs"
        
        if [[ -d "$logs_dir" ]]; then
            rm -rf "$logs_dir"
            print_message $GREEN "✓ 服务日志清理完成: $service ($logs_dir)"
        fi
    done
    
    print_message $GREEN "✓ 所有日志文件清理完成"
}

# ============ Docker 相关函数 ============

# 检查Docker是否可用
check_docker() {
    if ! command -v docker &> /dev/null; then
        print_message $RED "错误: Docker未安装或不在PATH中"
        return 1
    fi
    
    if ! docker info >/dev/null 2>&1; then
        print_message $RED "错误: Docker服务未运行或无权限访问"
        print_message $YELLOW "请检查:"
        print_message $YELLOW "  1. Docker服务是否启动: sudo systemctl status docker"
        print_message $YELLOW "  2. 当前用户是否在docker组: sudo usermod -aG docker \$USER"
        return 1
    fi
    
    return 0
}

# 创建Docker网络
create_docker_network() {
    if ! docker network ls | grep -q "$DOCKER_NETWORK"; then
        print_message $BLUE "创建Docker网络: $DOCKER_NETWORK"
        docker network create "$DOCKER_NETWORK" >/dev/null 2>&1
        if [[ $? -eq 0 ]]; then
            print_message $GREEN "✓ Docker网络创建成功"
        else
            print_message $RED "✗ Docker网络创建失败"
            return 1
        fi
    else
        print_message $YELLOW "Docker网络已存在: $DOCKER_NETWORK"
    fi
    return 0
}

# 启动Docker服务
start_docker_service() {
    local service=$1
    
    # 检查Docker环境
    if ! check_docker; then
        return 1
    fi
    
    # 检查服务配置是否存在
    local docker_info="${DOCKER_CONFIG[$service]}"
    if [[ -z "$docker_info" ]]; then
        print_message $RED "错误: 服务 $service 的Docker配置信息未找到"
        return 1
    fi
    
    # 解析Docker配置 (格式: image:tag:port)
    local image_and_tag="${docker_info%:*}"
    local port="${docker_info##*:}"
    
    # 检查容器是否已经在运行
    if docker ps | grep -q "$service"; then
        print_message $YELLOW "Docker容器 $service 已经在运行中"
        return 0
    fi
    
    # 停止并删除已存在的容器（如果存在）
    if docker ps -a | grep -q "$service"; then
        print_message $YELLOW "清理旧的Docker容器: $service"
        docker stop "$service" >/dev/null 2>&1
        docker rm "$service" >/dev/null 2>&1
    fi
    
    print_message $BLUE "正在启动Docker服务: $service"
    
    # 获取配置文件路径
    local config_info="${SERVICE_CONFIG[$service]}"
    local config_file="${config_info%%:*}"
    local work_dir="${config_info##*:}"
    
    # 检查配置文件是否存在
    if [[ ! -f "$config_file" ]]; then
        print_message $RED "错误: 配置文件 $config_file 不存在"
        return 1
    fi
    
    # 启动Docker容器
    # 确保日志目录存在并设置正确权限
    mkdir -p "$work_dir/logs"
    
    # 设置日志目录权限（Docker容器内 appuser 的 UID=1000, GID=1000）
    chown -R 1000:1000 "$work_dir/logs" 2>/dev/null || true
    chmod -R 755 "$work_dir/logs" 2>/dev/null || true
    
    # 尝试多种可能的可执行文件路径
    local possible_paths=(
        "/app/$service"
        "/$service" 
        "/usr/local/bin/$service"
        "/bin/$service"
        "./$service"
    )
    
    local success=false
    local container_started=false
    
    # 先尝试默认方式启动（让镜像使用默认的ENTRYPOINT/CMD）
    print_message $YELLOW "  尝试默认方式启动..."
    docker run -d \
        --name "$service" \
        --network "$DOCKER_NETWORK" \
        -p "$port:$port" \
        -v "$config_file":/app/etc/config.yaml:ro \
        -v "$work_dir/logs":/app/logs \
        --restart unless-stopped \
        --user root \
        --entrypoint="" \
        "$image_and_tag" \
        /bin/sh -c "mkdir -p /app/logs && chown -R appuser:appuser /app/logs && su appuser -c '/app/$service -f /app/etc/config.yaml'" >/dev/null 2>&1
    
    if [[ $? -eq 0 ]]; then
        container_started=true
        sleep 2
        # 检查容器是否仍在运行
        if docker ps | grep -q "$service"; then
            success=true
        else
            print_message $YELLOW "  默认启动失败，尝试指定可执行文件路径..."
            docker rm "$service" >/dev/null 2>&1
            container_started=false
        fi
    fi
    
    # 如果默认方式失败，尝试指定可执行文件路径
    if [[ "$container_started" == false ]]; then
        for exec_path in "${possible_paths[@]}"; do
            print_message $YELLOW "  尝试路径: $exec_path"
            
            # 清理之前的容器
            docker rm "$service" >/dev/null 2>&1
            
            docker run -d \
                --name "$service" \
                --network "$DOCKER_NETWORK" \
                -p "$port:$port" \
                -v "$config_file":/app/etc/config.yaml:ro \
                -v "$work_dir/logs":/app/logs \
                --restart unless-stopped \
                --workdir /app \
                --user root \
                --entrypoint="" \
                "$image_and_tag" \
                /bin/sh -c "mkdir -p /app/logs && chown -R appuser:appuser /app/logs && su appuser -c '$exec_path -f /app/etc/config.yaml'" >/dev/null 2>&1
            
            if [[ $? -eq 0 ]]; then
                sleep 2
                if docker ps | grep -q "$service"; then
                    success=true
                    print_message $GREEN "  ✓ 成功使用路径: $exec_path"
                    break
                fi
            fi
        done
    fi
    
    # 检查最终结果
    if [[ "$success" == true ]]; then
        print_message $GREEN "✓ Docker服务 $service 启动成功"
        return 0
    else
        print_message $RED "✗ Docker服务 $service 启动失败"
        print_message $YELLOW "  容器日志："
        if docker ps -a | grep -q "$service"; then
            docker logs "$service" 2>&1 | tail -10 | sed 's/^/    /'
        else
            print_message $YELLOW "    无法获取容器日志（容器未创建）"
        fi
        print_message $YELLOW "  建议运行诊断: $0 --docker diagnose $service"
        return 1
    fi
}

# 停止Docker服务
stop_docker_service() {
    local service=$1
    
    if ! docker ps | grep -q "$service"; then
        print_message $YELLOW "Docker容器 $service 没有运行"
        return 0
    fi
    
    print_message $BLUE "正在停止Docker服务: $service"
    
    # 优雅停止容器
    docker stop "$service" >/dev/null 2>&1
    
    if [[ $? -eq 0 ]]; then
        # 删除容器
        docker rm "$service" >/dev/null 2>&1
        print_message $GREEN "✓ Docker服务 $service 已停止"
        return 0
    else
        print_message $RED "✗ Docker服务 $service 停止失败"
        return 1
    fi
}

# 检查Docker服务状态
check_docker_status() {
    local service=$1
    
    if docker ps | grep -q "$service"; then
        local container_id=$(docker ps | grep "$service" | awk '{print $1}')
        echo -e "$service: ${GREEN}运行中${NC} (容器ID: $container_id)"
        return 0
    elif docker ps -a | grep -q "$service"; then
        echo -e "$service: ${RED}已停止${NC} (容器存在但未运行)"
        return 1
    else
        echo -e "$service: ${RED}未创建${NC}"
        return 1
    fi
}

# Docker模式启动服务组
start_docker_group() {
    local group_name=$1
    shift
    local services=("$@")
    
    print_message $BLUE "\n========== 启动 $group_name (Docker模式) ==========\n"
    
    for service in "${services[@]}"; do
        start_docker_service "$service"
        sleep 2
    done
    
    # 等待服务组稳定
    print_message $YELLOW "等待服务组稳定..."
    sleep 5
}

# Docker模式启动所有服务
start_all_docker() {
    print_message $GREEN "\n========================================"
    print_message $GREEN "    开始启动汇农金融微服务 (Docker模式)"
    print_message $GREEN "========================================\n"
    
    # 检查Docker环境
    if ! check_docker; then
        return 1
    fi
    
    # 创建Docker网络
    if ! create_docker_network; then
        print_message $RED "Docker网络创建失败，无法继续"
        return 1
    fi
    
    # 清理日志
    clean_logs
    
    # 第一阶段：启动RPC服务
    print_message $YELLOW "阶段1: 启动用户管理RPC服务"
    start_docker_group "用户管理RPC服务" "${FIRST_RPC_GROUP[@]}"
    
    print_message $YELLOW "阶段2: 启动产品管理RPC服务"
    start_docker_group "产品管理RPC服务" "${SECOND_RPC_GROUP[@]}"
    
    print_message $YELLOW "阶段3: 启动业务RPC服务"
    start_docker_group "业务RPC服务" "${THIRD_RPC_GROUP[@]}"
    
    # 第二阶段：启动API服务
    print_message $YELLOW "阶段4: 启动用户管理API服务"
    start_docker_group "用户管理API服务" "${FIRST_API_GROUP[@]}"
    
    print_message $YELLOW "阶段5: 启动产品管理API服务"
    start_docker_group "产品管理API服务" "${SECOND_API_GROUP[@]}"
    
    print_message $YELLOW "阶段6: 启动业务API服务"
    start_docker_group "业务API服务" "${THIRD_API_GROUP[@]}"
    
    print_message $GREEN "\n========================================"
    print_message $GREEN "      所有服务启动完成 (Docker模式)"
    print_message $GREEN "========================================\n"
    
    # 显示最终状态
    show_docker_status
}

# Docker模式停止所有服务
stop_all_docker() {
    print_message $RED "\n========================================"
    print_message $RED "    开始停止汇农金融微服务 (Docker模式)"
    print_message $RED "========================================\n"
    
    # 按相反顺序停止服务 - 先停API，再停RPC
    local all_services=()
    all_services+=("${THIRD_API_GROUP[@]}")
    all_services+=("${SECOND_API_GROUP[@]}")
    all_services+=("${FIRST_API_GROUP[@]}")
    all_services+=("${THIRD_RPC_GROUP[@]}")
    all_services+=("${SECOND_RPC_GROUP[@]}")
    all_services+=("${FIRST_RPC_GROUP[@]}")
    
    for service in "${all_services[@]}"; do
        stop_docker_service "$service"
        sleep 1
    done
    
    print_message $RED "\n========================================"
    print_message $RED "      所有服务已停止 (Docker模式)"
    print_message $RED "========================================\n"
}

# 显示Docker模式所有服务状态
show_docker_status() {
    print_message $BLUE "\n========================================"
    print_message $BLUE "      服务状态 (Docker模式)"
    print_message $BLUE "========================================\n"
    
    print_message $BLUE "RPC服务:"
    for service in "${FIRST_RPC_GROUP[@]}" "${SECOND_RPC_GROUP[@]}" "${THIRD_RPC_GROUP[@]}"; do
        echo "  $(check_docker_status "$service")"
    done
    
    print_message $BLUE "\nAPI服务:"
    for service in "${FIRST_API_GROUP[@]}" "${SECOND_API_GROUP[@]}" "${THIRD_API_GROUP[@]}"; do
        echo "  $(check_docker_status "$service")"
    done
    
    print_message $BLUE "\nDocker网络信息:"
    if docker network ls | grep -q "$DOCKER_NETWORK"; then
        echo -e "  网络: ${GREEN}$DOCKER_NETWORK${NC} (已创建)"
    else
        echo -e "  网络: ${RED}$DOCKER_NETWORK${NC} (未创建)"
    fi
    echo
}

# Docker模式重启所有服务
restart_all_docker() {
    print_message $YELLOW "\n========================================"
    print_message $YELLOW "    重启汇农金融微服务 (Docker模式)"
    print_message $YELLOW "========================================\n"
    
    stop_all_docker
    sleep 5
    start_all_docker
}

# 清理Docker资源
cleanup_docker() {
    print_message $YELLOW "\n========================================"
    print_message $YELLOW "        清理Docker资源"
    print_message $YELLOW "========================================\n"
    
    # 检查Docker环境
    if ! check_docker; then
        return 1
    fi
    
    # 停止所有相关容器
    print_message $BLUE "停止所有汇农微服务容器..."
    for service in "${!DOCKER_CONFIG[@]}"; do
        if docker ps | grep -q "$service"; then
            docker stop "$service" >/dev/null 2>&1
            print_message $GREEN "✓ 停止容器: $service"
        fi
    done
    
    # 删除所有相关容器
    print_message $BLUE "删除所有汇农微服务容器..."
    for service in "${!DOCKER_CONFIG[@]}"; do
        if docker ps -a | grep -q "$service"; then
            docker rm "$service" >/dev/null 2>&1
            print_message $GREEN "✓ 删除容器: $service"
        fi
    done
    
    # 删除Docker网络
    if docker network ls | grep -q "$DOCKER_NETWORK"; then
        print_message $BLUE "删除Docker网络: $DOCKER_NETWORK"
        docker network rm "$DOCKER_NETWORK" >/dev/null 2>&1
        if [[ $? -eq 0 ]]; then
            print_message $GREEN "✓ Docker网络删除成功"
        else
            print_message $YELLOW "⚠ Docker网络删除失败（可能仍有容器在使用）"
        fi
    fi
    
    # 清理未使用的Docker资源
    print_message $BLUE "清理未使用的Docker资源..."
    docker system prune -f >/dev/null 2>&1
    print_message $GREEN "✓ Docker资源清理完成"
    
    print_message $YELLOW "\n========================================"
    print_message $YELLOW "        Docker资源清理完成"
    print_message $YELLOW "========================================\n"
}

# Docker镜像诊断
diagnose_docker_image() {
    local service=$1
    
    # 检查Docker环境
    if ! check_docker; then
        return 1
    fi
    
    # 检查服务配置是否存在
    local docker_info="${DOCKER_CONFIG[$service]}"
    if [[ -z "$docker_info" ]]; then
        print_message $RED "错误: 服务 $service 的Docker配置信息未找到"
        return 1
    fi
    
    local image_and_tag="${docker_info%:*}"
    
    print_message $BLUE "\n========================================"
    print_message $BLUE "    Docker 镜像诊断: $service"
    print_message $BLUE "========================================\n"
    
    print_message $YELLOW "镜像信息: $image_and_tag"
    
    # 检查镜像是否存在
    if ! docker images | grep -q "${image_and_tag%%:*}"; then
        print_message $RED "✗ 镜像不存在，请先拉取镜像:"
        print_message $YELLOW "  docker pull $image_and_tag"
        return 1
    fi
    
    print_message $GREEN "✓ 镜像存在"
    
    # 创建临时容器查看文件结构
    print_message $BLUE "检查镜像内部文件结构..."
    local temp_container="${service}-temp-$(date +%s)"
    
    # 启动临时容器
    docker run --name "$temp_container" --entrypoint="" "$image_and_tag" ls -la /app >/tmp/docker_ls.log 2>&1
    
    if [[ $? -eq 0 ]]; then
        print_message $GREEN "✓ /app 目录内容:"
        cat /tmp/docker_ls.log | sed 's/^/    /'
    else
        print_message $RED "✗ 无法访问 /app 目录"
    fi
    
    # 检查根目录
    docker run --name "${temp_container}-root" --entrypoint="" "$image_and_tag" ls -la / >/tmp/docker_root.log 2>&1
    if [[ $? -eq 0 ]]; then
        print_message $BLUE "根目录内容:"
        cat /tmp/docker_root.log | sed 's/^/    /'
    fi
    
    # 检查是否有可执行文件
    print_message $BLUE "搜索可执行文件..."
    docker run --name "${temp_container}-find" --entrypoint="" "$image_and_tag" find / -name "*$service*" -type f 2>/dev/null | head -10 | sed 's/^/    /'
    
    # 清理临时容器
    docker rm "$temp_container" "${temp_container}-root" "${temp_container}-find" >/dev/null 2>&1
    rm -f /tmp/docker_ls.log /tmp/docker_root.log
    
    print_message $BLUE "\n========================================"
    print_message $BLUE "        诊断完成"
    print_message $BLUE "========================================\n"
}

# 诊断所有Docker镜像
diagnose_all_docker_images() {
    print_message $GREEN "\n========================================"
    print_message $GREEN "        诊断所有Docker镜像"
    print_message $GREEN "========================================\n"
    
    for service in "${!DOCKER_CONFIG[@]}"; do
        diagnose_docker_image "$service"
        sleep 1
    done
}

# ============ 原有函数 ============

# 启动所有服务
start_all() {
    print_message $GREEN "\n========================================"
    print_message $GREEN "        开始启动汇农金融微服务"
    print_message $GREEN "========================================\n"
    
    create_dirs
    clean_logs
    
    # 第一阶段：启动RPC服务
    print_message $YELLOW "阶段1: 启动用户管理RPC服务"
    start_group "用户管理RPC服务" "${FIRST_RPC_GROUP[@]}"
    
    print_message $YELLOW "阶段2: 启动产品管理RPC服务"
    start_group "产品管理RPC服务" "${SECOND_RPC_GROUP[@]}"
    
    print_message $YELLOW "阶段3: 启动业务RPC服务"
    start_group "业务RPC服务" "${THIRD_RPC_GROUP[@]}"
    
    # 第二阶段：启动API服务
    print_message $YELLOW "阶段4: 启动用户管理API服务"
    start_group "用户管理API服务" "${FIRST_API_GROUP[@]}"
    
    print_message $YELLOW "阶段5: 启动产品管理API服务"
    start_group "产品管理API服务" "${SECOND_API_GROUP[@]}"
    
    print_message $YELLOW "阶段6: 启动业务API服务"
    start_group "业务API服务" "${THIRD_API_GROUP[@]}"
    
    print_message $GREEN "\n========================================"
    print_message $GREEN "          所有服务启动完成"
    print_message $GREEN "========================================\n"
    
    # 显示最终状态
    show_status
}

# 停止所有服务
stop_all() {
    print_message $RED "\n========================================"
    print_message $RED "        开始停止汇农金融微服务"
    print_message $RED "========================================\n"
    
    # 按相反顺序停止服务 - 先停API，再停RPC
    local all_services=()
    all_services+=("${THIRD_API_GROUP[@]}")
    all_services+=("${SECOND_API_GROUP[@]}")
    all_services+=("${FIRST_API_GROUP[@]}")
    all_services+=("${THIRD_RPC_GROUP[@]}")
    all_services+=("${SECOND_RPC_GROUP[@]}")
    all_services+=("${FIRST_RPC_GROUP[@]}")
    
    for service in "${all_services[@]}"; do
        stop_service "$service"
        sleep 1
    done
    
    print_message $RED "\n========================================"
    print_message $RED "          所有服务已停止"
    print_message $RED "========================================\n"
}

# 显示所有服务状态
show_status() {
    print_message $BLUE "\n========================================"
    print_message $BLUE "          服务状态"
    print_message $BLUE "========================================\n"
    
    print_message $BLUE "RPC服务:"
    for service in "${FIRST_RPC_GROUP[@]}" "${SECOND_RPC_GROUP[@]}" "${THIRD_RPC_GROUP[@]}"; do
        echo "  $(check_status "$service")"
    done
    
    print_message $BLUE "\nAPI服务:"
    for service in "${FIRST_API_GROUP[@]}" "${SECOND_API_GROUP[@]}" "${THIRD_API_GROUP[@]}"; do
        echo "  $(check_status "$service")"
    done
    echo
}

# 重启所有服务
restart_all() {
    print_message $YELLOW "\n========================================"
    print_message $YELLOW "        重启汇农金融微服务"
    print_message $YELLOW "========================================\n"
    
    stop_all
    sleep 5
    start_all
}

# 列出所有可用服务
list_services() {
    echo "可用的服务："
    echo ""
    echo "RPC服务:"
    for service in "${FIRST_RPC_GROUP[@]}" "${SECOND_RPC_GROUP[@]}" "${THIRD_RPC_GROUP[@]}"; do
        echo "  $service"
    done
    echo ""
    echo "API服务:"
    for service in "${FIRST_API_GROUP[@]}" "${SECOND_API_GROUP[@]}" "${THIRD_API_GROUP[@]}"; do
        echo "  $service"
    done
}

# 显示帮助信息
show_help() {
    echo "用法: $0 [--docker] {start|stop|restart|status|list|help} [service-name]"
    echo ""
    echo "运行模式:"
    echo "  --docker      - 使用Docker模式运行（默认为本地模式）"
    echo ""
    echo "全局命令:"
    echo "  start         - 启动所有微服务（先启动RPC，再启动API）"
    echo "  stop          - 停止所有微服务"
    echo "  restart       - 重启所有微服务"
    echo "  status        - 查看所有微服务状态"
    echo "  list          - 列出所有可用的服务"
    echo "  cleanup       - 清理Docker资源（仅Docker模式）"
    echo "  diagnose      - 诊断Docker镜像（仅Docker模式）"
    echo "  help          - 显示此帮助信息"
    echo ""
    echo "单服务命令:"
    echo "  start <服务名>     - 启动指定的微服务"
    echo "  stop <服务名>      - 停止指定的微服务"
    echo "  restart <服务名>   - 重启指定的微服务"
    echo "  status <服务名>    - 查看指定微服务状态"
    echo "  diagnose <服务名>  - 诊断指定Docker镜像（仅Docker模式）"
    echo ""
    echo "启动顺序（全部启动时）:"
    echo "  1. 用户管理RPC服务: appuser-rpc, oauser-rpc"
    echo "  2. 产品管理RPC服务: loanproduct-rpc, leaseproduct-rpc"
    echo "  3. 业务RPC服务: loan-rpc, lease-rpc"
    echo "  4. 用户管理API服务: appuser-api, oauser-api"
    echo "  5. 产品管理API服务: loanproduct-api, leaseproduct-api"
    echo "  6. 业务API服务: loan-api, lease-api"
    echo ""
    echo "示例:"
    echo "  $0 start                      # 本地模式启动所有服务"
    echo "  $0 --docker start             # Docker模式启动所有服务"
    echo "  $0 start appuser-rpc          # 本地模式仅启动 appuser-rpc 服务"
    echo "  $0 --docker start appuser-rpc # Docker模式仅启动 appuser-rpc 服务"
    echo "  $0 stop loan-api              # 本地模式仅停止 loan-api 服务"
    echo "  $0 --docker stop loan-api     # Docker模式仅停止 loan-api 服务"
    echo "  $0 status                     # 本地模式查看所有服务状态"
    echo "  $0 --docker status            # Docker模式查看所有服务状态"
    echo "  $0 --docker cleanup           # 清理所有Docker资源"
    echo "  $0 --docker diagnose          # 诊断所有Docker镜像"
    echo "  $0 --docker diagnose oauser-api # 诊断指定Docker镜像"
    echo ""
    echo "本地模式文件位置:"
    echo "  日志位置: $LOG_DIR (脚本日志)"
    echo "  服务日志位置: 各服务工作目录下的 logs/ 文件夹"
    echo "  PID文件位置: $PID_DIR"
    echo ""
    echo "Docker模式配置:"
    echo "  网络名称: $DOCKER_NETWORK"
    echo "  容器日志: docker logs <服务名>"
    echo "  镜像仓库: registry.huinong.internal/huinong/"
}

# 验证服务名是否有效
validate_service() {
    local service=$1
    if [[ -z "${SERVICE_CONFIG[$service]}" ]]; then
        print_message $RED "错误: 无效的服务名 '$service'"
        echo ""
        list_services
        return 1
    fi
    return 0
}

# 主函数
main() {
    # 解析参数
    local docker_mode=false
    local args=()
    
    # 处理 --docker 参数
    for arg in "$@"; do
        case "$arg" in
            --docker)
                docker_mode=true
                STARTUP_MODE="docker"
                ;;
            *)
                args+=("$arg")
                ;;
        esac
    done
    
    local command="${args[0]}"
    local service="${args[1]}"
    
    # 显示当前运行模式
    if [[ "$docker_mode" == true ]]; then
        print_message $YELLOW "运行模式: Docker"
    else
        print_message $YELLOW "运行模式: 本地"
    fi
    
    case "$command" in
        start)
            if [[ -n "$service" ]]; then
                # 启动单个服务
                if validate_service "$service"; then
                    if [[ "$docker_mode" == true ]]; then
                        create_docker_network
                        clean_service_logs "$service"
                        print_message $BLUE "启动单个Docker服务: $service"
                        start_docker_service "$service"
                    else
                        create_dirs
                        clean_service_logs "$service"
                        print_message $BLUE "启动单个服务: $service"
                        start_service "$service"
                    fi
                fi
            else
                # 启动所有服务
                if [[ "$docker_mode" == true ]]; then
                    start_all_docker
                else
                    start_all
                fi
            fi
            ;;
        stop)
            if [[ -n "$service" ]]; then
                # 停止单个服务
                if validate_service "$service"; then
                    if [[ "$docker_mode" == true ]]; then
                        print_message $BLUE "停止单个Docker服务: $service"
                        stop_docker_service "$service"
                    else
                        print_message $BLUE "停止单个服务: $service"
                        stop_service "$service"
                    fi
                fi
            else
                # 停止所有服务
                if [[ "$docker_mode" == true ]]; then
                    stop_all_docker
                else
                    stop_all
                fi
            fi
            ;;
        restart)
            if [[ -n "$service" ]]; then
                # 重启单个服务
                if validate_service "$service"; then
                    if [[ "$docker_mode" == true ]]; then
                        print_message $BLUE "重启单个Docker服务: $service"
                        stop_docker_service "$service"
                        sleep 3
                        create_docker_network
                        clean_service_logs "$service"
                        start_docker_service "$service"
                    else
                        print_message $BLUE "重启单个服务: $service"
                        stop_service "$service"
                        sleep 3
                        create_dirs
                        clean_service_logs "$service"
                        start_service "$service"
                    fi
                fi
            else
                # 重启所有服务
                if [[ "$docker_mode" == true ]]; then
                    restart_all_docker
                else
                    restart_all
                fi
            fi
            ;;
        status)
            if [[ -n "$service" ]]; then
                # 查看单个服务状态
                if validate_service "$service"; then
                    if [[ "$docker_mode" == true ]]; then
                        print_message $BLUE "Docker服务状态: $service"
                        check_docker_status "$service"
                    else
                        print_message $BLUE "服务状态: $service"
                        check_status "$service"
                    fi
                fi
            else
                # 查看所有服务状态
                if [[ "$docker_mode" == true ]]; then
                    show_docker_status
                else
                    show_status
                fi
            fi
            ;;
        list)
            list_services
            ;;
        cleanup)
            if [[ "$docker_mode" == true ]]; then
                cleanup_docker
            else
                print_message $RED "错误: cleanup命令仅在Docker模式下可用"
                print_message $YELLOW "请使用: $0 --docker cleanup"
                exit 1
            fi
            ;;
        diagnose)
            if [[ "$docker_mode" == true ]]; then
                if [[ -n "$service" ]]; then
                    # 诊断单个服务
                    if validate_service "$service"; then
                        diagnose_docker_image "$service"
                    fi
                else
                    # 诊断所有服务
                    diagnose_all_docker_images
                fi
            else
                print_message $RED "错误: diagnose命令仅在Docker模式下可用"
                print_message $YELLOW "请使用: $0 --docker diagnose [服务名]"
                exit 1
            fi
            ;;
        help|--help|-h)
            show_help
            ;;
        "")
            print_message $RED "错误: 缺少命令参数"
            echo ""
            show_help
            exit 1
            ;;
        *)
            print_message $RED "错误: 无效的命令 '$command'"
            echo ""
            show_help
            exit 1
            ;;
    esac
}

# 执行主函数
main "$@"
