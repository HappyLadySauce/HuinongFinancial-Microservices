#!/bin/bash

# ===========================================
# 🧪 惠农金服微服务 - 完整API接口测试脚本
# ===========================================

set -e

# 颜色输出
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
PURPLE='\033[0;35m'
CYAN='\033[0;36m'
NC='\033[0m' # No Color

# 服务配置
AUTH_API="http://127.0.0.1:10003/api/v1/auth"
APPUSER_API="http://127.0.0.1:10001/api/v1/appuser"
OAUSER_API="http://127.0.0.1:10002/api/v1/oa"

# 全局变量存储Token
APP_ACCESS_TOKEN=""
APP_REFRESH_TOKEN=""
OA_ACCESS_TOKEN=""
OA_REFRESH_TOKEN=""

# 工具函数
print_header() {
    echo -e "\n${CYAN}===========================================${NC}"
    echo -e "${CYAN}$1${NC}"
    echo -e "${CYAN}===========================================${NC}"
}

print_step() {
    echo -e "\n${BLUE}📋 $1${NC}"
}

print_success() {
    echo -e "${GREEN}✅ $1${NC}"
}

print_error() {
    echo -e "${RED}❌ $1${NC}"
}

print_warning() {
    echo -e "${YELLOW}⚠️  $1${NC}"
}

# 通用HTTP请求函数
make_request() {
    local method=$1
    local url=$2
    local data=$3
    local headers=$4
    local description=$5
    
    echo -e "\n${PURPLE}🔄 $description${NC}"
    echo -e "${YELLOW}   $method $url${NC}"
    
    if [ -n "$data" ]; then
        echo -e "${YELLOW}   Data: $data${NC}"
    fi
    
    local cmd="curl -s -w \"\\n%{http_code}\" -X $method \"$url\""
    
    if [ -n "$headers" ]; then
        cmd="$cmd $headers"
    fi
    
    if [ -n "$data" ]; then
        cmd="$cmd -H \"Content-Type: application/json\" -d '$data'"
    fi
    
    local response=$(eval $cmd)
    local http_code=$(echo "$response" | tail -n1)
    local body=$(echo "$response" | sed '$d')
    
    echo -e "${CYAN}   Response: $body${NC}"
    echo -e "${CYAN}   HTTP Code: $http_code${NC}"
    
    if [[ $http_code -ge 200 && $http_code -lt 300 ]]; then
        print_success "$description - 成功"
        echo "$body"
    else
        print_error "$description - 失败 (HTTP $http_code)"
        echo "$body"
    fi
}

# 提取Token函数
extract_token() {
    local response=$1
    local token_type=$2
    echo "$response" | grep -o "\"$token_type\":\"[^\"]*\"" | cut -d'"' -f4
}

