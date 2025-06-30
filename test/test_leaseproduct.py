#!/usr/bin/env python3
# -*- coding: utf-8 -*-
"""
租赁产品服务API测试脚本
"""

import requests
import json
import time
from typing import Optional, Dict, Any
from config import LEASEPRODUCT_BASE_URL, OAUSER_BASE_URL, LEASEPRODUCT_TEST_DATA, OAUSER_TEST_DATA, REQUEST_TIMEOUT, REQUEST_DELAY

class LeaseProductTester:
    def __init__(self, base_url: str):
        self.base_url = base_url.rstrip('/')
        self.session = requests.Session()
        self.token = None
        
    def set_token(self, token: str):
        """设置JWT token"""
        self.token = token
        self.session.headers.update({'Authorization': f'Bearer {token}'})
        
    def make_request(self, method: str, endpoint: str, data: Optional[Dict] = None, params: Optional[Dict] = None) -> Dict[Any, Any]:
        """发送HTTP请求"""
        url = f"{self.base_url}{endpoint}"
        
        print(f"\n🔄 [LeaseProduct] {method.upper()} {endpoint}")
        if data:
            print(f"📤 Request: {json.dumps(data, indent=2, ensure_ascii=False)}")
        if params:
            print(f"📋 Params: {json.dumps(params, indent=2, ensure_ascii=False)}")
        
        try:
            if method.lower() == 'get':
                response = self.session.get(url, json=data, params=params, timeout=REQUEST_TIMEOUT)
            elif method.lower() == 'post':
                response = self.session.post(url, json=data, timeout=REQUEST_TIMEOUT)
            elif method.lower() == 'put':
                response = self.session.put(url, json=data, timeout=REQUEST_TIMEOUT)
            elif method.lower() == 'delete':
                response = self.session.delete(url, json=data, timeout=REQUEST_TIMEOUT)
            else:
                raise ValueError(f"Unsupported method: {method}")
            
            print(f"📊 Status: {response.status_code}")
            
            if response.headers.get('content-type', '').startswith('application/json'):
                result = response.json()
                print(f"📥 Response: {json.dumps(result, indent=2, ensure_ascii=False)}")
                return result
            else:
                print(f"📥 Response: {response.text}")
                return {"status_code": response.status_code, "text": response.text}
                
        except Exception as e:
            print(f"❌ Error: {str(e)}")
            return {"error": str(e)}

    def test_create_product(self, product_data: Dict) -> Optional[int]:
        """测试创建租赁产品"""
        print(f"\n{'='*50}")
        print(f"🏗️ 测试创建租赁产品")
        print(f"{'='*50}")
        
        result = self.make_request("POST", "/api/v1/admin/leaseproduct/products", product_data)
        
        if "data" in result and "id" in result["data"]:
            product_id = result["data"]["id"]
            print(f"✅ 产品创建成功，ID: {product_id}")
            return product_id
        else:
            print(f"❌ 产品创建失败")
            return None
    
    def test_list_products(self, params: Dict = None):
        """测试获取产品列表"""
        print(f"\n{'='*50}")
        print(f"📋 测试获取产品列表")
        print(f"{'='*50}")
        
        default_params = {"page": 1, "size": 10}
        if params:
            default_params.update(params)
            
        result = self.make_request("GET", "/api/v1/leaseproduct/products", params=default_params)
        return result
    
    def test_get_product_detail(self, product_code: str):
        """测试获取产品详情"""
        print(f"\n{'='*50}")
        print(f"🔍 测试获取产品详情 - {product_code}")
        print(f"{'='*50}")
        
        result = self.make_request("GET", f"/api/v1/leaseproduct/products/{product_code}")
        return result
    
    def test_update_product(self, product_code: str, update_data: Dict):
        """测试更新产品"""
        print(f"\n{'='*50}")
        print(f"✏️ 测试更新产品 - {product_code}")
        print(f"{'='*50}")
        
        result = self.make_request("PUT", f"/api/v1/admin/leaseproduct/products/{product_code}", update_data)
        return result
    
    def test_check_inventory(self, inventory_data: Dict):
        """测试检查库存可用性"""
        print(f"\n{'='*50}")    
        print(f"📦 测试检查库存可用性")
        print(f"{'='*50}")
        
        result = self.make_request("POST", "/api/v1/leaseproduct/products/check-inventory", inventory_data)
        return result
    
    def test_delete_product(self, product_code: str):
        """测试删除产品"""
        print(f"\n{'='*50}")
        print(f"🗑️ 测试删除产品 - {product_code}")
        print(f"{'='*50}")
        
        result = self.make_request("DELETE", f"/api/v1/admin/leaseproduct/products/{product_code}")
        return result

