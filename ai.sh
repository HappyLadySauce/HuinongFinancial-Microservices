#!/bin/bash

echo "🔧 修复惠农金服微服务 - 数据库字段映射问题"
echo "=================================================="

# 1. 修复AppUser模型中的字段映射
echo "📋 1. 修复AppUser模型字段映射..."

# 更新AppUser模型定义，将account字段改为phone
sed -i 's/Account.*string.*`db:"account"`.*\/\/ 账号(手机号)/Phone    string    `db:"phone"`    \/\/ 手机号/g' app/appuser/cmd/model/appUsersModel_gen.go

# 更新缓存键名
sed -i 's/cacheAppUsersAccountPrefix = "cache:appUsers:account:"/cacheAppUsersPhonePrefix = "cache:appUsers:phone:"/g' app/appuser/cmd/model/appUsersModel_gen.go

# 更新FindOneByAccount方法为FindOneByPhone
sed -i 's/FindOneByAccount(ctx context.Context, account string) (\*AppUsers, error)/FindOneByPhone(ctx context.Context, phone string) (*AppUsers, error)/g' app/appuser/cmd/model/appUsersModel_gen.go

# 更新方法实现
sed -i 's/func (m \*defaultAppUsersModel) FindOneByAccount(ctx context.Context, account string)/func (m *defaultAppUsersModel) FindOneByPhone(ctx context.Context, phone string)/g' app/appuser/cmd/model/appUsersModel_gen.go

# 更新方法内部的account变量为phone
sed -i 's/appUsersAccountKey := fmt.Sprintf("%s%v", cacheAppUsersAccountPrefix, account)/appUsersPhoneKey := fmt.Sprintf("%s%v", cacheAppUsersPhonePrefix, phone)/g' app/appuser/cmd/model/appUsersModel_gen.go

# 更新SQL查询
sed -i 's/where `account` = ? limit 1/where `phone` = ? limit 1/g' app/appuser/cmd/model/appUsersModel_gen.go

# 更新其他相关的account引用
sed -i 's/data.Account/data.Phone/g' app/appuser/cmd/model/appUsersModel_gen.go
sed -i 's/newData.Account/newData.Phone/g' app/appuser/cmd/model/appUsersModel_gen.go

echo "✅ AppUser模型字段映射已修复"

# 2. 修复OAUser模型中类似的问题
echo "📋 2. 检查OAUser模型..."

# 检查OAUser模型是否有类似问题
if grep -q "Account.*string.*db:\"account\"" app/oauser/cmd/model/oaUsersModel_gen.go 2>/dev/null; then
    echo "⚠️  OAUser模型也存在字段映射问题，正在修复..."
    
    # 更新OAUser模型的account字段为username
    sed -i 's/Account.*string.*`db:"account"`.*\/\/ 账号/Username  string    `db:"username"`  \/\/ 用户名/g' app/oauser/cmd/model/oaUsersModel_gen.go
    
    # 更新相关方法和变量
    sed -i 's/FindOneByAccount/FindOneByUsername/g' app/oauser/cmd/model/oaUsersModel_gen.go
    sed -i 's/cacheOaUsersAccountPrefix/cacheOaUsersUsernamePrefix/g' app/oauser/cmd/model/oaUsersModel_gen.go
    sed -i 's/where `account` = ? limit 1/where `username` = ? limit 1/g' app/oauser/cmd/model/oaUsersModel_gen.go
    sed -i 's/data.Account/data.Username/g' app/oauser/cmd/model/oaUsersModel_gen.go
    
    echo "✅ OAUser模型字段映射已修复"
else
    echo "✅ OAUser模型字段映射正常"
fi

# 3. 重启所有服务
echo "📋 3. 重启微服务..."

# 停止所有服务
pkill -f "auth.go"
pkill -f "authrpc.go" 
pkill -f "appuser.go"
pkill -f "appuserrpc.go"
pkill -f "oauser.go"
pkill -f "oauserrpc.go"

sleep 3

# 重新启动所有服务
echo "   🔑 启动认证服务..."
cd app/auth/cmd/rpc && nohup go run authrpc.go > /dev/null 2>&1 &
cd ../../../../
cd app/auth/cmd/api && nohup go run auth.go > /dev/null 2>&1 &
cd ../../../../

echo "   👤 启动C端用户服务..."
cd app/appuser/cmd/rpc && nohup go run appuserrpc.go > /dev/null 2>&1 &
cd ../../../../
cd app/appuser/cmd/api && nohup go run appuser.go > /dev/null 2>&1 &
cd ../../../../

echo "   👨‍💼 启动B端用户服务..."
cd app/oauser/cmd/rpc && nohup go run oauserrpc.go > /dev/null 2>&1 &
cd ../../../../
cd app/oauser/cmd/api && nohup go run oauser.go > /dev/null 2>&1 &
cd ../../../../

echo "⏳ 等待服务启动..."
sleep 8

# 4. 测试修复效果
echo "📋 4. 测试修复效果..."

# 快速测试
echo "🧪 测试C端用户服务..."
APP_TOKEN=$(curl -s "http://127.0.0.1:10003/api/v1/auth/app/login" \
    -H "Content-Type: application/json" \
    -d '{"account":"13800138000","password":"123456"}' | jq -r '.data.accessToken')

if [ "$APP_TOKEN" != "null" ] && [ "$APP_TOKEN" != "" ]; then
    echo "✅ 认证服务正常，Token: ${APP_TOKEN:0:20}..."
    
    # 测试用户信息接口
    USER_INFO=$(curl -s -H "Authorization: Bearer $APP_TOKEN" "http://127.0.0.1:10001/api/v1/appuser/info")
    echo "🔍 用户信息接口测试结果:"
    echo "$USER_INFO" | jq . 2>/dev/null || echo "$USER_INFO"
else
    echo "❌ 认证服务异常"
fi

echo ""
echo "🧪 测试B端用户服务..."
OA_TOKEN=$(curl -s "http://127.0.0.1:10003/api/v1/auth/oa/login" \
    -H "Content-Type: application/json" \
    -d '{"username":"admin","password":"123456"}' | jq -r '.data.accessToken')

if [ "$OA_TOKEN" != "null" ] && [ "$OA_TOKEN" != "" ]; then
    echo "✅ B端认证服务正常，Token: ${OA_TOKEN:0:20}..."
    
    # 测试管理员信息接口
    ADMIN_INFO=$(curl -s -H "Authorization: Bearer $OA_TOKEN" "http://127.0.0.1:10002/api/v1/oauser/info")
    echo "🔍 管理员信息接口测试结果:"
    echo "$ADMIN_INFO" | jq . 2>/dev/null || echo "$ADMIN_INFO"
else
    echo "❌ B端认证服务异常"
fi

echo ""
echo "🎉 修复完成！"
echo "=================================================="
echo "如果还有问题，请检查服务日志："
echo "- AppUser RPC日志: find app/appuser/cmd/rpc/logs -name '*.log' -exec tail -n 20 {} \\;"
echo "- OAUser RPC日志: find app/oauser/cmd/rpc/logs -name '*.log' -exec tail -n 20 {} \\;"