# ===========================================
# 🔐 认证服务测试
# ===========================================
test_auth_service() {
    print_header "🔐 认证服务 (Auth Service) 测试"
    
    # 1. C端用户登录测试
    print_step "1. C端用户登录测试"
    
    # 正确凭据登录
    local login_response=$(make_request "POST" "$AUTH_API/app/login" \
        '{"account":"13800138000","password":"123456"}' \
        "" "C端用户登录 - 正确凭据")
    
    APP_ACCESS_TOKEN=$(echo "$login_response" | grep -o '"accessToken":"[^"]*"' | cut -d'"' -f4)
    APP_REFRESH_TOKEN=$(echo "$login_response" | grep -o '"refreshToken":"[^"]*"' | cut -d'"' -f4)
    
    if [ -n "$APP_ACCESS_TOKEN" ]; then
        print_success "C端AccessToken获取成功: ${APP_ACCESS_TOKEN:0:20}..."
    else
        print_error "C端AccessToken获取失败"
    fi
    
    # 错误凭据登录
    make_request "POST" "$AUTH_API/app/login" \
        '{"account":"13800138000","password":"wrong_password"}' \
        "" "C端用户登录 - 错误密码"
    
    # 缺少参数
    make_request "POST" "$AUTH_API/app/login" \
        '{"account":"13800138000"}' \
        "" "C端用户登录 - 缺少密码"
    
    # 2. B端管理员登录测试
    print_step "2. B端管理员登录测试"
    
    # 正确凭据登录
    local oa_login_response=$(make_request "POST" "$AUTH_API/oa/login" \
        '{"username":"admin","password":"admin123"}' \
        "" "B端管理员登录 - 正确凭据")
    
    OA_ACCESS_TOKEN=$(echo "$oa_login_response" | grep -o '"accessToken":"[^"]*"' | cut -d'"' -f4)
    OA_REFRESH_TOKEN=$(echo "$oa_login_response" | grep -o '"refreshToken":"[^"]*"' | cut -d'"' -f4)
    
    if [ -n "$OA_ACCESS_TOKEN" ]; then
        print_success "B端AccessToken获取成功: ${OA_ACCESS_TOKEN:0:20}..."
    else
        print_error "B端AccessToken获取失败"
    fi
    
    # 错误凭据登录
    make_request "POST" "$AUTH_API/oa/login" \
        '{"username":"admin","password":"wrong_password"}' \
        "" "B端管理员登录 - 错误密码"
    
    # 3. Token刷新测试
    if [ -n "$APP_REFRESH_TOKEN" ]; then
        print_step "3. Token刷新测试"
        make_request "POST" "$AUTH_API/refresh" \
            "{\"refreshToken\":\"$APP_REFRESH_TOKEN\"}" \
            "" "刷新AccessToken"
    fi
    
    # 4. 用户登出测试
    if [ -n "$APP_ACCESS_TOKEN" ]; then
        print_step "4. 用户登出测试"
        make_request "POST" "$AUTH_API/logout" \
            "" \
            "-H \"Authorization: Bearer $APP_ACCESS_TOKEN\"" \
            "用户登出"
    fi
}

# ===========================================
# 👤 C端用户服务测试
# ===========================================
test_appuser_service() {
    print_header "👤 C端用户服务 (AppUser Service) 测试"
    
    if [ -z "$APP_ACCESS_TOKEN" ]; then
        print_warning "没有有效的C端AccessToken，重新获取..."
        local login_response=$(make_request "POST" "$AUTH_API/app/login" \
            '{"account":"13800138000","password":"123456"}' \
            "" "重新获取C端Token")
        APP_ACCESS_TOKEN=$(echo "$login_response" | grep -o '"accessToken":"[^"]*"' | cut -d'"' -f4)
    fi
    
    if [ -z "$APP_ACCESS_TOKEN" ]; then
        print_error "无法获取C端AccessToken，跳过用户服务测试"
        return
    fi
    
    # 1. 获取用户信息测试
    print_step "1. 获取用户信息测试"
    make_request "GET" "$APPUSER_API/info" \
        "" \
        "-H \"Authorization: Bearer $APP_ACCESS_TOKEN\"" \
        "获取当前用户信息"
    
    # 2. 更新用户档案测试
    print_step "2. 更新用户档案测试"
    make_request "PUT" "$APPUSER_API/profile" \
        '{"nickname":"测试昵称更新","age":33,"gender":1,"occupation":"高级农民","address":"北京市朝阳区测试地址","income":9500.00}' \
        "-H \"Authorization: Bearer $APP_ACCESS_TOKEN\"" \
        "更新用户档案"
    
    # 3. 无权限访问测试
    print_step "3. 无权限访问测试"
    make_request "GET" "$APPUSER_API/info" \
        "" \
        "-H \"Authorization: Bearer invalid_token\"" \
        "使用无效Token访问 - 应该失败"
    
    # 4. 缺少Token测试
    make_request "GET" "$APPUSER_API/info" \
        "" \
        "" \
        "无Token访问 - 应该失败"
}

