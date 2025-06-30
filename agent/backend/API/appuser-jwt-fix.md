# AppUser JWT 配置修复文档

## 问题描述

### 原始问题
- JWT Token 过期错误：`Token is expired`
- 配置不一致：AppUser 服务使用独立的 JWT 密钥

### 根本原因
1. **JWT 密钥不统一**：
   - AppUser 服务：`huinong-appuser-secret`
   - 其他服务：`huinong-auth-access-secret`

2. **Token 时间异常**：
   - Token 生成时间与系统时间不匹配
   - 导致 Token 已过期超过 23 小时

## 解决方案

### 配置统一化
```yaml
# app/appuser/cmd/api/etc/appuser.yaml
Auth:
  AccessSecret: "huinong-auth-access-secret"
  AccessExpire: 3600

# app/appuser/cmd/rpc/etc/appuserrpc.yaml
JwtAuth:
  AccessSecret: "huinong-auth-access-secret" 
  AccessExpire: 3600
```

## 影响的接口

### 需要重新获取 Token 的接口
1. `POST /api/v1/auth/login` - 用户登录
2. `POST /api/v1/auth/register` - 用户注册

### 受保护的接口（需要新 Token）
1. `GET /api/v1/user/info` - 获取用户信息
2. `PUT /api/v1/user/info` - 更新用户信息
3. `POST /api/v1/user/delete` - 删除用户
4. `POST /api/v1/auth/logout` - 用户注销
5. `POST /api/v1/auth/password` - 修改密码

## 验证步骤

1. **重启服务**
   ```bash
   # 重启 AppUser RPC 服务
   cd app/appuser/cmd/rpc
   go run appuserrpc.go
   
   # 重启 AppUser API 服务  
   cd app/appuser/cmd/api
   go run appuser.go
   ```

2. **重新登录获取新 Token**
   ```bash
   curl -X POST http://localhost:10001/api/v1/auth/login \
     -H "Content-Type: application/json" \
     -d '{"phone":"13452552349","password":"13452552349"}'
   ```

3. **使用新 Token 访问受保护接口**
   ```bash
   curl -X GET http://localhost:10001/api/v1/user/info \
     -H "Authorization: Bearer <new_token>" \
     -H "Content-Type: application/json" \
     -d '{"phone":"13452552349"}'
   ```

## 预期结果
- ✅ JWT Token 过期错误消失
- ✅ 所有微服务使用统一的 JWT 密钥
- ✅ Token 生成和验证正常工作 