# 微服务API测试脚本完善指南

## 概述
当前的 `test_api.py` 只测试了用户服务（appuser 和 oauser），需要完善其他4个微服务模块的测试：

1. **租赁产品服务** (leaseproduct)
2. **租赁业务服务** (lease)
3. **贷款产品服务** (loanproduct)  
4. **贷款业务服务** (loan)

## 配置文件更新

### 1. 更新 `config.py`
在配置文件中添加新服务的配置：

```python
# 新增服务配置
LEASEPRODUCT_BASE_URL = "http://localhost:10003"  # 租赁产品服务地址
LEASE_BASE_URL = "http://localhost:10004"         # 租赁业务服务地址
LOANPRODUCT_BASE_URL = "http://localhost:10005"   # 贷款产品服务地址
LOAN_BASE_URL = "http://localhost:10006"          # 贷款业务服务地址

# 租赁产品测试数据
LEASEPRODUCT_TEST_DATA = {
    "product_code": "LP001",
    "name": "挖掘机测试产品",
    "type": "挖掘机",
    "machinery": "大型挖掘机",
    "brand": "卡特彼勒",
    "model": "CAT320D",
    "daily_rate": 800.00,
    "deposit": 10000.00,
    "max_duration": 365,
    "min_duration": 1,
    "description": "高性能挖掘机，适用于大型工程项目",
    "inventory_count": 5
}

# 租赁申请测试数据
LEASE_TEST_DATA = {
    "product_id": 1,
    "product_code": "LP001",
    "name": "测试租赁申请",
    "type": "挖掘机",
    "machinery": "大型挖掘机",
    "start_date": "2024-01-01",
    "end_date": "2024-01-10",
    "duration": 10,
    "daily_rate": 800.00,
    "total_amount": 8000.00,
    "deposit": 10000.00,
    "delivery_address": "北京市朝阳区建国路1号",
    "contact_phone": "13800138000",
    "purpose": "建筑工程施工"
}

# 贷款产品测试数据
LOANPRODUCT_TEST_DATA = {
    "product_code": "LN001",
    "name": "个人信用贷款",
    "type": "信用贷款",
    "max_amount": 500000.00,
    "min_amount": 10000.00,
    "max_duration": 36,
    "min_duration": 6,
    "interest_rate": 0.08,
    "description": "无抵押个人信用贷款产品"
}

# 贷款申请测试数据
LOAN_TEST_DATA = {
    "product_id": 1,
    "name": "测试贷款申请",
    "type": "信用贷款",
    "amount": 100000.00,
    "duration": 12,
    "purpose": "个人消费"
}
```

## 测试类扩展

### 2. 租赁产品服务测试方法

```python
def test_leaseproduct_service():
    """测试租赁产品服务"""
    # 需要管理员权限
    # 测试接口：
    # - POST /api/v1/admin/leaseproduct/products (创建产品)
    # - GET /api/v1/leaseproduct/products (获取产品列表)
    # - GET /api/v1/leaseproduct/products/{productCode} (获取产品详情)
    # - PUT /api/v1/admin/leaseproduct/products/{productCode} (更新产品)
    # - DELETE /api/v1/admin/leaseproduct/products/{productCode} (删除产品)
    # - POST /api/v1/leaseproduct/products/check-inventory (检查库存)
```

### 3. 租赁业务服务测试方法

```python
def test_lease_service():
    """测试租赁业务服务"""
    # 需要用户登录 + 管理员审批
    # 测试接口：
    # - POST /api/v1/lease/applications (创建申请)
    # - GET /api/v1/lease/applications (获取我的申请列表)
    # - GET /api/v1/lease/applications/{id} (获取申请详情)
    # - PUT /api/v1/lease/applications/{id} (更新申请)
    # - POST /api/v1/lease/applications/{id}/cancel (取消申请)
    # - GET /api/v1/admin/lease/applications (管理员获取所有申请)
    # - POST /api/v1/admin/lease/applications/{id}/approve (管理员审批)
```

### 4. 贷款产品服务测试方法

```python
def test_loanproduct_service():
    """测试贷款产品服务"""
    # 需要管理员权限
    # 测试接口：
    # - POST /api/v1/admin/loanproduct/products (创建产品)
    # - GET /api/v1/loanproduct/products (获取产品列表)
    # - GET /api/v1/loanproduct/products/{id} (获取产品详情)
    # - PUT /api/v1/admin/loanproduct/products/{id} (更新产品)
    # - PUT /api/v1/admin/loanproduct/products/{id}/status (更新状态)
    # - DELETE /api/v1/admin/loanproduct/products/{id} (删除产品)
```

### 5. 贷款业务服务测试方法

```python
def test_loan_service():
    """测试贷款业务服务"""
    # 需要用户登录 + 管理员审批
    # 测试接口：
    # - POST /api/v1/loan/applications (创建申请)
    # - GET /api/v1/loan/applications (获取我的申请列表)
    # - GET /api/v1/loan/applications/{id} (获取申请详情)
    # - PUT /api/v1/loan/applications/{id} (更新申请)
    # - POST /api/v1/loan/applications/{id}/cancel (取消申请)
    # - GET /api/v1/admin/loan/applications (管理员获取所有申请)
    # - POST /api/v1/admin/loan/applications/{id}/approve (管理员审批)
```

## 测试流程设计

### 权限管理
- **产品管理接口**: 需要B端管理员token
- **业务申请接口**: 需要C端用户token
- **审批接口**: 需要B端管理员token

### 测试顺序
1. 用户服务测试（已完成）
2. 产品服务测试（租赁产品、贷款产品）
3. 业务服务测试（租赁申请、贷款申请）

### 数据依赖
- 业务申请测试依赖于产品数据
- 建议在测试中动态创建和清理测试数据

## 实现建议

### 1. 模块化设计
```python
class ProductTester(APITester):
    """产品服务测试基类"""
    
class BusinessTester(APITester): 
    """业务服务测试基类"""
```

### 2. 测试数据管理
- 使用测试前创建、测试后清理的策略
- 保存创建的资源ID，用于后续测试

### 3. 错误处理
- 添加完善的异常处理
- 记录测试结果和失败原因

### 4. 并发测试
- 考虑添加并发测试场景
- 测试系统在高并发下的表现

## 完整测试脚本结构

```
test/
├── config.py                    # 配置文件
├── test_api.py                 # 原有用户服务测试
├── test_leaseproduct.py        # 租赁产品服务测试
├── test_lease.py               # 租赁业务服务测试  
├── test_loanproduct.py         # 贷款产品服务测试
├── test_loan.py                # 贷款业务服务测试
├── test_all.py                 # 完整集成测试
└── README_TEST_ENHANCEMENT.md  # 本指导文档
```

## 运行方式

```bash
# 单独测试某个服务
python test_leaseproduct.py

# 完整测试所有服务
python test_all.py

# 原有用户服务测试
python test_api.py
```

## 注意事项

1. **服务依赖**: 确保所有微服务都在运行
2. **端口配置**: 根据实际部署调整服务端口
3. **数据清理**: 测试后清理测试数据
4. **权限验证**: 确保token有效性和权限正确性
5. **异常处理**: 处理网络超时、服务不可用等异常情况

这个指南提供了完善测试脚本的完整框架，可以根据实际需求进行调整和实现。 