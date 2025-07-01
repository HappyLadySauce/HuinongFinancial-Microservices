 # C端前端微服务重构文档

## 概述

本文档记录了C端前端（ui/users）从单体后端迁移到微服务后端的重构过程和主要变更。

## 主要变更

### 1. 认证机制变更

**旧版本（Cookie/Session认证）:**
- 使用 Cookie 进行认证
- 后端返回 `session_id`
- 依赖服务器端会话管理

**新版本（JWT Token认证）:**
- 使用 JWT Token 进行认证
- 后端返回 `token` 字段
- 客户端自主管理token状态
- 请求头添加 `Authorization: Bearer <token>`

### 2. API端点变更

**API基础地址:**
- 旧: `/api/v1`
- 新: `http://127.0.0.1:8080/api/v1`

**用户认证接口:**
- 注册: `/auth/register` (新增)
- 登录: `/auth/login` 
- 登出: `/auth/logout`
- 修改密码: `/auth/password` (新增)

**用户信息接口:**
- 获取用户信息: `/user/info` (GET方法修改为需要body参数)  
- 更新用户信息: `/user/info` (PUT)
- 删除用户: `/user/delete`

**业务接口:**
- 租赁申请: `/lease/applications`
- 贷款申请: `/loan/applications`
- 租赁产品: `/leaseproduct/products`
- 贷款产品: `/loanproduct/products`

### 3. 响应格式变更

**登录响应:**
```typescript
// 旧格式
{
  code: 200,
  message: "登录成功",
  data: { session_id: "xxx" }
}

// 新格式
{
  token: "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
}
```

**其他接口响应:**
大部分接口保持原有的 `{ code, message, data }` 格式不变。

### 4. 数据结构变更

**租赁申请 (LeaseApplication):**
```typescript
// 新增字段
{
  application_id: string,      // 申请编号
  product_code: string,        // 产品编码  
  machinery: string,           // 机械设备
  daily_rate: number,          // 日租金
  total_amount: number,        // 总金额
  deposit: number,             // 押金
  delivery_address: string,    // 交付地址
  contact_phone: string,       // 联系电话
  created_at: number,          // 时间戳格式
  updated_at: number
}
```

**贷款申请 (LoanApplication):**
```typescript
// 新增字段
{
  application_id: string,      // 申请编号
  created_at: number,          // 时间戳格式
  updated_at: number
}
```

## 文件变更清单

### 修改的文件

1. **src/services/api.ts** - 完全重构
   - 更新API基础地址
   - 重构认证机制为JWT
   - 更新所有接口定义
   - 新增产品相关API

2. **src/services/auth.ts** - 重构认证服务
   - 移除Session相关方法
   - 新增JWT token处理方法
   - 更新认证流程

3. **src/stores/user.ts** - 更新用户状态管理
   - JWT token存储和验证
   - token过期检查
   - 状态初始化优化

## 兼容性处理

为了保持与现有代码的兼容性，提供了以下类型别名：

```typescript
// 类型别名
export type LoanApproval = LoanApplication
export type LeaseApproval = LeaseApplication  
export type LoanApprovalRequest = LoanApplicationRequest
export type LeaseApprovalRequest = LeaseApplicationRequest

// API别名
export const loanApi = loanApprovalApi
export const leaseApi = leaseApprovalApi
```

## 使用说明

### 认证流程

```typescript
import { authApi } from '@/services/api'
import { useUserStore } from '@/stores/user'

// 登录
const loginResult = await authApi.login(phone, password)
const userStore = useUserStore()
userStore.login(loginResult) // 自动存储JWT token

// 后续请求会自动添加Authorization头
```

### 业务接口使用

```typescript
import { loanApprovalApi, leaseApprovalApi } from '@/services/api'

// 创建贷款申请
const loanResult = await loanApprovalApi.create({
  product_id: 1,
  name: "个人消费贷款",
  type: "信用贷款",
  amount: 100000,
  duration: 12,
  purpose: "购车"
})

// 创建租赁申请  
const leaseResult = await leaseApprovalApi.create({
  product_id: 1,
  product_code: "CAT001",
  name: "挖掘机租赁",
  type: "挖掘机",
  machinery: "大型挖掘机",
  start_date: "2025-01-01",
  end_date: "2025-01-10", 
  duration: 10,
  daily_rate: 800,
  total_amount: 8000,
  deposit: 10000,
  delivery_address: "北京市朝阳区",
  contact_phone: "13800138000",
  purpose: "工程施工"
})
```

## 注意事项

1. **JWT Token管理**: 前端需要妥善管理token的存储和过期处理
2. **错误处理**: 401错误会自动清除token并跳转到登录页
3. **API兼容性**: 部分旧接口可能需要调整参数格式
4. **时间格式**: 新API使用Unix时间戳而非ISO字符串

## 后续工作

1. 更新前端页面组件以使用新的API接口
2. 测试所有业务流程确保功能正常
3. 优化错误处理和用户体验
4. 添加更多的token验证和刷新机制