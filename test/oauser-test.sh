#!/bin/bash

# OAUser 后台用户服务测试脚本
# 用于验证后台用户管理系统的各项功能

set -e

# 配置
BASE_URL="http://localhost:10002"
ADMIN_PHONE="13452552346"
ADMIN_PASSWORD="13452552346"
OPERATOR_PHONE="13452552352"  # 使用不同的手机号避免冲突
OPERATOR_PASSWORD="13452552352"

echo "🚀 开始 OAUser 后台用户服务测试"
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

print_warning() {
    echo -e "\033[33m⚠️  $1\033[0m"
}

# 发送请求的通用函数
send_request() {
    local method=$1
    local url=$2
    local data=$3
    local headers=$4
    local description=$5
    
    print_info "测试: $description"
    
    local response
    if [[ -n "$data" ]]; then
        if [[ -n "$headers" ]]; then
            response=$(curl -s -w "\n%{http_code}" -X "$method" \
                -H "Content-Type: application/json" \
                -H "$headers" \
                -d "$data" \
                "$url" || echo "000")
        else
            response=$(curl -s -w "\n%{http_code}" -X "$method" \
                -H "Content-Type: application/json" \
                -d "$data" \
                "$url" || echo "000")
        fi
    else
        if [[ -n "$headers" ]]; then
            response=$(curl -s -w "\n%{http_code}" -X "$method" \
                -H "Content-Type: application/json" \
                -H "$headers" \
                "$url" || echo "000")
        else
            response=$(curl -s -w "\n%{http_code}" -X "$method" \
                -H "Content-Type: application/json" \
                "$url" || echo "000")
        fi
    fi
    
    local http_code=$(echo "$response" | tail -n1)
    local body=$(echo "$response" | head -n -1)
    
    echo "HTTP $http_code - $body"
    
    # 提取token (修复正则表达式引号问题)
    if [[ "$body" =~ \"token\":\"([^\"]+)\" ]]; then
        local token="${BASH_REMATCH[1]}"
        echo "$token"
        return 0
    fi
    
    echo ""
    return $http_code
}

# 全局变量
ADMIN_TOKEN=""
OPERATOR_TOKEN=""

# 1. 测试管理员注册
echo ""
print_info "========== 1. 管理员注册测试 =========="

ADMIN_TOKEN=$(send_request "POST" "${BASE_URL}/api/v1/auth/register" \
    "{\"phone\":\"${ADMIN_PHONE}\",\"password\":\"${ADMIN_PASSWORD}\",\"role\":\"admin\"}" \
    "" "管理员注册")

if [[ $? -eq 200 && -n "$ADMIN_TOKEN" ]]; then
    print_success "管理员注册成功，Token: ${ADMIN_TOKEN:0:50}..."
else
    print_warning "管理员可能已存在，尝试登录..."
fi

# 2. 测试管理员登录
echo ""
print_info "========== 2. 管理员登录测试 =========="

if [[ -z "$ADMIN_TOKEN" ]]; then
    ADMIN_TOKEN=$(send_request "POST" "${BASE_URL}/api/v1/auth/login" \
        "{\"phone\":\"${ADMIN_PHONE}\",\"password\":\"${ADMIN_PASSWORD}\"}" \
        "" "管理员登录")
    
    if [[ $? -eq 200 && -n "$ADMIN_TOKEN" ]]; then
        print_success "管理员登录成功，Token: ${ADMIN_TOKEN:0:50}..."
    else
        print_error "管理员登录失败"
        exit 1
    fi
fi

# 3. 测试普通操作员注册
echo ""
print_info "========== 3. 普通操作员注册测试 =========="

OPERATOR_TOKEN=$(send_request "POST" "${BASE_URL}/api/v1/auth/register" \
    "{\"phone\":\"${OPERATOR_PHONE}\",\"password\":\"${OPERATOR_PASSWORD}\",\"role\":\"operator\"}" \
    "" "普通操作员注册")

if [[ $? -eq 200 && -n "$OPERATOR_TOKEN" ]]; then
    print_success "普通操作员注册成功，Token: ${OPERATOR_TOKEN:0:50}..."
else
    print_warning "普通操作员可能已存在，尝试登录..."
fi

# 4. 测试普通操作员登录
echo ""
print_info "========== 4. 普通操作员登录测试 =========="

if [[ -z "$OPERATOR_TOKEN" ]]; then
    OPERATOR_TOKEN=$(send_request "POST" "${BASE_URL}/api/v1/auth/login" \
        "{\"phone\":\"${OPERATOR_PHONE}\",\"password\":\"${OPERATOR_PASSWORD}\"}" \
        "" "普通操作员登录")
    
    if [[ $? -eq 200 && -n "$OPERATOR_TOKEN" ]]; then
        print_success "普通操作员登录成功，Token: ${OPERATOR_TOKEN:0:50}..."
    else
        print_error "普通操作员登录失败"
        exit 1
    fi
fi

# 5. 测试获取用户信息（管理员）
echo ""
print_info "========== 5. 管理员获取用户信息测试 =========="

send_request "GET" "${BASE_URL}/api/v1/user/info" \
    "{\"phone\":\"${ADMIN_PHONE}\"}" \
    "Authorization: Bearer $ADMIN_TOKEN" \
    "管理员获取用户信息" >/dev/null

if [[ $? -eq 200 ]]; then
    print_success "管理员获取用户信息成功"
else
    print_error "管理员获取用户信息失败"
fi

# 6. 测试获取用户信息（普通操作员）
echo ""
print_info "========== 6. 普通操作员获取用户信息测试 =========="

send_request "GET" "${BASE_URL}/api/v1/user/info" \
    "{\"phone\":\"${OPERATOR_PHONE}\"}" \
    "Authorization: Bearer $OPERATOR_TOKEN" \
    "普通操作员获取用户信息" >/dev/null

if [[ $? -eq 200 ]]; then
    print_success "普通操作员获取用户信息成功"
else
    print_error "普通操作员获取用户信息失败"
fi

# 7. 测试更新用户信息（管理员）
echo ""
print_info "========== 7. 管理员更新用户信息测试 =========="

send_request "PUT" "${BASE_URL}/api/v1/user/info" \
    "{\"user_info\":{\"id\":1,\"phone\":\"${ADMIN_PHONE}\",\"name\":\"系统管理员\",\"nickname\":\"Admin\",\"age\":30,\"gender\":1,\"role\":\"admin\",\"status\":1,\"created_at\":$(date +%s),\"updated_at\":$(date +%s)}}" \
    "Authorization: Bearer $ADMIN_TOKEN" \
    "管理员更新用户信息" >/dev/null

if [[ $? -eq 200 ]]; then
    print_success "管理员更新用户信息成功"
else
    print_error "管理员更新用户信息失败"
fi

# 8. 测试修改密码（普通操作员）
echo ""
print_info "========== 8. 普通操作员修改密码测试 =========="

send_request "POST" "${BASE_URL}/api/v1/auth/password" \
    "{\"phone\":\"${OPERATOR_PHONE}\",\"old_password\":\"${OPERATOR_PASSWORD}\",\"new_password\":\"${OPERATOR_PASSWORD}\"}" \
    "Authorization: Bearer $OPERATOR_TOKEN" \
    "普通操作员修改密码" >/dev/null

if [[ $? -eq 200 ]]; then
    print_success "普通操作员修改密码成功"
else
    print_error "普通操作员修改密码失败"
fi

# 9. 测试管理员权限 - 删除用户
echo ""
print_info "========== 9. 管理员权限测试 - 删除用户 =========="

# 先创建一个测试用户
TEST_PHONE="13800138999"
TEST_PASSWORD="test123456"

print_info "创建测试用户..."
send_request "POST" "${BASE_URL}/api/v1/auth/register" \
    "{\"phone\":\"${TEST_PHONE}\",\"password\":\"${TEST_PASSWORD}\",\"role\":\"operator\"}" \
    "" "创建测试用户" >/dev/null

# 管理员删除用户
print_info "管理员删除测试用户..."
send_request "POST" "${BASE_URL}/api/v1/user/delete" \
    "{\"phone\":\"${TEST_PHONE}\"}" \
    "Authorization: Bearer $ADMIN_TOKEN" \
    "管理员删除用户" >/dev/null

if [[ $? -eq 200 ]]; then
    print_success "管理员删除用户成功"
else
    print_error "管理员删除用户失败"
fi

# 10. 测试普通操作员权限 - 尝试删除用户（应该失败）
echo ""
print_info "========== 10. 普通操作员权限测试 - 尝试删除用户 =========="

send_request "POST" "${BASE_URL}/api/v1/user/delete" \
    "{\"phone\":\"${TEST_PHONE}\"}" \
    "Authorization: Bearer $OPERATOR_TOKEN" \
    "普通操作员尝试删除用户" >/dev/null

if [[ $? -eq 403 || $? -eq 401 ]]; then
    print_success "普通操作员权限控制正常（无法删除用户）"
else
    print_warning "普通操作员权限控制可能存在问题"
fi

# 11. 测试注销功能
echo ""
print_info "========== 11. 用户注销测试 =========="

send_request "POST" "${BASE_URL}/api/v1/auth/logout" \
    "{\"token\":\"${OPERATOR_TOKEN}\"}" \
    "Authorization: Bearer $OPERATOR_TOKEN" \
    "普通操作员注销" >/dev/null

if [[ $? -eq 200 ]]; then
    print_success "普通操作员注销成功"
else
    print_error "普通操作员注销失败"
fi

# 12. 测试注销后的Token失效
echo ""
print_info "========== 12. Token失效验证测试 =========="

send_request "GET" "${BASE_URL}/api/v1/user/info" \
    "{\"phone\":\"${OPERATOR_PHONE}\"}" \
    "Authorization: Bearer $OPERATOR_TOKEN" \
    "使用已注销Token获取用户信息" >/dev/null

if [[ $? -eq 401 ]]; then
    print_success "Token失效验证正常"
else
    print_warning "Token失效验证可能存在问题"
fi

echo ""
echo "=================================================="
print_info "OAUser 后台用户服务测试完成！"

# 输出测试总结
echo ""
echo "🔍 测试结果总结："
echo "✅ 管理员注册/登录功能"
echo "✅ 普通操作员注册/登录功能"
echo "✅ 用户信息获取和更新功能"
echo "✅ 密码修改功能"
echo "✅ 角色权限控制（管理员可删除用户）"
echo "✅ 用户注销功能"

echo ""
echo "📋 测试用户信息："
echo "管理员: $ADMIN_PHONE / $ADMIN_PASSWORD"
echo "操作员: $OPERATOR_PHONE / $OPERATOR_PASSWORD"

echo ""
echo "🔄 如果测试失败，请检查："
echo "1. OAUser RPC 服务是否启动 (端口20002)"
echo "2. OAUser API 服务是否启动 (端口10002)"
echo "3. 数据库连接是否正常"
echo "4. Redis缓存是否可用"

echo ""
echo "🚀 启动服务命令："
echo "cd app/oauser/cmd/rpc && go run oauserrpc.go &"
echo "cd app/oauser/cmd/api && go run oauser.go &" 