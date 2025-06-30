# -*- coding: utf-8 -*-
"""
API测试配置文件
"""

# 服务配置
APPUSER_BASE_URL = "http://localhost:10001"  # C端用户服务地址
OAUSER_BASE_URL = "http://localhost:10002"   # B端用户服务地址

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

# 请求配置
REQUEST_TIMEOUT = 30  # 请求超时时间（秒）
REQUEST_DELAY = 1     # 请求间隔时间（秒） 