def login_admin():
    """登录管理员账户获取token"""
    print("🔐 正在登录管理员账户...")
    
    admin_user = OAUSER_TEST_DATA[0]
    login_data = {
        "phone": admin_user["phone"],
        "password": admin_user["password"]
    }
    
    try:
        response = requests.post(
            f"{OAUSER_BASE_URL}/api/v1/auth/login",
            json=login_data,
            timeout=REQUEST_TIMEOUT
        )
        
        if response.status_code == 200:
            result = response.json()
            if "token" in result:
                print("✅ 管理员登录成功")
                return result["token"]
        
        print("❌ 管理员登录失败")
        return None
        
    except Exception as e:
        print(f"❌ 登录出错: {str(e)}")
        return None

def test_leaseproduct_service():
    """测试租赁产品服务完整流程"""
    print(f"\n{'#'*60}")
    print(f"🌟 开始测试 租赁产品服务 (leaseproduct)")
    print(f"🔍 服务地址: {LEASEPRODUCT_BASE_URL}")
    print(f"{'#'*60}")
    
    # 获取管理员token
    token = login_admin()
    if not token:
        print("❌ 无法获取管理员token，测试中止")
        return
    
    tester = LeaseProductTester(LEASEPRODUCT_BASE_URL)
    tester.set_token(token)
    
    try:
        # 1. 测试创建产品
        print("\n🔸 步骤1: 创建租赁产品")
        product_id = tester.test_create_product(LEASEPRODUCT_TEST_DATA)
        if not product_id:
            print("❌ 产品创建失败，跳过后续测试")
            return
        time.sleep(REQUEST_DELAY)
        
        # 2. 测试获取产品列表
        print("\n🔸 步骤2: 获取产品列表")
        tester.test_list_products()
        time.sleep(REQUEST_DELAY)
        
        # 3. 测试获取产品详情
        print("\n🔸 步骤3: 获取产品详情")
        tester.test_get_product_detail(LEASEPRODUCT_TEST_DATA["product_code"])
        time.sleep(REQUEST_DELAY)
        
        # 4. 测试更新产品
        print("\n🔸 步骤4: 更新产品信息")
        update_data = {
            "name": "更新后的挖掘机",
            "type": "挖掘机",
            "machinery": "大型挖掘机",
            "brand": "卡特彼勒",
            "model": "CAT320D",
            "daily_rate": 850.00,
            "deposit": 12000.00,
            "max_duration": 300,
            "min_duration": 3,
            "description": "更新后的产品描述",
            "status": 1
        }
        tester.test_update_product(LEASEPRODUCT_TEST_DATA["product_code"], update_data)
        time.sleep(REQUEST_DELAY)
        
        # 5. 测试检查库存
        print("\n🔸 步骤5: 检查库存可用性")
        inventory_check = {
            "product_code": LEASEPRODUCT_TEST_DATA["product_code"],
            "quantity": 2,
            "start_date": "2024-01-01",
            "end_date": "2024-01-10"
        }
        tester.test_check_inventory(inventory_check)
        time.sleep(REQUEST_DELAY)
        
        # 6. 测试产品筛选查询
        print("\n🔸 步骤6: 测试产品筛选查询")
        filter_params = {
            "type": "挖掘机",
            "brand": "卡特彼勒",
            "status": 1,
            "keyword": "挖掘机"
        }
        tester.test_list_products(filter_params)
        time.sleep(REQUEST_DELAY)
        
        # 7. 测试删除产品
        print("\n🔸 步骤7: 删除产品")
        tester.test_delete_product(LEASEPRODUCT_TEST_DATA["product_code"])
        
        print(f"\n{'='*60}")
        print("✅ 租赁产品服务测试完成！")
        print(f"{'='*60}")
        
    except Exception as e:
        print(f"❌ 租赁产品测试过程中出现错误: {str(e)}")

if __name__ == "__main__":
    test_leaseproduct_service() 