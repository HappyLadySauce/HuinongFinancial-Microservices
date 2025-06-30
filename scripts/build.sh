#!/bin/bash

# 微服务构建脚本 - 基于go-zero微服务架构，集成SkyWalking监控
# 使用方法：./scripts/build.sh [service_name] [type] [options]
# 示例：./scripts/build.sh appuser api
#       ./scripts/build.sh appuser rpc
#       ./scripts/build.sh all api --release

set -e

# 颜色定义
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
PURPLE='\033[0;35m'
CYAN='\033[0;36m'
NC='\033[0m' # No Color

# SkyWalking配置
SKYWALKING_AGENT_PATH="app/skywalking-go-agent/skywalking-go-agent-0.6.0-linux-amd64"
SKYWALKING_IMPORT='_ "github.com/apache/skywalking-go"'

# 帮助信息
show_help() {
    echo -e "${BLUE}微服务构建脚本 - HuinongFinancial微服务 + SkyWalking监控${NC}"
    echo ""
    echo "功能特性:"
    echo "  - 支持构建go-zero微服务应用"
    echo "  - 自动集成SkyWalking Go Agent监控"
    echo "  - 支持Release和Debug两种构建模式"
    echo "  - 自动代码注入SkyWalking依赖"
    echo "  - 支持单服务和批量构建"
    echo "  - 构建产物统一输出到bin目录"
    echo ""
    echo "项目结构:"
    echo "  app/service/cmd/api/       # API服务源代码"
    echo "  app/service/cmd/rpc/       # RPC服务源代码"
    echo "  bin/                       # 构建产物输出目录"
    echo "  └── service-api            # API服务可执行文件"
    echo "  └── service-rpc            # RPC服务可执行文件"
    echo ""
    echo "使用方法:"
    echo "  $0 [service_name] [type] [options]"
    echo ""
    echo "参数说明:"
    echo "  service_name: 服务名称 (appuser|oauser|loan|loanproduct|leaseproduct|lease)"
    echo "  type:         构建类型 (api|rpc|all|clean|check)"
    echo "  options:      构建选项 (--release|--debug|--force|--no-skywalking)"
    echo ""
    echo "示例:"
    echo "  $0 appuser api              # 构建appuser的API服务"
    echo "  $0 appuser rpc              # 构建appuser的RPC服务"
    echo "  $0 appuser all              # 构建appuser的API和RPC服务"
    echo "  $0 all api                  # 构建所有服务的API"
    echo "  $0 all rpc                  # 构建所有服务的RPC"
    echo "  $0 all all                  # 构建所有服务的API和RPC"
    echo "  $0 appuser api --release    # Release模式构建appuser API"
    echo "  $0 appuser api --debug      # Debug模式构建appuser API"
    echo "  $0 appuser api --no-skywalking  # 不集成SkyWalking构建"
    echo "  $0 all clean                # 清理所有构建产物"
    echo "  $0 appuser check            # 检查appuser服务的构建环境"
    echo ""
    echo "构建选项:"
    echo "  --release        Release模式：优化编译，去除调试信息，缩小体积"
    echo "  --debug          Debug模式：包含调试信息，不优化（默认）"
    echo "  --force          强制重新构建，忽略文件时间戳"
    echo "  --no-skywalking  不集成SkyWalking监控"
    echo ""
    echo "支持的服务:"
    echo "  appuser      - App用户服务"
    echo "  oauser       - OA用户服务"
    echo "  loan         - 贷款服务"
    echo "  loanproduct  - 贷款产品服务"
    echo "  leaseproduct - 租赁产品服务"
    echo "  lease        - 租赁服务"
    echo ""
    echo "SkyWalking集成说明:"
    echo "  - Agent路径: $SKYWALKING_AGENT_PATH"
    echo "  - 自动注入导入: $SKYWALKING_IMPORT"
    echo "  - 编译参数: -toolexec=\"/path/to/agent\" -a"
    echo "  - 运行时通过环境变量配置服务名"
}

