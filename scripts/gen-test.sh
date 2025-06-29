#!/bin/bash

# 测试脚本生成器 - 基于OpenAPI/Swagger文档
# 使用方法：./scripts/gen-test.sh [yaml_file] [output_dir]
# 示例：./scripts/gen-test.sh docs/appuser/appuser-api.yaml tests/

set -e

# 颜色定义
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# 帮助信息
show_help() {
    echo -e "${BLUE}测试脚本生成器 - HuinongFinancial微服务${NC}"
    echo ""
    echo "根据OpenAPI/Swagger YAML文档生成curl测试脚本模板"
    echo ""
    echo "使用方法:"
    echo "  $0 [yaml_file] [output_dir]"
    echo ""
    echo "参数说明:"
    echo "  yaml_file:    API文档YAML文件路径"
    echo "  output_dir:   输出目录 (可选，默认为tests/)"
    echo ""
    echo "示例:"
    echo "  $0 docs/appuser/appuser-api.yaml"
    echo "  $0 docs/appuser/appuser-api.yaml tests/appuser/"
    echo ""
    echo "生成的文件:"
    echo "  service-name_test.sh - curl命令测试脚本"
}

# 检查参数
if [[ $# -eq 0 ]] || [[ "$1" == "-h" ]] || [[ "$1" == "--help" ]]; then
    show_help
    exit 0
fi

YAML_FILE=$1
OUTPUT_DIR=${2:-tests}

# 验证参数
if [[ ! -f "$YAML_FILE" ]]; then
    echo -e "${RED}错误: YAML文件不存在: $YAML_FILE${NC}"
    exit 1
fi

# 获取项目根目录和创建输出目录
PROJECT_ROOT=$(cd "$(dirname "$0")/.." && pwd)
OUTPUT_DIR="$PROJECT_ROOT/$OUTPUT_DIR"
mkdir -p "$OUTPUT_DIR"

# 从文件路径提取服务名
SERVICE_NAME=$(basename "$YAML_FILE" | sed 's/-api\.yaml$//' | sed 's/\.yaml$//')

echo -e "${BLUE}项目根目录: $PROJECT_ROOT${NC}"
echo -e "${BLUE}API文档文件: $YAML_FILE${NC}"
echo -e "${BLUE}服务名称: $SERVICE_NAME${NC}"
echo -e "${BLUE}输出目录: $OUTPUT_DIR${NC}"

# 生成curl测试脚本
generate_curl_script() {
    local service_name="$1"
    local output_dir="$2"
    local yaml_file="$3"
    
    local script_path="$output_dir/${service_name}_test.sh"
    
    echo -e "${BLUE}生成curl测试脚本...${NC}"
    
    # 提取基础路径
    local base_path=$(grep "^basePath:" "$yaml_file" | head -1 | sed 's/basePath:[[:space:]]*//' | tr -d '"' || echo "")
    
    cat > "$script_path" << EOF
#!/bin/bash

# $service_name 服务API测试脚本
# 使用方法: ./${service_name}_test.sh [base_url]
# 示例: ./${service_name}_test.sh http://localhost:8080

set -e

# 颜色定义
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# 基础URL
BASE_URL=\${1:-http://localhost:8080}
echo -e "\${BLUE}测试 $(echo $service_name | tr '[:lower:]' '[:upper:]') 服务API\${NC}"
echo -e "\${BLUE}Base URL: \$BASE_URL\${NC}"
echo ""

# 全局变量
TOKEN=""
USER_PHONE="13800138000"
USER_PASSWORD="123456"

# 工具函数
log_info() {
    echo -e "\${BLUE}[INFO] \$1\${NC}"
}

log_success() {
    echo -e "\${GREEN}[SUCCESS] \$1\${NC}"
}

log_error() {
    echo -e "\${RED}[ERROR] \$1\${NC}"
}

log_warning() {
    echo -e "\${YELLOW}[WARNING] \$1\${NC}"
}

# 发送HTTP请求的通用函数
send_request() {
    local method=\$1
    local url=\$2
    local data=\$3
    local headers=\$4
    local description=\$5
    
    echo ""
    log_info "测试: \$description"
    log_info "方法: \$method"
    log_info "URL: \$url"
    
    if [[ -n "\$data" ]]; then
        log_info "请求数据: \$data"
    fi
    
    local response
    if [[ -n "\$data" ]]; then
        response=\$(curl -s -w "\\n%{http_code}" -X "\$method" \\
            -H "Content-Type: application/json" \\
            \$headers \\
            -d "\$data" \\
            "\$url" || echo "000")
    else
        response=\$(curl -s -w "\\n%{http_code}" -X "\$method" \\
            -H "Content-Type: application/json" \\
            \$headers \\
            "\$url" || echo "000")
    fi
    
    local http_code=\$(echo "\$response" | tail -n1)
    local body=\$(echo "\$response" | sed '\$d')
    
    if [[ "\$http_code" -ge 200 && "\$http_code" -lt 300 ]]; then
        log_success "HTTP \$http_code - 请求成功"
        echo "响应: \$body"
        
        # 尝试提取token
        if [[ "\$body" =~ \\"token\\":\\"([^"]+)\\" ]]; then
            TOKEN="\${BASH_REMATCH[1]}"
            log_info "提取到Token: \$TOKEN"
        fi
    else
        log_error "HTTP \$http_code - 请求失败"
        echo "响应: \$body"
    fi
    
    echo "----------------------------------------"
}

EOF

    # 添加API测试函数（基于appuser API文档的已知接口）
    if [[ "$service_name" == "appuser" ]]; then
        cat >> "$script_path" << 'EOF'
# 用户注册
test_register() {
    local url="$BASE_URL/api/v1/auth/register"
    local method="POST"
    local description="用户注册"
    local data='{"phone": "'"$USER_PHONE"'", "password": "'"$USER_PASSWORD"'"}'
    local headers=""
    send_request "$method" "$url" "$data" "$headers" "$description"
}

# 用户登录
test_login() {
    local url="$BASE_URL/api/v1/auth/login"
    local method="POST"
    local description="用户登录"
    local data='{"phone": "'"$USER_PHONE"'", "password": "'"$USER_PASSWORD"'"}'
    local headers=""
    send_request "$method" "$url" "$data" "$headers" "$description"
}

# 修改密码
test_change_password() {
    local url="$BASE_URL/api/v1/auth/password"
    local method="POST"
    local description="修改密码"
    local data='{"phone": "'"$USER_PHONE"'", "old_password": "'"$USER_PASSWORD"'", "new_password": "newpass123"}'
    local headers=""
    if [[ -n "$TOKEN" ]]; then
        headers="-H \"Authorization: Bearer $TOKEN\""
    fi
    send_request "$method" "$url" "$data" "$headers" "$description"
}

# 用户登出
test_logout() {
    local url="$BASE_URL/api/v1/auth/logout"
    local method="POST"
    local description="用户登出"
    local data='{"token": "'"$TOKEN"'"}'
    local headers=""
    if [[ -n "$TOKEN" ]]; then
        headers="-H \"Authorization: Bearer $TOKEN\""
    fi
    send_request "$method" "$url" "$data" "$headers" "$description"
}

# 获取用户信息
test_get_user_info() {
    local url="$BASE_URL/api/v1/user/info"
    local method="GET"
    local description="获取用户信息"
    local data='{"phone": "'"$USER_PHONE"'"}'
    local headers=""
    if [[ -n "$TOKEN" ]]; then
        headers="-H \"Authorization: Bearer $TOKEN\""
    fi
    send_request "$method" "$url" "$data" "$headers" "$description"
}

# 更新用户信息
test_update_user_info() {
    local url="$BASE_URL/api/v1/user/info"
    local method="PUT"
    local description="更新用户信息"
    local data='{"user_info": {"id": 1, "phone": "'"$USER_PHONE"'", "name": "测试用户", "nickname": "test", "age": 25, "gender": 1, "occupation": "工程师", "address": "北京市", "income": 10000.0, "status": 1, "created_at": 1640995200, "updated_at": 1640995200}}'
    local headers=""
    if [[ -n "$TOKEN" ]]; then
        headers="-H \"Authorization: Bearer $TOKEN\""
    fi
    send_request "$method" "$url" "$data" "$headers" "$description"
}

# 删除用户
test_delete_user() {
    local url="$BASE_URL/api/v1/user/delete"
    local method="POST"
    local description="删除用户"
    local data='{"phone": "'"$USER_PHONE"'"}'
    local headers=""
    if [[ -n "$TOKEN" ]]; then
        headers="-H \"Authorization: Bearer $TOKEN\""
    fi
    send_request "$method" "$url" "$data" "$headers" "$description"
}

EOF
    else
        # 对于其他服务，添加通用的API测试函数模板
        cat >> "$script_path" << 'EOF'
# TODO: 请根据具体的API文档添加测试函数
# 示例函数模板:
# test_example_api() {
#     local url="$BASE_URL/api/v1/example"
#     local method="POST"
#     local description="示例API"
#     local data='{"key": "value"}'
#     local headers=""
#     if [[ -n "$TOKEN" ]]; then
#         headers="-H \"Authorization: Bearer $TOKEN\""
#     fi
#     send_request "$method" "$url" "$data" "$headers" "$description"
# }

EOF
    fi

    # 添加主测试流程
    if [[ "$service_name" == "appuser" ]]; then
        cat >> "$script_path" << 'EOF'
# 主测试流程
main() {
    log_info "开始 appuser 服务API测试"
    
    # 按照逻辑顺序执行测试
    test_register
    test_login
    test_get_user_info
    test_update_user_info
    test_change_password
    test_logout
    # test_delete_user  # 注释掉删除用户测试，避免影响后续测试
    
    log_success "appuser 服务API测试完成"
}

EOF
    else
        cat >> "$script_path" << EOF
# 主测试流程
main() {
    log_info "开始 $service_name 服务API测试"
    
    # TODO: 按照逻辑顺序执行测试
    # test_example_api
    
    log_success "$service_name 服务API测试完成"
}

EOF
    fi

    # 添加脚本执行入口
    cat >> "$script_path" << 'EOF'
# 执行主函数
main "$@"
EOF

    chmod +x "$script_path"
    echo -e "${GREEN}✓ 生成curl测试脚本: $script_path${NC}"
}

# 主函数
main() {
    echo -e "${BLUE}开始生成测试脚本...${NC}"
    
    # 生成curl测试脚本
    generate_curl_script "$SERVICE_NAME" "$OUTPUT_DIR" "$YAML_FILE"
    
    echo ""
    echo -e "${GREEN}🎉 测试脚本生成完成!${NC}"
    echo -e "${BLUE}输出目录: $OUTPUT_DIR${NC}"
    echo ""
    echo -e "${YELLOW}使用说明:${NC}"
    echo "  执行curl测试: ./$OUTPUT_DIR/${SERVICE_NAME}_test.sh [base_url]"
    echo "  示例: ./$OUTPUT_DIR/${SERVICE_NAME}_test.sh http://localhost:8080"
    echo ""
    echo -e "${YELLOW}注意事项:${NC}"
    if [[ "$SERVICE_NAME" == "appuser" ]]; then
        echo "  - 已为appuser服务生成完整的API测试函数"
        echo "  - 测试顺序: 注册 -> 登录 -> 获取信息 -> 更新信息 -> 修改密码 -> 登出"
        echo "  - 删除用户测试默认被注释，如需测试请手动启用"
    else
        echo "  - 请根据具体的API文档完善测试函数"
        echo "  - 当前为通用模板，需要手动添加具体的API测试逻辑"
    fi
}

# 执行主函数
main "$@"
