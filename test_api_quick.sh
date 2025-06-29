#!/bin/bash

# ===========================================
# 🚀 惠农金服微服务 - 快速API测试脚本
# ===========================================

set -e

# 颜色输出
GREEN='\033[0;32m'
RED='\033[0;31m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m'

# 服务配置
AUTH_API="http://127.0.0.1:10003/api/v1/auth"
APPUSER_API="http://127.0.0.1:10001/api/v1/appuser"
OAUSER_API="http://127.0.0.1:10002/api/v1/oa"

echo -e "${BLUE}🚀 惠农金服 - 快速API测试${NC}"
echo "=========================================="

# 1. 测试C端用户登录
echo -e "\n${YELLOW}📱 测试C端用户登录...${NC}"
APP_LOGIN_RESPONSE=$(curl -s -X POST "$AUTH_API/app/login" \
  -H "Content-Type: application/json" \
  -d '{"account":"13800138000","password":"123456"}')

echo "Response: $APP_LOGIN_RESPONSE"

# 提取Token
APP_TOKEN=$(echo "$APP_LOGIN_RESPONSE" | grep -o '"accessToken":"[^"]*"' | cut -d'"' -f4)

if [ -n "$APP_TOKEN" ]; then
    echo -e "${GREEN}✅ C端登录成功，Token: ${APP_TOKEN:0:20}...${NC}"
else
    echo -e "${RED}❌ C端登录失败${NC}"
    exit 1
fi

# 2. 测试B端管理员登录
echo -e "\n${YELLOW}👨‍💼 测试B端管理员登录...${NC}"
OA_LOGIN_RESPONSE=$(curl -s -X POST "$AUTH_API/oa/login" \
  -H "Content-Type: application/json" \
  -d '{"username":"admin","password":"admin123"}')

echo "Response: $OA_LOGIN_RESPONSE"

# 提取Token
OA_TOKEN=$(echo "$OA_LOGIN_RESPONSE" | grep -o '"accessToken":"[^"]*"' | cut -d'"' -f4)

if [ -n "$OA_TOKEN" ]; then
    echo -e "${GREEN}✅ B端登录成功，Token: ${OA_TOKEN:0:20}...${NC}"
else
    echo -e "${RED}❌ B端登录失败${NC}"
fi

# 3. 测试获取用户信息
echo -e "\n${YELLOW}👤 测试获取用户信息...${NC}"
USER_INFO_RESPONSE=$(curl -s -X GET "$APPUSER_API/info" \
  -H "Authorization: Bearer $APP_TOKEN")

echo "Response: $USER_INFO_RESPONSE"

if echo "$USER_INFO_RESPONSE" | grep -q '"code":200'; then
    echo -e "${GREEN}✅ 获取用户信息成功${NC}"
else
    echo -e "${RED}❌ 获取用户信息失败${NC}"
fi

# 4. 测试更新用户档案
echo -e "\n${YELLOW}📝 测试更新用户档案...${NC}"
UPDATE_RESPONSE=$(curl -s -X PUT "$APPUSER_API/profile" \
  -H "Authorization: Bearer $APP_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"nickname":"快速测试用户","age":35,"gender":1,"income":10000.00}')

echo "Response: $UPDATE_RESPONSE"

if echo "$UPDATE_RESPONSE" | grep -q '"code":200'; then
    echo -e "${GREEN}✅ 更新用户档案成功${NC}"
else
    echo -e "${RED}❌ 更新用户档案失败${NC}"
fi

# 5. 测试管理员接口（如果有OA Token）
if [ -n "$OA_TOKEN" ]; then
    echo -e "\n${YELLOW}🛡️ 测试管理员获取用户列表...${NC}"
    ADMIN_RESPONSE=$(curl -s -X GET "$OAUSER_API/users/?page=1&size=5" \
      -H "Authorization: Bearer $OA_TOKEN")
    
    echo "Response: $ADMIN_RESPONSE"
    
    if echo "$ADMIN_RESPONSE" | grep -q '"code":200'; then
        echo -e "${GREEN}✅ 管理员接口调用成功${NC}"
    else
        echo -e "${RED}❌ 管理员接口调用失败${NC}"
    fi
fi

# 6. 测试错误场景
echo -e "\n${YELLOW}❌ 测试错误登录...${NC}"
ERROR_RESPONSE=$(curl -s -X POST "$AUTH_API/app/login" \
  -H "Content-Type: application/json" \
  -d '{"account":"13800138000","password":"wrong_password"}')

echo "Response: $ERROR_RESPONSE"

if echo "$ERROR_RESPONSE" | grep -q -v '"code":200'; then
    echo -e "${GREEN}✅ 错误处理正常${NC}"
else
    echo -e "${RED}❌ 错误处理异常${NC}"
fi

echo -e "\n${BLUE}=========================================="
echo -e "🎉 快速测试完成！${NC}" 