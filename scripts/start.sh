#!/bin/bash

# 微服务启动停止脚本
# 作者: Auto Generated
# 描述: 管理汇农金融微服务的启动、停止和状态查看

# 配置区域
BASE_DIR="/root/HuinongFinancial-Microservices"
BIN_DIR="$BASE_DIR/bin"
PID_DIR="/tmp/huinong_pids"
LOG_DIR="/var/log/huinong"

# SkyWalking配置
SW_AGENT_SERVER=${SW_AGENT_SERVER:-"skywalking-oap:11800"}

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
    
    nohup env SW_AGENT_NAME="$service" SW_AGENT_SERVER="$SW_AGENT_SERVER" "$BIN_DIR/$service" -f "$config_file" > "$LOG_DIR/$service.log" 2>&1 &
    local pid=$!
    
    # 保存PID
    echo "$pid" > "$pid_file"
    
    # 等待服务启动
    sleep 3
    
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

# 启动所有服务
start_all() {
    print_message $GREEN "\n========================================"
    print_message $GREEN "        开始启动汇农金融微服务"
    print_message $GREEN "========================================\n"
    
    create_dirs
    
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

# 显示帮助信息
show_help() {
    echo "用法: $0 {start|stop|restart|status|help}"
    echo ""
    echo "命令说明:"
    echo "  start   - 启动所有微服务（先启动RPC，再启动API）"
    echo "  stop    - 停止所有微服务"
    echo "  restart - 重启所有微服务"
    echo "  status  - 查看所有微服务状态"
    echo "  help    - 显示此帮助信息"
    echo ""
    echo "启动顺序:"
    echo "  1. 用户管理RPC服务: appuser-rpc, oauser-rpc"
    echo "  2. 产品管理RPC服务: loanproduct-rpc, leaseproduct-rpc"
    echo "  3. 业务RPC服务: loan-rpc, lease-rpc"
    echo "  4. 用户管理API服务: appuser-api, oauser-api"
    echo "  5. 产品管理API服务: loanproduct-api, leaseproduct-api"
    echo "  6. 业务API服务: loan-api, lease-api"
    echo ""
    echo "环境变量:"
    echo "  SW_AGENT_SERVER - SkyWalking OAP服务器地址 (默认: skywalking-oap:11800)"
    echo ""
    echo "日志位置: $LOG_DIR"
    echo "PID文件位置: $PID_DIR"
}

# 主函数
main() {
    case "$1" in
        start)
            start_all
            ;;
        stop)
            stop_all
            ;;
        restart)
            restart_all
            ;;
        status)
            show_status
            ;;
        help|--help|-h)
            show_help
            ;;
        *)
            print_message $RED "错误: 无效的参数 '$1'"
            echo ""
            show_help
            exit 1
            ;;
    esac
}

# 执行主函数
main "$@"
