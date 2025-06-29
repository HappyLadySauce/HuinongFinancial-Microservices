#!/bin/bash

# 代码生成脚本 - 基于go-zero微服务架构
# 使用方法：./scripts/gen-code.sh [service_name] [type]
# 示例：./scripts/gen-code.sh appuser api
#       ./scripts/gen-code.sh appuser rpc

set -e

# 颜色定义
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# 帮助信息
show_help() {
    echo -e "${BLUE}代码生成脚本 - HuinongFinancial微服务${NC}"
    echo ""
    echo "项目结构设计:"
    echo "  app/service/               # 服务根目录"
    echo "  ├── service.api           # API定义文件"
    echo "  ├── service.proto         # RPC定义文件"
    echo "  ├── service.sql           # 数据库初始化文件"
    echo "  └── cmd/                  # 生成代码目录"
    echo "      ├── api/              # API服务代码"
    echo "      ├── rpc/              # RPC服务代码"
    echo "      └── model/            # 数据模型代码(API和RPC共用)"
    echo ""
    echo "使用方法:"
    echo "  $0 [service_name] [type]"
    echo ""
    echo "参数说明:"
    echo "  service_name: 服务名称 (appuser|oauser|loan|loanproduct|leaseproduct|lease)"
    echo "  type:         生成类型 (api|rpc|model|migrate|docs|all|clean|tidy|workspace)"
    echo ""
    echo "示例:"
    echo "  $0 appuser api           # 生成appuser服务的API代码"
    echo "  $0 appuser rpc           # 生成appuser服务的RPC代码"
    echo "  $0 loan model            # 生成loan服务的Model代码"
    echo "  $0 oauser migrate        # 执行oauser服务数据库迁移(先清空再迁移)"
    echo "  $0 appuser docs          # 生成appuser服务的API/RPC文档"
    echo "  $0 all api               # 生成所有服务的API代码"
    echo "  $0 all rpc               # 生成所有服务的RPC代码"
    echo "  $0 all model             # 生成所有服务的Model代码"
    echo "  $0 all migrate           # 执行所有服务数据库迁移(先清空再迁移)"
    echo "  $0 all docs              # 生成所有服务的API/RPC文档"
    echo "  $0 all all               # 生成所有服务的API、RPC、Model代码并迁移数据库"
    echo "  $0 all clean             # 清理所有服务生成的代码"
    echo "  $0 appuser all           # 生成appuser服务的所有代码"
    echo "  $0 appuser clean         # 清理appuser服务生成的代码"
    echo "  $0 appuser tidy          # 对appuser服务的API和RPC目录执行go mod tidy"
    echo "  $0 all tidy              # 对所有服务的API和RPC目录执行go mod tidy"
    echo "  $0 appuser workspace     # 设置appuser服务的Go工作区"
    echo "  $0 all workspace         # 设置所有服务的Go工作区"
    echo ""
    echo "支持的服务:"
    echo "  appuser      - App用户服务"
    echo "  oauser       - OA用户服务"
    echo "  loan         - 贷款服务"
    echo "  loanproduct  - 贷款产品服务"
    echo "  leaseproduct - 租赁产品服务"
    echo "  lease        - 租赁服务"
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
TYPES=("api" "rpc" "model" "migrate" "docs" "all" "clean" "tidy" "workspace")

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

# 获取项目根目录
PROJECT_ROOT=$(cd "$(dirname "$0")/.." && pwd)
echo -e "${BLUE}项目根目录: $PROJECT_ROOT${NC}"

# 检查goctl版本
check_goctl() {
    if ! command -v goctl &> /dev/null; then
        echo -e "${RED}错误: goctl 未安装${NC}"
        echo "请先安装goctl: go install github.com/zeromicro/go-zero/tools/goctl@latest"
        exit 1
    fi
    
    GOCTL_VERSION=$(goctl --version | grep -oE '[0-9]+\.[0-9]+\.[0-9]+' | head -1)
    echo -e "${GREEN}goctl版本: $GOCTL_VERSION${NC}"
}

# 清理生成的代码
clean_service() {
    local service=$1
    local cmd_dir="$PROJECT_ROOT/app/$service/cmd"
    local docs_dir="$PROJECT_ROOT/docs/$service"
    
    local cleaned=false
    
    if [[ -d "$cmd_dir" ]]; then
        echo -e "${YELLOW}清理 $service 服务生成的代码...${NC}"
        rm -rf "$cmd_dir"
        echo -e "${GREEN}✓ $service 代码清理完成${NC}"
        cleaned=true
    fi
    
    if [[ -d "$docs_dir" ]]; then
        echo -e "${YELLOW}清理 $service 服务生成的文档...${NC}"
        rm -rf "$docs_dir"
        echo -e "${GREEN}✓ $service 文档清理完成${NC}"
        cleaned=true
    fi
    
    if [[ "$cleaned" == false ]]; then
        echo -e "${YELLOW}$service 服务无生成的代码或文档需要清理${NC}"
    fi
}

# 执行go mod tidy
run_go_mod_tidy() {
    local target_dir=$1
    local service_name=$2
    local code_type=$3
    
    if [[ -d "$target_dir" ]]; then
        echo -e "${YELLOW}执行 $service_name $code_type go mod tidy...${NC}"
        cd "$target_dir"
        
        # 检查是否存在go.mod文件，如果不存在则初始化
        if [[ ! -f "go.mod" ]]; then
            echo -e "${BLUE}初始化 go module...${NC}"
            go mod init "$service_name-$code_type" 2>/dev/null || true
        fi
        
        # 执行go mod tidy
        if go mod tidy 2>/dev/null; then
            echo -e "${GREEN}✓ $service_name $code_type go mod tidy 完成${NC}"
        else
            echo -e "${YELLOW}⚠ $service_name $code_type go mod tidy 出现警告，但不影响使用${NC}"
        fi
    fi
}

# 设置Go工作区
setup_go_workspace() {
    local service=$1
    local cmd_dir="$PROJECT_ROOT/app/$service/cmd"
    
    echo -e "${GREEN}设置 $service 服务的Go工作区...${NC}"
    
    # 进入cmd目录
    cd "$cmd_dir"
    
    # 检查是否已存在go.work文件
    if [[ -f "go.work" ]]; then
        echo -e "${YELLOW}go.work 文件已存在，重新初始化...${NC}"
        rm -f "go.work"
    fi
    
    # 初始化工作区
    echo -e "${BLUE}初始化 Go 工作区...${NC}"
    go work init
    
    # 添加所有存在的模块到工作区
    local modules=()
    
    if [[ -d "api" ]] && [[ -f "api/go.mod" ]]; then
        modules+=("api")
    fi
    
    if [[ -d "rpc" ]] && [[ -f "rpc/go.mod" ]]; then
        modules+=("rpc")
    fi
    
    if [[ -d "model" ]] && [[ -f "model/go.mod" ]]; then
        modules+=("model")
    fi
    
    # 添加模块到工作区
    for module in "${modules[@]}"; do
        echo -e "${BLUE}添加 $module 模块到工作区...${NC}"
        go work use "$module"
    done
    
    if [[ ${#modules[@]} -gt 0 ]]; then
        echo -e "${GREEN}✓ $service Go工作区设置完成，包含模块: ${modules[*]}${NC}"
        
        # 显示工作区信息
        echo -e "${BLUE}工作区配置:${NC}"
        cat go.work | sed 's/^/    /'
        
        # 执行go work sync同步依赖
        echo -e "${BLUE}同步工作区依赖...${NC}"
        if go work sync 2>/dev/null; then
            echo -e "${GREEN}✓ 工作区依赖同步完成${NC}"
        else
            echo -e "${YELLOW}⚠ 工作区依赖同步出现警告，但不影响使用${NC}"
        fi
    else
        echo -e "${YELLOW}未找到任何可用模块，跳过工作区设置${NC}"
    fi
}

# 生成API代码
generate_api() {
    local service=$1
    local api_file="$PROJECT_ROOT/app/$service/$service-api.api"
    local api_dir="$PROJECT_ROOT/app/$service/cmd/api"
    
    if [[ ! -f "$api_file" ]]; then
        echo -e "${YELLOW}跳过 $service: API文件不存在 ($service-api.api)${NC}"
        return
    fi
    
    echo -e "${GREEN}生成 $service API代码...${NC}"
    echo -e "${BLUE}源文件: $api_file${NC}"
    echo -e "${BLUE}目标目录: $api_dir${NC}"
    
    # 创建目录
    mkdir -p "$api_dir"
    
    # 进入API目录
    cd "$api_dir"
    
    # 生成API代码
    if [[ $(echo "$GOCTL_VERSION" | cut -d. -f1) -ge 1 ]] && [[ $(echo "$GOCTL_VERSION" | cut -d. -f2) -ge 3 ]]; then
        # goctl >= 1.3
        goctl api go -api "$api_file" -dir . -style=goZero
    else
        # goctl < 1.3
        goctl api go -api "$api_file" -dir . -style=goZero
    fi
    
    echo -e "${GREEN}✓ $service API代码生成完成${NC}"
    
    # 执行go mod tidy
    run_go_mod_tidy "$api_dir" "$service" "api"
}

# 生成RPC代码
generate_rpc() {
    local service=$1
    local proto_file="$PROJECT_ROOT/app/$service/$service-rpc.proto"
    local rpc_dir="$PROJECT_ROOT/app/$service/cmd/rpc"

    
    if [[ ! -f "$proto_file" ]]; then
        echo -e "${YELLOW}跳过 $service: Proto文件不存在 ($service-rpc.proto)${NC}"
        return
    fi
    
    echo -e "${GREEN}生成 $service RPC代码...${NC}"
    echo -e "${BLUE}源文件: $proto_file${NC}"
    echo -e "${BLUE}目标目录: $rpc_dir${NC}"
    
    # 创建目录
    mkdir -p "$rpc_dir"
    
    # 生成RPC代码
    if [[ $(echo "$GOCTL_VERSION" | cut -d. -f1) -ge 1 ]] && [[ $(echo "$GOCTL_VERSION" | cut -d. -f2) -ge 3 ]]; then
        # goctl >= 1.3
        cd "$rpc_dir"
        goctl rpc protoc "$proto_file" --go_out=. --go-grpc_out=. --zrpc_out=. --proto_path="$(dirname "$proto_file")"
        
        # 移除omitempty标签
        if [[ -f "$rpc_dir/$service.pb.go" ]]; then
            sed -i.bak 's/,omitempty//g' "$rpc_dir/$service.pb.go" && rm -f "$rpc_dir/$service.pb.go.bak" 2>/dev/null || \
            sed -i 's/,omitempty//g' "$rpc_dir/$service.pb.go" 2>/dev/null || true
        fi
    else
        # goctl < 1.3
        cd "$rpc_dir"
        goctl rpc proto -src "$proto_file" -dir . -style=goZero
        
        # 移除omitempty标签 - goctl < 1.3 版本会自动创建pb目录并生成文件
        pb_files=$(find . -name "*.pb.go" 2>/dev/null || true)
        for pb_file in $pb_files; do
            if [[ -f "$pb_file" ]]; then
                sed -i.bak 's/,omitempty//g' "$pb_file" && rm -f "$pb_file.bak" 2>/dev/null || \
                sed -i 's/,omitempty//g' "$pb_file" 2>/dev/null || true
            fi
        done
    fi
    
    echo -e "${GREEN}✓ $service RPC代码生成完成${NC}"
    
    # 执行go mod tidy
    run_go_mod_tidy "$rpc_dir" "$service" "rpc"
}

# 生成Model代码
generate_model() {
    local service=$1
    local sql_file="$PROJECT_ROOT/app/$service/$service.sql"
    local model_dir="$PROJECT_ROOT/app/$service/cmd/model"
    
    if [[ ! -f "$sql_file" ]]; then
        echo -e "${YELLOW}跳过 $service: SQL文件不存在 ($service.sql)${NC}"
        return
    fi
    
    echo -e "${GREEN}生成 $service Model代码...${NC}"
    echo -e "${BLUE}源文件: $sql_file${NC}"
    echo -e "${BLUE}目标目录: $model_dir${NC}"
    
    # 创建目录
    mkdir -p "$model_dir"
    
    # 使用DDL方式从SQL文件生成Model代码
    cd "$model_dir"
    goctl model mysql ddl --src "$sql_file" --dir . --cache=true --style=goZero
    
    echo -e "${GREEN}✓ $service Model代码生成完成${NC}"
    echo -e "${YELLOW}注意: Model代码由API和RPC服务共用${NC}"
    
    # 执行go mod tidy
    run_go_mod_tidy "$model_dir" "$service" "model"
}

# 数据库迁移功能
migrate_database() {
    local service=$1
    local sql_file="$PROJECT_ROOT/app/$service/$service.sql"
    
    if [[ ! -f "$sql_file" ]]; then
        echo -e "${YELLOW}跳过 $service: SQL文件不存在 ($service.sql)${NC}"
        return
    fi
    
    # 数据库连接配置
    local db_host="10.10.10.6"
    local db_port="3306"
    local db_user=""
    local db_pass=""
    local db_name=""
    
    case $service in
        "appuser")
            db_user="appuser"
            db_pass="appuser"
            db_name="appuser"
            ;;
        "oauser")
            db_user="oauser"
            db_pass="oauser"
            db_name="oauser"
            ;;
        "loan")
            db_user="loan"
            db_pass="loan"
            db_name="loan"
            ;;
        "loanproduct")
            db_user="loanproduct"
            db_pass="loanproduct"
            db_name="loanproduct"
            ;;
        "leaseproduct")
            db_user="leaseproduct"
            db_pass="leaseproduct"
            db_name="leaseproduct"
            ;;
        "lease")
            db_user="lease"
            db_pass="lease"
            db_name="lease"
            ;;
        *)
            echo -e "${RED}错误: 未配置 $service 的数据库连接${NC}"
            return
            ;;
    esac
    
    echo -e "${GREEN}执行 $service 数据库迁移...${NC}"
    echo -e "${BLUE}数据库: $db_user@$db_host:$db_port/$db_name${NC}"
    echo -e "${BLUE}SQL文件: $sql_file${NC}"
    
    # 检查数据库连接
    if ! mysql -h "$db_host" -P "$db_port" -u "$db_user" -p"$db_pass" -e "SELECT 1;" 2>/dev/null; then
        echo -e "${RED}错误: 无法连接到数据库 $db_name${NC}"
        echo -e "${YELLOW}请检查数据库服务是否启动，用户权限是否正确${NC}"
        return
    fi
    
    # 清空数据库
    echo -e "${YELLOW}正在清空数据库 $db_name...${NC}"
    
    # 获取所有表名
    TABLES=$(mysql -h "$db_host" -P "$db_port" -u "$db_user" -p"$db_pass" "$db_name" -e "SHOW TABLES;" 2>/dev/null | grep -v "Tables_in_" | tr '\n' ' ')
    
    if [[ -n "$TABLES" ]]; then
        echo -e "${BLUE}发现表: $TABLES${NC}"
        
        # 禁用外键检查
        mysql -h "$db_host" -P "$db_port" -u "$db_user" -p"$db_pass" "$db_name" -e "SET FOREIGN_KEY_CHECKS = 0;" 2>/dev/null
        
        # 删除所有表
        for table in $TABLES; do
            echo -e "${YELLOW}删除表: $table${NC}"
            mysql -h "$db_host" -P "$db_port" -u "$db_user" -p"$db_pass" "$db_name" -e "DROP TABLE IF EXISTS \`$table\`;" 2>/dev/null
        done
        
        # 启用外键检查
        mysql -h "$db_host" -P "$db_port" -u "$db_user" -p"$db_pass" "$db_name" -e "SET FOREIGN_KEY_CHECKS = 1;" 2>/dev/null
        
        echo -e "${GREEN}✓ 数据库清空完成${NC}"
    else
        echo -e "${BLUE}数据库为空，无需清理${NC}"
    fi
    
    # 执行SQL迁移
    echo -e "${YELLOW}正在执行SQL迁移...${NC}"
    if mysql -h "$db_host" -P "$db_port" -u "$db_user" -p"$db_pass" "$db_name" < "$sql_file" 2>/dev/null; then
        echo -e "${GREEN}✓ $service 数据库迁移完成${NC}"
        
        # 显示创建的表
        echo -e "${BLUE}数据库表列表:${NC}"
        mysql -h "$db_host" -P "$db_port" -u "$db_user" -p"$db_pass" "$db_name" -e "SHOW TABLES;" 2>/dev/null | grep -v "Tables_in_"
    else
        echo -e "${RED}✗ $service 数据库迁移失败${NC}"
        echo -e "${YELLOW}请检查SQL文件语法是否正确${NC}"
    fi
}

