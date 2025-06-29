#!/bin/bash

# oauser 服务API测试脚本
# 使用方法: ./oauser_test.sh [base_url]
# 示例: ./oauser_test.sh http://localhost:8080

set -e

# 颜色定义
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# 基础URL
BASE_URL=${1:-http://localhost:10002}
echo -e "${BLUE}测试 OAUSER 服务API${NC}"
echo -e "${BLUE}Base URL: $BASE_URL${NC}"
echo ""

# 全局变量
TOKEN="eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyX2lkIjoyMDA1LCJwaG9uZSI6IjEzNDUyNTUyMzQ5IiwidXNlcl90eXBlIjoib2EiLCJpc3MiOiJodWlub25nLWZpbmFuY2lhbCIsImV4cCI6MTc1MTIxNDU3NywibmJmIjoxNzUxMjEwOTc3LCJpYXQiOjE3NTEyMTA5Nzd9.U-Wj7iz9ngnmXFsUpG3t7KLrh1yzQsEYlSErW7ghWbg"
USER_PHONE="13452552349"
USER_PASSWORD="13452552349"

# 工具函数
log_info() {
    echo -e "${BLUE}[INFO] $1${NC}"
}

log_success() {
    echo -e "${GREEN}[SUCCESS] $1${NC}"
}

log_error() {
    echo -e "${RED}[ERROR] $1${NC}"
}

log_warning() {
    echo -e "${YELLOW}[WARNING] $1${NC}"
}

# 发送HTTP请求的通用函数
send_request() {
    local method=$1
    local url=$2
    local data=$3
    local headers=$4
    local description=$5
    
    echo ""
    log_info "测试: $description"
    log_info "方法: $method"
    log_info "URL: $url"
    
    if [[ -n "$data" ]]; then
        log_info "请求数据: $data"
    fi
    
    local response
    if [[ -n "$data" ]]; then
        response=$(curl -s -w "\n%{http_code}" -X "$method" \
            -H "Content-Type: application/json" \
            $headers \
            -d "$data" \
            "$url" || echo "000")
    else
        response=$(curl -s -w "\n%{http_code}" -X "$method" \
            -H "Content-Type: application/json" \
            $headers \
            "$url" || echo "000")
    fi
    
    local http_code=$(echo "$response" | tail -n1)
    local body=$(echo "$response" | sed '$d')
    
    if [[ "$http_code" -ge 200 && "$http_code" -lt 300 ]]; then
        log_success "HTTP $http_code - 请求成功"
        echo "响应: $body"
        
        # 尝试提取token
        if [[ "$body" =~ \"token\":\"([^"]+)\" ]]; then
            TOKEN="${BASH_REMATCH[1]}"
            log_info "提取到Token: $TOKEN"
        fi
    else
        log_error "HTTP $http_code - 请求失败"
        echo "响应: $body"
    fi
    
    echo "----------------------------------------"
}

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

# 主测试流程
main() {
    log_info "开始 oauser 服务API测试"
    
    # TODO: 按照逻辑顺序执行测试
    # test_example_api
    
    log_success "oauser 服务API测试完成"
}

# 执行主函数
main "$@"