# ===========================================
# 🛡️ B端管理服务测试
# ===========================================
test_oauser_service() {
    print_header "🛡️ B端管理服务 (OAUser Service) 测试"
    
    if [ -z "$OA_ACCESS_TOKEN" ]; then
        print_warning "没有有效的B端AccessToken，重新获取..."
        local oa_login_response=$(make_request "POST" "$AUTH_API/oa/login" \
            '{"username":"admin","password":"admin123"}' \
            "" "重新获取B端Token")
        OA_ACCESS_TOKEN=$(echo "$oa_login_response" | grep -o '"accessToken":"[^"]*"' | cut -d'"' -f4)
    fi
    
    if [ -z "$OA_ACCESS_TOKEN" ]; then
        print_error "无法获取B端AccessToken，跳过管理服务测试"
        return
    fi
    
    # 1. 获取用户列表测试
    print_step "1. 获取用户列表测试"
    make_request "GET" "$OAUSER_API/users/?page=1&size=10" \
        "" \
        "-H \"Authorization: Bearer $OA_ACCESS_TOKEN\"" \
        "获取用户列表 - 默认分页"
    
    # 带搜索条件
    make_request "GET" "$OAUSER_API/users/?page=1&size=5&keyword=admin&status=1" \
        "" \
        "-H \"Authorization: Bearer $OA_ACCESS_TOKEN\"" \
        "获取用户列表 - 带搜索条件"
    
    # 2. 创建用户测试
    print_step "2. 创建用户测试"
    make_request "POST" "$OAUSER_API/users/" \
        '{"username":"testuser001","password":"test123456","name":"测试用户","email":"test@huinong.com","mobile":"13900000001","roles":["editor"]}' \
        "-H \"Authorization: Bearer $OA_ACCESS_TOKEN\"" \
        "创建新用户"
    
    # 3. 获取用户详情测试
    print_step "3. 获取用户详情测试"
    make_request "GET" "$OAUSER_API/users/testuser001" \
        "" \
        "-H \"Authorization: Bearer $OA_ACCESS_TOKEN\"" \
        "获取用户详情"
    
    # 4. 更新用户测试
    print_step "4. 更新用户测试"
    make_request "PUT" "$OAUSER_API/users/1" \
        '{"name":"更新的用户名","status":1,"roles":["admin","editor"],"email":"updated@huinong.com","mobile":"13900000002"}' \
        "-H \"Authorization: Bearer $OA_ACCESS_TOKEN\"" \
        "更新用户信息"
    
    # 5. 删除用户测试
    print_step "5. 删除用户测试"
    make_request "DELETE" "$OAUSER_API/users/testuser001" \
        "" \
        "-H \"Authorization: Bearer $OA_ACCESS_TOKEN\"" \
        "删除用户"
    
    # 6. 权限测试 - 使用C端Token访问管理接口
    if [ -n "$APP_ACCESS_TOKEN" ]; then
        print_step "6. 权限测试 - C端Token访问管理接口"
        make_request "GET" "$OAUSER_API/users/?page=1&size=10" \
            "" \
            "-H \"Authorization: Bearer $APP_ACCESS_TOKEN\"" \
            "C端用户访问管理接口 - 应该失败"
    fi
}

# ===========================================
# 🧪 数据边界测试
# ===========================================
test_edge_cases() {
    print_header "🧪 数据边界和异常测试"
    
    print_step "1. 参数验证测试"
    
    # 空数据测试
    make_request "POST" "$AUTH_API/app/login" \
        '{}' \
        "" "登录空参数测试"
    
    # 超长数据测试
    local long_string=$(printf "a%.0s" {1..1000})
    make_request "POST" "$AUTH_API/app/login" \
        "{\"account\":\"$long_string\",\"password\":\"123456\"}" \
        "" "登录超长账号测试"
    
    # SQL注入测试
    make_request "POST" "$AUTH_API/app/login" \
        '{"account":"admin'\''OR 1=1--","password":"any"}' \
        "" "SQL注入防护测试"
    
    # 2. 数据类型错误测试
    print_step "2. 数据类型错误测试"
    make_request "PUT" "$APPUSER_API/profile" \
        '{"age":"不是数字","gender":"不是数字"}' \
        "-H \"Authorization: Bearer $APP_ACCESS_TOKEN\"" \
        "错误数据类型测试" || true
    
    # 3. 并发请求测试
    print_step "3. 并发请求测试"
    echo "发起10个并发登录请求..."
    for i in {1..10}; do
        (make_request "POST" "$AUTH_API/app/login" \
            '{"account":"13800138000","password":"123456"}' \
            "" "并发登录测试 $i") &
    done
    wait
    print_success "并发测试完成"
}

