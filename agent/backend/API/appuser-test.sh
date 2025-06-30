#!/bin/bash

# AppUser JWT 配置修复测试脚本
# 用于验证 JWT 配置修复后的接口功能

set -e

# 配置
BASE_URL="http://localhost:10001"
PHONE="13452552349"
PASSWORD="13452552349"

echo "🚀 开始 AppUser JWT 修复验证测试"
echo "=================================================="

# 颜色输出函数
print_success() {
    echo -e "\033[32m✅ $1\033[0m"
}

print_error() {
    echo -e "\033[31m❌ $1\033[0m"
}

print_info() {
    echo -e "\033[34mℹ️  $1\033[0m"
}

# 1. 测试用户登录
echo ""
print_info "1. 测试用户登录获取新 Token..."

LOGIN_RESPONSE=$(curl -s -w "\n%{http_code}" -X POST \
  "${BASE_URL}/api/v1/auth/login" \
  -H "Content-Type: application/json" \
  -d "{\"phone\":\"${PHONE}\",\"password\":\"${PASSWORD}\"}")

HTTP_CODE=$(echo "$LOGIN_RESPONSE" | tail -n1)
RESPONSE_BODY=$(echo "$LOGIN_RESPONSE" | head -n -1)

if [ "$HTTP_CODE" = "200" ]; then
    print_success "登录成功 (HTTP $HTTP_CODE)"
    echo "响应: $RESPONSE_BODY"
    
    # 提取 token
    TOKEN=$(echo "$RESPONSE_BODY" | grep -o '"token":"[^"]*' | cut -d'"' -f4)
    
    if [ -n "$TOKEN" ]; then
        print_success "Token 获取成功"
        echo "Token: ${TOKEN:0:50}..."
    else
        print_error "Token 提取失败"
        exit 1
    fi
else
    print_error "登录失败 (HTTP $HTTP_CODE)"
    echo "响应: $RESPONSE_BODY"
    exit 1
fi

# 2. 测试获取用户信息（需要 JWT 认证）
echo ""
print_info "2. 测试获取用户信息（JWT 认证）..."

USER_INFO_RESPONSE=$(curl -s -w "\n%{http_code}" -X GET \
  "${BASE_URL}/api/v1/user/info" \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d "{\"phone\":\"${PHONE}\"}")

HTTP_CODE=$(echo "$USER_INFO_RESPONSE" | tail -n1)
RESPONSE_BODY=$(echo "$USER_INFO_RESPONSE" | head -n -1)

if [ "$HTTP_CODE" = "200" ]; then
    print_success "获取用户信息成功 (HTTP $HTTP_CODE)"
    echo "响应: $RESPONSE_BODY"
else
    print_error "获取用户信息失败 (HTTP $HTTP_CODE)"
    echo "响应: $RESPONSE_BODY"
    
    if echo "$RESPONSE_BODY" | grep -q "expired"; then
        print_error "Token 仍然过期！配置可能未生效，请重启服务"
    fi
fi

# 3. 测试修改密码（需要 JWT 认证）
echo ""
print_info "3. 测试修改密码接口（JWT 认证）..."

CHANGE_PWD_RESPONSE=$(curl -s -w "\n%{http_code}" -X POST \
  "${BASE_URL}/api/v1/auth/password" \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d "{\"phone\":\"${PHONE}\",\"old_password\":\"${PASSWORD}\",\"new_password\":\"${PASSWORD}\"}")

HTTP_CODE=$(echo "$CHANGE_PWD_RESPONSE" | tail -n1)
RESPONSE_BODY=$(echo "$CHANGE_PWD_RESPONSE" | head -n -1)

if [ "$HTTP_CODE" = "200" ]; then
    print_success "修改密码接口调用成功 (HTTP $HTTP_CODE)"
    echo "响应: $RESPONSE_BODY"
elif [ "$HTTP_CODE" = "401" ]; then
    print_error "修改密码接口认证失败 (HTTP $HTTP_CODE)"
    echo "响应: $RESPONSE_BODY"
    
    if echo "$RESPONSE_BODY" | grep -q "expired"; then
        print_error "Token 仍然过期！配置可能未生效，请重启服务"
    fi
else
    print_info "修改密码接口响应 (HTTP $HTTP_CODE)"
    echo "响应: $RESPONSE_BODY"
fi

echo ""
echo "=================================================="
print_info "测试完成！"

# 输出总结
echo ""
echo "🔍 问题排查建议："
echo "1. 如果仍有 Token 过期错误，请重启 AppUser 服务"
echo "2. 确认配置文件修改已生效"
echo "3. 检查服务器时间同步"

echo ""
echo "🔄 重启服务命令："
echo "cd app/appuser/cmd/rpc && go run appuserrpc.go &"
echo "cd app/appuser/cmd/api && go run appuser.go &" 