# 检查参数
if [[ $# -eq 0 ]] || [[ "$1" == "-h" ]] || [[ "$1" == "--help" ]]; then
    show_help
    exit 0
fi

SERVICE_NAME=$1
TYPE=$2

# 支持的服务列表
SERVICES=("appuser" "oauser" "loan" "loanproduct" "leaseproduct" "lease")
TYPES=("api" "rpc" "all" "clean" "check")

# 验证服务名称
if [[ "$SERVICE_NAME" != "all" ]] && [[ ! " ${SERVICES[@]} " =~ " ${SERVICE_NAME} " ]]; then
    echo -e "${RED}错误: 不支持的服务名称 '$SERVICE_NAME'${NC}"
    echo -e "支持的服务: ${SERVICES[*]}"
    exit 1
fi

# 验证类型
if [[ ! " ${TYPES[@]} " =~ " ${TYPE} " ]]; then
    echo -e "${RED}错误: 不支持的类型 '$TYPE'${NC}"
    echo -e "支持的类型: ${TYPES[*]}"
    exit 1
fi

# 解析构建选项
BUILD_MODE="debug"
FORCE_BUILD=false
ENABLE_SKYWALKING=true

shift 2  # 移除前两个参数
while [[ $# -gt 0 ]]; do
    case $1 in
        --release)
            BUILD_MODE="release"
            shift
            ;;
        --debug)
            BUILD_MODE="debug"
            shift
            ;;
        --force)
            FORCE_BUILD=true
            shift
            ;;
        --no-skywalking)
            ENABLE_SKYWALKING=false
            shift
            ;;
        *)
            echo -e "${RED}错误: 未知选项 '$1'${NC}"
            exit 1
            ;;
    esac
done

# 获取项目根目录
PROJECT_ROOT=$(cd "$(dirname "$0")/.." && pwd)
BIN_DIR="$PROJECT_ROOT/bin"
SKYWALKING_AGENT_FULL_PATH="$PROJECT_ROOT/$SKYWALKING_AGENT_PATH"

echo -e "${BLUE}========================================${NC}"
echo -e "${BLUE}  HuinongFinancial 微服务构建工具      ${NC}"
echo -e "${BLUE}========================================${NC}"
echo -e "${BLUE}项目根目录: $PROJECT_ROOT${NC}"
echo -e "${BLUE}构建模式: $BUILD_MODE${NC}"
echo -e "${BLUE}SkyWalking: $([ "$ENABLE_SKYWALKING" == true ] && echo "启用" || echo "禁用")${NC}"
echo -e "${BLUE}输出目录: $BIN_DIR${NC}"

# 检查Go环境
check_go_env() {
    if ! command -v go &> /dev/null; then
        echo -e "${RED}错误: Go 未安装或不在PATH中${NC}"
        exit 1
    fi
    
    GO_VERSION=$(go version | grep -oE 'go[0-9]+\.[0-9]+\.[0-9]+' | sed 's/go//')
    echo -e "${GREEN}Go版本: $GO_VERSION${NC}"
    
    # 检查Go版本是否满足要求（至少1.18）
    GO_MAJOR=$(echo "$GO_VERSION" | cut -d. -f1)
    GO_MINOR=$(echo "$GO_VERSION" | cut -d. -f2)
    
    if [[ $GO_MAJOR -lt 1 ]] || [[ $GO_MAJOR -eq 1 && $GO_MINOR -lt 18 ]]; then
        echo -e "${YELLOW}警告: 建议使用Go 1.18+版本以获得最佳兼容性${NC}"
    fi
}

# 检查SkyWalking Agent
check_skywalking_agent() {
    if [[ "$ENABLE_SKYWALKING" != true ]]; then
        echo -e "${YELLOW}SkyWalking监控已禁用${NC}"
        return 0
    fi
    
    if [[ ! -f "$SKYWALKING_AGENT_FULL_PATH" ]]; then
        echo -e "${RED}错误: SkyWalking Agent未找到${NC}"
        echo -e "${YELLOW}预期路径: $SKYWALKING_AGENT_FULL_PATH${NC}"
        echo -e "${YELLOW}请确保已下载并放置SkyWalking Go Agent${NC}"
        echo -e "${YELLOW}下载地址: https://skywalking.apache.org/downloads/${NC}"
        exit 1
    fi
    
    if [[ ! -x "$SKYWALKING_AGENT_FULL_PATH" ]]; then
        echo -e "${YELLOW}设置SkyWalking Agent执行权限...${NC}"
        chmod +x "$SKYWALKING_AGENT_FULL_PATH"
    fi
    
    echo -e "${GREEN}✓ SkyWalking Agent检查通过: $SKYWALKING_AGENT_FULL_PATH${NC}"
}

