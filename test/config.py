# -*- coding: utf-8 -*-
"""
API测试配置文件
"""

# 服务配置
APPUSER_BASE_URL = "http://localhost:10001"      # C端用户服务地址
OAUSER_BASE_URL = "http://localhost:10002"       # B端用户服务地址
LEASEPRODUCT_BASE_URL = "http://localhost:10006" # 租赁产品服务地址 (修复: 10003 -> 10006)
LEASE_BASE_URL = "http://localhost:10004"        # 租赁业务服务地址
LOANPRODUCT_BASE_URL = "http://localhost:10005"  # 贷款产品服务地址
LOAN_BASE_URL = "http://localhost:10003"         # 贷款业务服务地址 (修复: 10006 -> 10003)

# C端测试用户配置
APPUSER_TEST_DATA = {
    "phone": "13452552490",
    "password": "13452552490"
}

# B端测试用户配置
OAUSER_TEST_DATA = [
    {
        "phone": "13452552490",
        "password": "13452552490", 
        "role": "admin"
    },
    {
        "phone": "13452552491", 
        "password": "13452552491",
        "role": "operator"
    }
]

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

# 租赁申请测试数据 - 修复日期问题：使用未来日期
LEASE_TEST_DATA = {
    "product_id": 1,  # 将根据实际创建的产品ID更新
    "product_code": "LP001",
    "name": "测试租赁申请",
    "type": "挖掘机",
    "machinery": "大型挖掘机",
    "start_date": "2025-07-01",  # 修改为未来日期
    "end_date": "2025-07-10",    # 修改为未来日期
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
    "product_id": 1,  # 将根据实际创建的产品ID更新
    "name": "测试贷款申请",
    "type": "信用贷款",
    "amount": 100000.00,
    "duration": 12,
    "purpose": "个人消费"
}

# 请求配置
REQUEST_TIMEOUT = 30  # 请求超时时间（秒）
REQUEST_DELAY = 1     # 请求间隔时间（秒） 