# ===========================================
# 📊 性能测试
# ===========================================
test_performance() {
    print_header "📊 简单性能测试"
    
    if [ -z "$APP_ACCESS_TOKEN" ]; then
        print_warning "获取AccessToken用于性能测试..."
        local login_response=$(make_request "POST" "$AUTH_API/app/login" \
            '{"account":"13800138000","password":"123456"}' \
            "" "获取性能测试Token")
        APP_ACCESS_TOKEN=$(echo "$login_response" | grep -o '"accessToken":"[^"]*"' | cut -d'"' -f4)
    fi
    
    if [ -n "$APP_ACCESS_TOKEN" ]; then
        print_step "1. 用户信息查询性能测试 (50次请求)"
        
        local start_time=$(date +%s.%N)
        for i in {1..50}; do
            curl -s -H "Authorization: Bearer $APP_ACCESS_TOKEN" \
                "$APPUSER_API/info" > /dev/null
        done
        local end_time=$(date +%s.%N)
        
        local duration=$(echo "$end_time - $start_time" | bc)
        local qps=$(echo "scale=2; 50 / $duration" | bc)
        
        print_success "50次请求完成，总耗时: ${duration}秒，QPS: $qps"
    fi
}

# ===========================================
# 🎯 综合业务流程测试
# ===========================================
test_business_flow() {
    print_header "🎯 综合业务流程测试"
    
    print_step "1. 完整用户注册→登录→更新资料→登出流程"
    
    # 注意：由于当前没有注册接口，这里模拟已存在用户的完整流程
    
    # Step 1: 登录
    local login_response=$(make_request "POST" "$AUTH_API/app/login" \
        '{"account":"13800138001","password":"123456"}' \
        "" "业务流程 - 用户登录")
    
    local flow_token=$(echo "$login_response" | grep -o '"accessToken":"[^"]*"' | cut -d'"' -f4)
    
    # Step 2: 获取用户信息
    make_request "GET" "$APPUSER_API/info" \
        "" \
        "-H \"Authorization: Bearer $flow_token\"" \
        "业务流程 - 获取用户信息"
    
    # Step 3: 更新用户资料
    make_request "PUT" "$APPUSER_API/profile" \
        '{"nickname":"业务流程测试","age":30,"gender":2,"occupation":"业务测试员","income":8800.00}' \
        "-H \"Authorization: Bearer $flow_token\"" \
        "业务流程 - 更新用户资料"
    
    # Step 4: 再次获取用户信息验证更新
    make_request "GET" "$APPUSER_API/info" \
        "" \
        "-H \"Authorization: Bearer $flow_token\"" \
        "业务流程 - 验证资料更新"
    
    # Step 5: 登出
    make_request "POST" "$AUTH_API/logout" \
        "" \
        "-H \"Authorization: Bearer $flow_token\"" \
        "业务流程 - 用户登出"
    
    print_success "完整业务流程测试完成"
}

# ===========================================
# 🎨 主程序入口
# ===========================================
main() {
    print_header "🚀 惠农金服微服务API测试开始"
    
    echo -e "${CYAN}测试环境信息:${NC}"
    echo -e "  Auth API:    $AUTH_API"
    echo -e "  AppUser API: $APPUSER_API"
    echo -e "  OAUser API:  $OAUSER_API"
    echo -e "  时间:        $(date)"
    
    # 检查服务是否可用
    print_step "检查服务可用性..."
    
    for service in "Auth:$AUTH_API" "AppUser:$APPUSER_API" "OAUser:$OAUSER_API"; do
        local name=$(echo $service | cut -d: -f1)
        local url=$(echo $service | cut -d: -f2)
        
        if curl -s --connect-timeout 5 "$url" > /dev/null 2>&1; then
            print_success "$name 服务可用"
        else
            print_warning "$name 服务不可用，可能影响测试结果"
        fi
    done
    
    # 执行测试套件
    test_auth_service
    test_appuser_service  
    test_oauser_service
    test_edge_cases
    test_performance
    test_business_flow
    
    print_header "🎉 所有测试完成"
    echo -e "${GREEN}如果您看到错误，请检查:${NC}"
    echo -e "  1. 相关微服务是否正常启动"
    echo -e "  2. 数据库和Redis连接是否正常"
    echo -e "  3. 测试数据是否已初始化"
    echo -e "  4. 网络连接是否正常"
}

# 执行主程序
main "$@" 