# 检查服务目录结构
check_service_structure() {
    local service=$1
    local code_type=$2
    local service_dir="$PROJECT_ROOT/app/$service"
    local cmd_dir="$service_dir/cmd/$code_type"
    
    if [[ ! -d "$service_dir" ]]; then
        echo -e "${RED}错误: 服务目录不存在 - $service_dir${NC}"
        return 1
    fi
    
    if [[ ! -d "$cmd_dir" ]]; then
        echo -e "${RED}错误: $service $code_type 目录不存在 - $cmd_dir${NC}"
        echo -e "${YELLOW}提示: 请先运行 ./scripts/gen-code.sh $service $code_type 生成代码${NC}"
        return 1
    fi
    
    # 查找main.go文件
    local main_file=""
    if [[ -f "$cmd_dir/$service.go" ]]; then
        main_file="$cmd_dir/$service.go"
    elif [[ -f "$cmd_dir/main.go" ]]; then
        main_file="$cmd_dir/main.go"
    elif [[ -f "$cmd_dir/$service$code_type.go" ]]; then
        # 处理 appuserrpc.go, appuserapi.go 这种命名模式
        main_file="$cmd_dir/$service$code_type.go"
    else
        # 查找任何.go文件作为入口点
        local go_files=($(find "$cmd_dir" -maxdepth 1 -name "*.go" -type f | head -1))
        if [[ ${#go_files[@]} -gt 0 ]]; then
            main_file="${go_files[0]}"
            echo -e "${YELLOW}警告: 使用找到的Go文件作为入口点: $(basename "$main_file")${NC}"
        else
            echo -e "${RED}错误: 未找到任何Go入口文件 - $cmd_dir${NC}"
            return 1
        fi
    fi
    
    echo "$main_file"
    return 0
}

# 注入SkyWalking依赖
inject_skywalking_import() {
    local main_file=$1
    local service=$2
    local code_type=$3
    local cmd_dir=$(dirname "$main_file")
    
    if [[ "$ENABLE_SKYWALKING" != true ]]; then
        return 0
    fi
    
    # 检查是否已经注入
    if grep -q "skywalking-go" "$main_file"; then
        echo -e "${BLUE}SkyWalking代码依赖已存在，跳过注入${NC}"
    else
        echo -e "${YELLOW}注入SkyWalking依赖到 $main_file${NC}"
        
        # 创建备份
        cp "$main_file" "$main_file.bak"
        
        # 查找import块，如果存在则在其中添加，否则在package声明后添加
        if grep -q "^import (" "$main_file"; then
            # 在import块中添加
            sed -i.tmp '/^import (/a\
\t'"$SKYWALKING_IMPORT" "$main_file"
        elif grep -q "^import " "$main_file"; then
            # 在单行import后添加
            sed -i.tmp '/^import /a\
import '"$SKYWALKING_IMPORT" "$main_file"
        else
            # 在package声明后添加
            sed -i.tmp '/^package /a\
\
import '"$SKYWALKING_IMPORT" "$main_file"
        fi
        
        # 清理临时文件
        rm -f "$main_file.tmp"
        
        echo -e "${GREEN}✓ SkyWalking代码依赖注入完成${NC}"
    fi
    
    # 检查并添加go.mod依赖
    local go_mod_file="$cmd_dir/go.mod"
    if [[ -f "$go_mod_file" ]]; then
        if ! grep -q "skywalking-go" "$go_mod_file"; then
            echo -e "${YELLOW}添加SkyWalking模块依赖到 go.mod${NC}"
            
            # 创建go.mod备份
            cp "$go_mod_file" "$go_mod_file.skywalking_bak"
            
            # 添加依赖到go.mod的require块中
            if grep -q "require (" "$go_mod_file"; then
                # 在require块中添加
                sed -i.tmp '/require (/a\
\tgithub.com/apache/skywalking-go v0.6.0' "$go_mod_file"
            else
                # 在文件末尾添加require块
                echo "" >> "$go_mod_file"
                echo "require (" >> "$go_mod_file"
                echo -e "\tgithub.com/apache/skywalking-go v0.6.0" >> "$go_mod_file"
                echo ")" >> "$go_mod_file"
            fi
            
            # 清理临时文件
            rm -f "$go_mod_file.tmp"
            
            echo -e "${GREEN}✓ SkyWalking模块依赖添加完成${NC}"
            
            # 执行go mod tidy来下载依赖
            echo -e "${BLUE}下载SkyWalking依赖包...${NC}"
            if go mod tidy 2>/dev/null; then
                echo -e "${GREEN}✓ 依赖包下载完成${NC}"
            else
                echo -e "${YELLOW}⚠ 依赖包下载出现警告，但不影响构建${NC}"
            fi
        else
            echo -e "${BLUE}SkyWalking模块依赖已存在${NC}"
        fi
    else
        echo -e "${YELLOW}⚠ 未找到go.mod文件，跳过模块依赖添加${NC}"
    fi
}

# 恢复SkyWalking依赖注入
restore_main_file() {
    local main_file=$1
    local cmd_dir=$(dirname "$main_file")
    local go_mod_file="$cmd_dir/go.mod"
    
    if [[ -f "$main_file.bak" ]]; then
        mv "$main_file.bak" "$main_file"
        echo -e "${BLUE}恢复main文件: $main_file${NC}"
    fi
    
    # 如果存在go.mod备份，也要恢复
    if [[ -f "$go_mod_file.skywalking_bak" ]]; then
        mv "$go_mod_file.skywalking_bak" "$go_mod_file"
        echo -e "${BLUE}恢复go.mod文件: $go_mod_file${NC}"
    fi
}

# 构建服务
build_service() {
    local service=$1
    local code_type=$2
    
    echo -e "${GREEN}构建 $service $code_type 服务...${NC}"
    
    # 检查服务结构并获取main文件路径
    local main_file
    if ! main_file=$(check_service_structure "$service" "$code_type"); then
        return 1
    fi
    
    local cmd_dir=$(dirname "$main_file")
    local output_name="$service-$code_type"
    local output_path="$BIN_DIR/$output_name"
    
    echo -e "${BLUE}源目录: $cmd_dir${NC}"
    echo -e "${BLUE}主文件: $main_file${NC}"
    echo -e "${BLUE}输出文件: $output_path${NC}"
    
    # 检查是否需要重新构建
    if [[ "$FORCE_BUILD" != true ]] && [[ -f "$output_path" ]]; then
        if [[ "$output_path" -nt "$cmd_dir" ]]; then
            echo -e "${YELLOW}$output_name 已是最新，跳过构建（使用 --force 强制重建）${NC}"
            return 0
        fi
    fi
    
    # 进入源代码目录
    cd "$cmd_dir"
    
    # 注入SkyWalking依赖
    inject_skywalking_import "$main_file" "$service" "$code_type"
    
    # 构建Go编译参数
    local build_args=()
    
    # SkyWalking编译参数
    if [[ "$ENABLE_SKYWALKING" == true ]]; then
        build_args+=("-toolexec=$SKYWALKING_AGENT_FULL_PATH")
        build_args+=("-a")  # 强制重新构建所有包
    fi
    
    # 构建模式参数
    case $BUILD_MODE in
        "release")
            build_args+=("-ldflags=-s -w")  # 去除调试信息和符号表
            build_args+=("-trimpath")       # 去除文件路径信息
            ;;
        "debug")
            build_args+=("-gcflags=all=-N -l")  # 禁用优化和内联
            ;;
    esac
    
    # 输出参数
    build_args+=("-o" "$output_path")
    
    # 执行构建
    echo -e "${BLUE}执行构建命令...${NC}"
    echo -e "${CYAN}go build ${build_args[*]} .${NC}"
    
    local build_start_time=$(date +%s)
    
    if go build "${build_args[@]}" .; then
        local build_end_time=$(date +%s)
        local build_duration=$((build_end_time - build_start_time))
        
        echo -e "${GREEN}✓ $service $code_type 构建成功${NC}"
        echo -e "${BLUE}构建时间: ${build_duration}秒${NC}"
        
        # 显示文件信息
        if [[ -f "$output_path" ]]; then
            local file_size=$(du -h "$output_path" | cut -f1)
            echo -e "${BLUE}文件大小: $file_size${NC}"
            echo -e "${BLUE}文件路径: $output_path${NC}"
        fi
        
        # 恢复main文件
        restore_main_file "$main_file"
        
        return 0
    else
        echo -e "${RED}✗ $service $code_type 构建失败${NC}"
        
        # 恢复main文件
        restore_main_file "$main_file"
        
        return 1
    fi
}

# 清理构建产物
clean_service() {
    local service=$1
    
    if [[ "$service" == "all" ]]; then
        echo -e "${YELLOW}清理所有构建产物...${NC}"
        if [[ -d "$BIN_DIR" ]]; then
            rm -rf "$BIN_DIR"/*
            echo -e "${GREEN}✓ 所有构建产物清理完成${NC}"
        else
            echo -e "${BLUE}构建产物目录不存在，无需清理${NC}"
        fi
    else
        echo -e "${YELLOW}清理 $service 服务构建产物...${NC}"
        local api_file="$BIN_DIR/$service-api"
        local rpc_file="$BIN_DIR/$service-rpc"
        local cleaned=false
        
        if [[ -f "$api_file" ]]; then
            rm -f "$api_file"
            echo -e "${GREEN}✓ 清理 $service-api${NC}"
            cleaned=true
        fi
        
        if [[ -f "$rpc_file" ]]; then
            rm -f "$rpc_file"
            echo -e "${GREEN}✓ 清理 $service-rpc${NC}"
            cleaned=true
        fi
        
        if [[ "$cleaned" != true ]]; then
            echo -e "${BLUE}$service 服务无构建产物需要清理${NC}"
        fi
    fi
}

# 检查服务环境
check_service() {
    local service=$1
    
    echo -e "${GREEN}检查 $service 服务构建环境...${NC}"
    
    local service_dir="$PROJECT_ROOT/app/$service"
    
    if [[ ! -d "$service_dir" ]]; then
        echo -e "${RED}✗ 服务目录不存在: $service_dir${NC}"
        return 1
    fi
    
    echo -e "${GREEN}✓ 服务目录存在${NC}"
    
    # 检查API服务
    local api_dir="$service_dir/cmd/api"
    if [[ -d "$api_dir" ]]; then
        echo -e "${GREEN}✓ API服务代码存在${NC}"
        
        # 检查go.mod
        if [[ -f "$api_dir/go.mod" ]]; then
            echo -e "${GREEN}✓ API go.mod存在${NC}"
        else
            echo -e "${YELLOW}⚠ API go.mod不存在${NC}"
        fi
        
        # 检查main文件
        if [[ -f "$api_dir/$service.go" ]] || [[ -f "$api_dir/main.go" ]] || [[ -f "$api_dir/${service}api.go" ]] || [[ -n "$(find "$api_dir" -maxdepth 1 -name "*.go" -type f 2>/dev/null | head -1)" ]]; then
            echo -e "${GREEN}✓ API main文件存在${NC}"
        else
            echo -e "${RED}✗ API main文件不存在${NC}"
        fi
    else
        echo -e "${YELLOW}⚠ API服务代码不存在${NC}"
        echo -e "${BLUE}  生成API代码: ./scripts/gen-code.sh $service api${NC}"
    fi
    
    # 检查RPC服务
    local rpc_dir="$service_dir/cmd/rpc"
    if [[ -d "$rpc_dir" ]]; then
        echo -e "${GREEN}✓ RPC服务代码存在${NC}"
        
        # 检查go.mod
        if [[ -f "$rpc_dir/go.mod" ]]; then
            echo -e "${GREEN}✓ RPC go.mod存在${NC}"
        else
            echo -e "${YELLOW}⚠ RPC go.mod不存在${NC}"
        fi
        
        # 检查main文件
        if [[ -f "$rpc_dir/$service.go" ]] || [[ -f "$rpc_dir/main.go" ]] || [[ -f "$rpc_dir/${service}rpc.go" ]] || [[ -n "$(find "$rpc_dir" -maxdepth 1 -name "*.go" -type f 2>/dev/null | head -1)" ]]; then
            echo -e "${GREEN}✓ RPC main文件存在${NC}"
        else
            echo -e "${RED}✗ RPC main文件不存在${NC}"
        fi
    else
        echo -e "${YELLOW}⚠ RPC服务代码不存在${NC}"
        echo -e "${BLUE}  生成RPC代码: ./scripts/gen-code.sh $service rpc${NC}"
    fi
    
    echo -e "${GREEN}✓ $service 服务环境检查完成${NC}"
}

# 构建单个服务的指定类型
build_single_service() {
    local service=$1
    local type=$2
    
    case $type in
        "api")
            build_service "$service" "api"
            ;;
        "rpc")
            build_service "$service" "rpc"
            ;;
        "all")
            local success=true
            build_service "$service" "api" || success=false
            build_service "$service" "rpc" || success=false
            [[ "$success" == true ]]
            ;;
        "clean")
            clean_service "$service"
            ;;
        "check")
            check_service "$service"
            ;;
    esac
}

# 构建所有服务
build_all_services() {
    local type=$1
    
    if [[ "$type" == "clean" ]]; then
        clean_service "all"
        return $?
    fi
    
    echo -e "${BLUE}开始构建所有服务...${NC}"
    
    local total_success=true
    local success_count=0
    local total_count=0
    
    for service in "${SERVICES[@]}"; do
        echo ""
        echo -e "${PURPLE}========== 构建 $service 服务 ==========${NC}"
        
        case $type in
            "api"|"rpc")
                ((total_count++))
                if build_single_service "$service" "$type"; then
                    ((success_count++))
                else
                    total_success=false
                fi
                ;;
            "all")
                ((total_count += 2))  # api + rpc
                local service_success=true
                build_service "$service" "api" || service_success=false
                build_service "$service" "rpc" || service_success=false
                
                if [[ "$service_success" == true ]]; then
                    ((success_count += 2))
                else
                    total_success=false
                    # 计算实际成功的数量
                    [[ -f "$BIN_DIR/$service-api" ]] && ((success_count++)) || true
                    [[ -f "$BIN_DIR/$service-rpc" ]] && ((success_count++)) || true
                fi
                ;;
            "check")
                check_service "$service"
                ;;
        esac
    done
    
    if [[ "$type" != "check" ]]; then
        echo ""
        echo -e "${BLUE}========== 构建汇总 ==========${NC}"
        echo -e "${BLUE}成功: $success_count/$total_count${NC}"
        
        if [[ "$total_success" == true ]]; then
            echo -e "${GREEN}✓ 所有服务构建成功！${NC}"
        else
            echo -e "${YELLOW}⚠ 部分服务构建失败${NC}"
        fi
    fi
    
    [[ "$total_success" == true ]]
}

# 显示构建结果
show_build_results() {
    if [[ ! -d "$BIN_DIR" ]] || [[ -z "$(ls -A "$BIN_DIR" 2>/dev/null)" ]]; then
        echo -e "${YELLOW}无构建产物${NC}"
        return
    fi
    
    echo ""
    echo -e "${BLUE}========== 构建产物 ==========${NC}"
    echo -e "${BLUE}输出目录: $BIN_DIR${NC}"
    
    for file in "$BIN_DIR"/*; do
        if [[ -f "$file" ]]; then
            local filename=$(basename "$file")
            local filesize=$(du -h "$file" | cut -f1)
            local filetime=$(stat -c %y "$file" 2>/dev/null || stat -f %Sm "$file" 2>/dev/null || echo "未知时间")
            echo -e "${GREEN}  $filename${NC} (${CYAN}$filesize${NC}, $filetime)"
        fi
    done
    
    echo ""
    echo -e "${BLUE}========== 运行说明 ==========${NC}"
    
    if [[ "$ENABLE_SKYWALKING" == true ]]; then
        echo -e "${YELLOW}SkyWalking监控已集成，运行时请设置环境变量:${NC}"
        echo -e "${CYAN}export SW_AGENT_NAME=服务名${NC}"
        echo -e "${CYAN}export SW_AGENT_SERVER=skywalking-oap地址${NC}"
        echo ""
        echo -e "${YELLOW}示例运行命令:${NC}"
        echo -e "${CYAN}SW_AGENT_NAME=appuser-api ./bin/appuser-api${NC}"
        echo -e "${CYAN}SW_AGENT_NAME=appuser-rpc ./bin/appuser-rpc${NC}"
    else
        echo -e "${YELLOW}直接运行构建的服务:${NC}"
        echo -e "${CYAN}./bin/appuser-api${NC}"
        echo -e "${CYAN}./bin/appuser-rpc${NC}"
    fi
}

# 主函数
main() {
    # 进入项目根目录
    cd "$PROJECT_ROOT"
    
    # 检查环境
    check_go_env
    check_skywalking_agent
    
    # 创建输出目录
    mkdir -p "$BIN_DIR"
    
    # 执行构建
    if [[ "$SERVICE_NAME" == "all" ]]; then
        build_all_services "$TYPE"
    else
        build_single_service "$SERVICE_NAME" "$TYPE"
    fi
    
    local exit_code=$?
    
    # 显示结果
    if [[ "$TYPE" != "clean" ]] && [[ "$TYPE" != "check" ]]; then
        show_build_results
    fi
    
    echo ""
    echo -e "${GREEN}========================================${NC}"
    if [[ $exit_code -eq 0 ]]; then
        echo -e "${GREEN}          构建完成！                   ${NC}"
    else
        echo -e "${RED}          构建失败！                   ${NC}"
    fi
    echo -e "${GREEN}========================================${NC}"
    
    exit $exit_code
}

# 执行主函数
main "$@"