# 生成文档
generate_docs() {
    local service=$1
    local api_file="$PROJECT_ROOT/app/$service/$service-api.api"
    local proto_file="$PROJECT_ROOT/app/$service/$service-rpc.proto"
    local docs_dir="$PROJECT_ROOT/docs/$service"
    
    echo -e "${GREEN}生成 $service 服务文档...${NC}"
    
    # 创建文档目录
    mkdir -p "$docs_dir"
    
    local has_docs=false
    
    # 生成API Swagger文档
    if [[ -f "$api_file" ]]; then
        echo -e "${BLUE}生成API Swagger文档...${NC}"
        echo -e "${BLUE}源文件: $api_file${NC}"
        echo -e "${BLUE}输出目录: $docs_dir${NC}"
        
        # 检查goctl版本是否支持swagger
        if goctl api swagger --help >/dev/null 2>&1; then
            # 生成JSON格式的swagger文档
            if goctl api swagger --api "$api_file" --dir "$docs_dir" 2>/dev/null; then
                echo -e "${GREEN}✓ Swagger JSON文档生成完成${NC}"
                
                # 生成YAML格式的swagger文档
                if goctl api swagger --api "$api_file" --dir "$docs_dir" --yaml 2>/dev/null; then
                    echo -e "${GREEN}✓ Swagger YAML文档生成完成${NC}"
                else
                    echo -e "${YELLOW}⚠ Swagger YAML文档生成失败，但JSON文档已生成${NC}"
                fi
                
                has_docs=true
            else
                echo -e "${YELLOW}⚠ Swagger文档生成失败，请检查API文件语法${NC}"
            fi
        else
            echo -e "${YELLOW}⚠ 当前goctl版本不支持swagger生成，需要goctl >= 1.8.2${NC}"
        fi
    else
        echo -e "${YELLOW}跳过API文档: $service-api.api 文件不存在${NC}"
    fi
    
    # 生成Proto文档
    if [[ -f "$proto_file" ]]; then
        echo -e "${BLUE}生成RPC Proto文档...${NC}"
        echo -e "${BLUE}源文件: $proto_file${NC}"
        
        # 检查是否安装了protoc-gen-doc
        if command -v protoc-gen-doc &> /dev/null; then
            # 生成HTML格式的protobuf文档
            if protoc --doc_out="$docs_dir" --doc_opt=html,"$service-rpc.html" --proto_path="$(dirname "$proto_file")" "$proto_file" 2>/dev/null; then
                echo -e "${GREEN}✓ Proto HTML文档生成完成${NC}"
                has_docs=true
            else
                echo -e "${YELLOW}⚠ Proto HTML文档生成失败${NC}"
            fi
            
            # 生成Markdown格式的protobuf文档
            if protoc --doc_out="$docs_dir" --doc_opt=markdown,"$service-rpc.md" --proto_path="$(dirname "$proto_file")" "$proto_file" 2>/dev/null; then
                echo -e "${GREEN}✓ Proto Markdown文档生成完成${NC}"
                has_docs=true
            else
                echo -e "${YELLOW}⚠ Proto Markdown文档生成失败${NC}"
            fi
        else
            # 创建简单的Proto说明文档
            cat > "$docs_dir/$service-rpc-proto.md" << EOF
# $service RPC 服务协议文档

## Protocol Buffer 定义文件
- 文件路径: \`app/$service/$service-rpc.proto\`
- 生成时间: $(date '+%Y-%m-%d %H:%M:%S')

## 说明
本文档由 \`$service-rpc.proto\` 文件自动生成。

要查看完整的服务定义，请直接查看 Protocol Buffer 文件：
\`\`\`bash
cat app/$service/$service-rpc.proto
\`\`\`

## 安装 protoc-gen-doc 生成更详细文档
要生成更详细的RPC文档，请安装 protoc-gen-doc：
\`\`\`bash
go install github.com/pseudomuto/protoc-gen-doc/cmd/protoc-gen-doc@latest
\`\`\`

然后重新运行文档生成命令。
EOF
            echo -e "${GREEN}✓ Proto说明文档生成完成${NC}"
            echo -e "${YELLOW}提示: 安装 protoc-gen-doc 可生成更详细的RPC文档${NC}"
            has_docs=true
        fi
    else
        echo -e "${YELLOW}跳过RPC文档: $service-rpc.proto 文件不存在${NC}"
    fi
    
    # 生成文档索引
    if [[ "$has_docs" == true ]]; then
        cat > "$docs_dir/README.md" << EOF
# $service 服务文档

本目录包含 $service 服务的API和RPC文档。

## 文档列表

### API文档 (Swagger)
EOF
        
        if [[ -f "$docs_dir/swagger.json" ]]; then
            echo "- [Swagger JSON](./swagger.json) - OpenAPI 3.0 JSON格式" >> "$docs_dir/README.md"
        fi
        
        if [[ -f "$docs_dir/swagger.yaml" ]]; then
            echo "- [Swagger YAML](./swagger.yaml) - OpenAPI 3.0 YAML格式" >> "$docs_dir/README.md"
        fi
        
        cat >> "$docs_dir/README.md" << EOF

### RPC文档 (Protocol Buffer)
EOF
        
        if [[ -f "$docs_dir/$service-rpc.html" ]]; then
            echo "- [$service-rpc.html](./$service-rpc.html) - HTML格式RPC文档" >> "$docs_dir/README.md"
        fi
        
        if [[ -f "$docs_dir/$service-rpc.md" ]]; then
            echo "- [$service-rpc.md](./$service-rpc.md) - Markdown格式RPC文档" >> "$docs_dir/README.md"
        fi
        
        if [[ -f "$docs_dir/$service-rpc-proto.md" ]]; then
            echo "- [$service-rpc-proto.md](./$service-rpc-proto.md) - Proto文件说明" >> "$docs_dir/README.md"
        fi
        
        cat >> "$docs_dir/README.md" << EOF

## 使用说明

### 查看Swagger文档
可以使用以下方式查看Swagger文档：
1. 将JSON/YAML文件导入到 [Swagger Editor](https://editor.swagger.io/)
2. 使用IDE插件（如VSCode的OpenAPI扩展）
3. 使用swagger-ui等工具

### 更新文档
当修改了API或RPC定义文件后，运行以下命令更新文档：
\`\`\`bash
./scripts/gen-code.sh $service docs
\`\`\`

生成时间: $(date '+%Y-%m-%d %H:%M:%S')
EOF
        
        echo -e "${GREEN}✓ 文档索引生成完成${NC}"
        echo -e "${BLUE}文档目录: $docs_dir${NC}"
        
        # 显示生成的文档列表
        echo -e "${BLUE}生成的文档文件:${NC}"
        ls -la "$docs_dir" | grep -v "^total" | tail -n +2 | sed 's/^/    /'
    else
        echo -e "${YELLOW}未生成任何文档，请检查API/RPC定义文件是否存在${NC}"
        # 删除空的文档目录
        rmdir "$docs_dir" 2>/dev/null || true
    fi
}

# 单独执行 go mod tidy
tidy_service() {
    local service=$1
    echo -e "${GREEN}对 $service 服务执行 go mod tidy...${NC}"
    
    local api_dir="$PROJECT_ROOT/app/$service/cmd/api"
    local rpc_dir="$PROJECT_ROOT/app/$service/cmd/rpc"
    local model_dir="$PROJECT_ROOT/app/$service/cmd/model"
    
    # 对API目录执行go mod tidy
    if [[ -d "$api_dir" ]]; then
        echo -e "${BLUE}处理API目录: $api_dir${NC}"
        run_go_mod_tidy "$api_dir" "$service" "api"
    else
        echo -e "${YELLOW}跳过 $service API: 目录不存在${NC}"
    fi
    
    # 对RPC目录执行go mod tidy
    if [[ -d "$rpc_dir" ]]; then
        echo -e "${BLUE}处理RPC目录: $rpc_dir${NC}"
        run_go_mod_tidy "$rpc_dir" "$service" "rpc"
    else
        echo -e "${YELLOW}跳过 $service RPC: 目录不存在${NC}"
    fi
    
    # 对Model目录执行go mod tidy
    if [[ -d "$model_dir" ]]; then
        echo -e "${BLUE}处理Model目录: $model_dir${NC}"
        run_go_mod_tidy "$model_dir" "$service" "model"
    else
        echo -e "${YELLOW}跳过 $service Model: 目录不存在${NC}"
    fi
    
    echo -e "${GREEN}✓ $service 服务 go mod tidy 完成${NC}"
    
    # 设置Go工作区
    setup_go_workspace $service
}

# 生成所有代码
generate_all() {
    local service=$1
    echo -e "${BLUE}生成 $service 服务的所有代码...${NC}"
    generate_api $service
    generate_rpc $service
    migrate_database $service
    generate_model $service
    generate_docs $service
    
    # 设置Go工作区
    setup_go_workspace $service
}

# 为所有服务生成指定类型的代码
generate_all_services() {
    local type=$1
    echo -e "${BLUE}为所有服务生成 $type 代码...${NC}"
    
    for service in "${SERVICES[@]}"; do
        case $type in
            "api")
                generate_api $service
                ;;
            "rpc")
                generate_rpc $service
                ;;
            "model")
                generate_model $service
                ;;
            "migrate")
                migrate_database $service
                ;;
            "docs")
                generate_docs $service
                ;;
            "all")
                generate_all $service
                ;;
            "clean")
                clean_service $service
                ;;
            "tidy")
                tidy_service $service
                ;;
            "workspace")
                setup_go_workspace $service
                ;;
        esac
        echo ""
    done
}

# 显示项目结构
show_structure() {
    local service=$1
    echo -e "${BLUE}$service 服务目录结构:${NC}"
    echo "app/$service/"
    echo "├── $service-api.api      # API定义文件"
    echo "├── $service-rpc.proto    # RPC定义文件"
    echo "├── $service.sql          # 数据库初始化文件"
    echo "└── cmd/                  # 生成代码目录"
    echo "    ├── api/              # API服务代码"
    echo "    │   ├── etc/          # 配置文件"
    echo "    │   ├── internal/     # 内部代码"
    echo "    │   └── $service.go   # 入口文件"
    echo "    ├── rpc/              # RPC服务代码"
    echo "    │   ├── etc/          # 配置文件"
    echo "    │   ├── internal/     # 内部代码"
    echo "    │   ├── *.pb.go       # Protocol Buffer生成文件"
    echo "    │   └── $service.go   # 入口文件"
    echo "    └── model/            # 数据模型(API和RPC共用)"
    echo "        └── *.go          # 模型文件"
    echo ""
    echo "docs/$service/            # 服务文档目录"
    echo "├── README.md             # 文档索引"
    echo "├── swagger.json          # API Swagger JSON文档"
    echo "├── swagger.yaml          # API Swagger YAML文档"
    echo "├── $service-rpc.html     # RPC HTML文档"
    echo "├── $service-rpc.md       # RPC Markdown文档"
    echo "└── $service-rpc-proto.md # RPC Proto说明文档"
}

# 主函数
main() {
    echo -e "${BLUE}========================================${NC}"
    echo -e "${BLUE}  HuinongFinancial 微服务代码生成工具  ${NC}"
    echo -e "${BLUE}========================================${NC}"
    
    check_goctl
    
    # 进入项目根目录
    cd "$PROJECT_ROOT"
    
    if [[ "$SERVICE_NAME" == "all" ]]; then
        generate_all_services $TYPE
    else
        case $TYPE in
            "api")
                generate_api $SERVICE_NAME
                ;;
            "rpc")
                generate_rpc $SERVICE_NAME
                ;;
            "model")
                generate_model $SERVICE_NAME
                ;;
            "migrate")
                migrate_database $SERVICE_NAME
                ;;
            "docs")
                generate_docs $SERVICE_NAME
                ;;
            "all")
                generate_all $SERVICE_NAME
                ;;
            "clean")
                clean_service $SERVICE_NAME
                ;;
            "tidy")
                tidy_service $SERVICE_NAME
                ;;
            "workspace")
                setup_go_workspace $SERVICE_NAME
                ;;
        esac
        
        if [[ "$TYPE" != "clean" ]] && [[ "$TYPE" != "tidy" ]] && [[ "$TYPE" != "workspace" ]] && [[ "$TYPE" != "docs" ]]; then
            echo ""
            show_structure $SERVICE_NAME
        fi
    fi
    
    echo ""
    echo -e "${GREEN}========================================${NC}"
    echo -e "${GREEN}          操作完成！                   ${NC}"
    echo -e "${GREEN}========================================${NC}"
    
    if [[ "$TYPE" != "clean" ]] && [[ "$TYPE" != "tidy" ]] && [[ "$TYPE" != "workspace" ]] && [[ "$TYPE" != "docs" ]]; then
        echo ""
        echo -e "${YELLOW}下一步操作:${NC}"
        echo "1. 检查生成的代码是否正确"
        echo "2. 实现业务逻辑 (handler、logic、rpc server)"
        echo "3. 配置服务连接 (数据库、Redis、RPC等)"
        echo "4. 启动服务进行测试"
        echo ""
        echo -e "${YELLOW}目录管理:${NC}"
        echo "- 修改定义文件: app/$SERVICE_NAME/$SERVICE_NAME.{api,proto,sql}"
        echo "- 实现业务逻辑: app/$SERVICE_NAME/cmd/{api,rpc}/internal/"
        echo "- 清理生成代码: $0 $SERVICE_NAME clean"
        echo ""
        echo -e "${YELLOW}数据库操作:${NC}"
        echo "- 执行数据库迁移(先清空): $0 $SERVICE_NAME migrate"
        echo "- 批量迁移所有服务(先清空): $0 all migrate"
        echo "- 重新生成Model: $0 $SERVICE_NAME model"
        echo ""
        echo -e "${YELLOW}依赖管理:${NC}"
        echo "- 执行go mod tidy: $0 $SERVICE_NAME tidy"
        echo "- 批量执行go mod tidy: $0 all tidy"
        echo "- 设置Go工作区: $0 $SERVICE_NAME workspace"
        echo "- 批量设置Go工作区: $0 all workspace"
        echo ""
        echo -e "${YELLOW}文档生成:${NC}"
        echo "- 生成API/RPC文档: $0 $SERVICE_NAME docs"
        echo "- 批量生成文档: $0 all docs"
        echo "- 查看文档: docs/$SERVICE_NAME/README.md"
    elif [[ "$TYPE" == "clean" ]]; then
        echo ""
        echo -e "${YELLOW}代码已清理，可重新生成${NC}"
    elif [[ "$TYPE" == "tidy" ]]; then
        echo ""
        echo -e "${YELLOW}go mod tidy 已完成，依赖已更新${NC}"
    elif [[ "$TYPE" == "workspace" ]]; then
        echo ""
        echo -e "${YELLOW}Go工作区已设置完成，现在可以在cmd目录中进行跨模块开发${NC}"
    elif [[ "$TYPE" == "docs" ]]; then
        echo ""
        echo -e "${YELLOW}文档生成完成！${NC}"
        echo ""
        echo -e "${BLUE}查看方式:${NC}"
        if [[ "$SERVICE_NAME" == "all" ]]; then
            echo "- 查看所有服务文档: find docs/ -name 'README.md'"
            echo "- 浏览器打开Swagger: 将 docs/服务名/swagger.json 导入 https://editor.swagger.io/"
        else
            echo "- 查看文档索引: docs/$SERVICE_NAME/README.md"
            echo "- 浏览器打开Swagger: 将 docs/$SERVICE_NAME/swagger.json 导入 https://editor.swagger.io/"
            echo "- 查看RPC文档: docs/$SERVICE_NAME/$SERVICE_NAME-rpc.html (如果存在)"
        fi
        echo ""
        echo -e "${BLUE}工具推荐:${NC}"
        echo "- Swagger Editor: https://editor.swagger.io/"
        echo "- VSCode OpenAPI 扩展: 可直接预览swagger文件"
        echo "- 安装 protoc-gen-doc 获得更丰富的RPC文档:"
        echo "  go install github.com/pseudomuto/protoc-gen-doc/cmd/protoc-gen-doc@latest"
        echo ""
        echo -e "${BLUE}更新文档:${NC}"
        echo "- 修改API定义后: $0 $SERVICE_NAME docs"
        echo "- 修改RPC定义后: $0 $SERVICE_NAME docs"
    fi
}

# 执行主函数
main 