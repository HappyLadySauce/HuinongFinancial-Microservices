#!/bin/bash

# appuser 服务API测试脚本
# 使用方法: ./appuser_test.sh [base_url]
# 示例: ./appuser_test.sh http://localhost:8080

set -e

# 颜色定义
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# 基础URL
BASE_URL=${1:-http://localhost:10001}
echo -e "${BLUE}测试 APPUSER 服务API${NC}"
echo -e "${BLUE}Base URL: $BASE_URL${NC}"
echo ""

# 全局变量
TOKEN="eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyX2lkIjoxMDA5LCJwaG9uZSI6IjEzNDUyNTUyMzQ5IiwidXNlcl90eXBlIjoiYXBwIiwiaXNzIjoiaHVpbm9uZy1maW5hbmNpYWwiLCJleHAiOjE3NTEyMDk4MDcsIm5iZiI6MTc1MTIwNjIwNywiaWF0IjoxNzUxMjA2MjA3fQ.zcYAcZczymWQ74YiaShy2J_ZU4AEhfT4MYUQnzysyVM"
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
    local use_auth=$4
    local description=$5
    
    echo ""
    log_info "测试: $description"
    log_info "方法: $method"
    log_info "URL: $url"
    
    if [[ -n "$data" ]]; then
        log_info "请求数据: $data"
    fi
    
    local response
    if [[ "$use_auth" == "true" && -n "$TOKEN" ]]; then
        log_info "使用JWT认证: Bearer $TOKEN"
        if [[ -n "$data" ]]; then
            response=$(curl -s -w "\n%{http_code}" -X "$method" \
                -H "Content-Type: application/json" \
                -H "Authorization: Bearer $TOKEN" \
                -d "$data" \
                "$url" || echo "000")
        else
            response=$(curl -s -w "\n%{http_code}" -X "$method" \
                -H "Content-Type: application/json" \
                -H "Authorization: Bearer $TOKEN" \
                "$url" || echo "000")
        fi
    else
        if [[ -n "$data" ]]; then
            response=$(curl -s -w "\n%{http_code}" -X "$method" \
                -H "Content-Type: application/json" \
                -d "$data" \
                "$url" || echo "000")
        else
            response=$(curl -s -w "\n%{http_code}" -X "$method" \
                -H "Content-Type: application/json" \
                "$url" || echo "000")
        fi
    fi
    
    local http_code=$(echo "$response" | tail -n1)
    local body=$(echo "$response" | sed '$d')
    
    if [[ "$http_code" -ge 200 && "$http_code" -lt 300 ]]; then
        log_success "HTTP $http_code - 请求成功"
        echo "响应: $body"
        
        # 尝试提取token
        if [[ "$body" =~ \"token\":\"([^\"]+)\" ]]; then
            TOKEN="${BASH_REMATCH[1]}"
            log_info "提取到Token: $TOKEN"
        fi
    else
        log_error "HTTP $http_code - 请求失败"
        echo "响应: $body"
    fi
    
    echo "----------------------------------------"
}

# 用户注册
test_register() {
    local url="$BASE_URL/api/v1/auth/register"
    local method="POST"
    local description="用户注册"
    local data='{"phone": "'"$USER_PHONE"'", "password": "'"$USER_PASSWORD"'"}'
    send_request "$method" "$url" "$data" "false" "$description"
}

# 用户登录
test_login() {
    local url="$BASE_URL/api/v1/auth/login"
    local method="POST"
    local description="用户登录"
    local data='{"phone": "'"$USER_PHONE"'", "password": "'"$USER_PASSWORD"'"}'
    send_request "$method" "$url" "$data" "false" "$description"
}

# 修改密码
test_change_password() {
    local url="$BASE_URL/api/v1/auth/password"
    local method="POST"
    local description="修改密码"
    local data='{"phone": "'"$USER_PHONE"'", "old_password": "'"$USER_PASSWORD"'", "new_password": "newpass123"}'
    send_request "$method" "$url" "$data" "true" "$description"
}

# 用户登出
test_logout() {
    local url="$BASE_URL/api/v1/auth/logout"
    local method="POST"
    local description="用户登出"
    local data='{"token": "'"$TOKEN"'"}'
    send_request "$method" "$url" "$data" "true" "$description"
}

# 获取用户信息
test_get_user_info() {
    local url="$BASE_URL/api/v1/user/info?phone=$USER_PHONE"
    local method="GET"
    local description="获取用户信息"
    local data=""
    send_request "$method" "$url" "$data" "true" "$description"
}

# 更新用户信息
test_update_user_info() {
    local url="$BASE_URL/api/v1/user/info"
    local method="PUT"
    local description="更新用户信息"
    local data='{"user_info": {"id": 1, "phone": "'"$USER_PHONE"'", "name": "测试用户", "nickname": "test", "age": 25, "gender": 1, "occupation": "工程师", "address": "北京市", "income": 10000.0, "status": 1, "created_at": 1640995200, "updated_at": 1640995200}}'
    send_request "$method" "$url" "$data" "true" "$description"
}

# 删除用户
test_delete_user() {
    local url="$BASE_URL/api/v1/user/delete"
    local method="POST"
    local description="删除用户"
    local data='{"phone": "'"$USER_PHONE"'"}'
    send_request "$method" "$url" "$data" "true" "$description"
}

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

# 执行主函数
main "$@"
