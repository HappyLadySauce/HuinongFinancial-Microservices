#!/usr/bin/env python3
# -*- coding: utf-8 -*-
"""
优化后的完整API测试脚本
解决数据依赖和测试顺序问题：
1. 分离CRUD测试和业务流程测试  
2. 保留必要的测试数据
3. 优化测试顺序
4. 添加数据准备和清理逻辑
"""

import requests
import json
import time
from typing import Optional, Dict, Any
from config import (
    APPUSER_BASE_URL, OAUSER_BASE_URL, LEASEPRODUCT_BASE_URL, LEASE_BASE_URL, 
    LOANPRODUCT_BASE_URL, LOAN_BASE_URL, APPUSER_TEST_DATA, OAUSER_TEST_DATA,
    LEASEPRODUCT_TEST_DATA, LEASE_TEST_DATA, LOANPRODUCT_TEST_DATA, LOAN_TEST_DATA,
    REQUEST_TIMEOUT, REQUEST_DELAY
)

class APITester:
    def __init__(self, base_url: str, service_name: str):
        self.base_url = base_url.rstrip('/')
        self.service_name = service_name
        self.session = requests.Session()
        self.token = None
        
    def set_token(self, token: str):
        """设置JWT token"""
        self.token = token
        self.session.headers.update({'Authorization': f'Bearer {token}'})
        
    def clear_token(self):
        """清除token"""
        self.token = None
        if 'Authorization' in self.session.headers:
            del self.session.headers['Authorization']
    
    def make_request(self, method: str, endpoint: str, data: Optional[Dict] = None, params: Optional[Dict] = None) -> Dict[Any, Any]:
        """发送HTTP请求"""
        url = f"{self.base_url}{endpoint}"
        
        print(f"\n🔄 [{self.service_name}] {method.upper()} {endpoint}")
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

    def register_or_login(self, phone: str, password: str, role: str = None) -> bool:
        """尝试注册，如果用户已存在则登录"""
        print(f"\n{'='*50}")
        print(f"🔐 尝试注册或登录 - {phone}")
        print(f"{'='*50}")
        
        # 先尝试注册
        data = {"phone": phone, "password": password}
        if role:
            data["role"] = role
            
        result = self.make_request("POST", "/api/v1/auth/register", data)
        
        if "token" in result:
            self.set_token(result["token"])
            print(f"✅ 注册成功，已设置token")
            return True
        elif "用户已存在" in str(result) or "already exists" in str(result).lower():
            print(f"ℹ️ 用户已存在，尝试登录...")
            return self.test_login(phone, password)
        else:
            print(f"❌ 注册失败: {result}")
            return False
    
    def test_login(self, phone: str, password: str) -> bool:
        """测试登录接口"""
        print(f"\n{'='*50}")
        print(f"🔐 测试登录 - {phone}")
        print(f"{'='*50}")
        
        data = {"phone": phone, "password": password}
        result = self.make_request("POST", "/api/v1/auth/login", data)
        
        if "token" in result:
            self.set_token(result["token"])
            print(f"✅ 登录成功，已设置token")
            return True
        else:
            print(f"❌ 登录失败")
            return False

# 全局变量保存测试状态
test_data = {
    "lease_product_id": None,
    "lease_product_code": None,
    "loan_product_id": None,
    "lease_application_id": None,
    "loan_application_id": None,
    "admin_token": None,
    "user_token": None
}

def setup_test_data():
    """数据准备阶段：创建测试所需的基础数据"""
    print(f"\n{'#'*60}")
    print(f"🏗️ 数据准备阶段：创建测试基础数据")
    print(f"{'#'*60}")
    
    # 1. 登录管理员账户
    admin_tester = APITester(OAUSER_BASE_URL, "OAUser")
    admin_user = OAUSER_TEST_DATA[0]
    
    if not admin_tester.register_or_login(admin_user["phone"], admin_user["password"], admin_user["role"]):
        print("❌ 管理员登录失败")
        return False
    
    test_data["admin_token"] = admin_tester.token
    
    # 2. 创建租赁产品（供后续测试使用）
    print(f"\n{'='*50}")
    print(f"🏗️ 创建租赁测试产品")
    print(f"{'='*50}")
    
    leaseproduct_tester = APITester(LEASEPRODUCT_BASE_URL, "LeaseProduct")
    leaseproduct_tester.set_token(admin_tester.token)
    
    result = leaseproduct_tester.make_request("POST", "/api/v1/admin/leaseproduct/products", LEASEPRODUCT_TEST_DATA)
    if "data" in result and "id" in result["data"]:
        test_data["lease_product_id"] = result["data"]["id"]
        test_data["lease_product_code"] = LEASEPRODUCT_TEST_DATA["product_code"]
        print(f"✅ 租赁产品创建成功，ID: {test_data['lease_product_id']}")
    else:
        print(f"❌ 租赁产品创建失败")
        return False
    
    time.sleep(REQUEST_DELAY)
    
    # 3. 创建贷款产品（供后续测试使用）
    print(f"\n{'='*50}")
    print(f"🏗️ 创建贷款测试产品")
    print(f"{'='*50}")
    
    loanproduct_tester = APITester(LOANPRODUCT_BASE_URL, "LoanProduct")
    loanproduct_tester.set_token(admin_tester.token)
    
    result = loanproduct_tester.make_request("POST", "/api/v1/admin/loanproduct/products", LOANPRODUCT_TEST_DATA)
    if "data" in result and "id" in result["data"]:
        test_data["loan_product_id"] = result["data"]["id"]
        print(f"✅ 贷款产品创建成功，ID: {test_data['loan_product_id']}")
    else:
        print(f"❌ 贷款产品创建失败")
        return False
    
    time.sleep(REQUEST_DELAY)
    
    # 4. 登录用户账户
    user_tester = APITester(APPUSER_BASE_URL, "AppUser")
    user_data = APPUSER_TEST_DATA
    
    if user_tester.register_or_login(user_data["phone"], user_data["password"]):
        test_data["user_token"] = user_tester.token
        print(f"✅ 用户登录成功")
    else:
        print(f"❌ 用户登录失败")
        return False
    
    print(f"✅ 数据准备阶段完成")
    return True

def test_user_services():
    """测试用户服务（基础功能）"""
    print(f"\n{'#'*60}")
    print(f"🌟 阶段1：测试用户服务")
    print(f"{'#'*60}")
    
    # 测试C端用户服务
    print(f"\n{'*'*40}")
    print(f"测试C端用户服务 (appuser)")
    print(f"{'*'*40}")
    
    user_tester = APITester(APPUSER_BASE_URL, "AppUser")
    user_tester.set_token(test_data["user_token"])
    
    phone = APPUSER_TEST_DATA["phone"]
    user_info_result = user_tester.make_request("GET", "/api/v1/user/info", {"phone": phone})
    
    if "user_info" in user_info_result:
        updated_user_info = user_info_result["user_info"].copy()
        updated_user_info.update({
            "name": "测试用户更新",
            "nickname": "测试昵称",
            "age": 25,
            "gender": 1,
            "occupation": "软件工程师",
            "address": "北京市朝阳区",
            "income": 15000.00
        })
        # 移除不需要更新的字段
        for field in ["status", "role"]:
            if field in updated_user_info:
                del updated_user_info[field]
        
        user_tester.make_request("PUT", "/api/v1/user/info", {"user_info": updated_user_info})
    
    time.sleep(REQUEST_DELAY)
    
    # 测试B端用户服务
    print(f"\n{'*'*40}")
    print(f"测试B端用户服务 (oauser)")
    print(f"{'*'*40}")
    
    for i, user in enumerate(OAUSER_TEST_DATA[:2]):  # 只测试前2个用户，减少时间
        admin_tester = APITester(OAUSER_BASE_URL, "OAUser")
        if admin_tester.register_or_login(user["phone"], user["password"], user["role"]):
            admin_tester.make_request("GET", "/api/v1/user/info", {"phone": user["phone"]})
            admin_tester.make_request("PUT", "/api/v1/user/status", {"phone": user["phone"], "status": 1})
            if i == 0:  # 保存第一个管理员token
                test_data["admin_token"] = admin_tester.token
        time.sleep(REQUEST_DELAY)

def test_product_services():
    """测试产品服务（CRUD功能，但保留测试数据）"""
    print(f"\n{'#'*60}")
    print(f"🌟 阶段2：测试产品管理服务")
    print(f"{'#'*60}")
    
    # 测试租赁产品服务
    print(f"\n{'*'*40}")
    print(f"测试租赁产品服务 (leaseproduct)")
    print(f"{'*'*40}")
    
    leaseproduct_tester = APITester(LEASEPRODUCT_BASE_URL, "LeaseProduct")
    leaseproduct_tester.set_token(test_data["admin_token"])
    
    # 查询产品列表
    leaseproduct_tester.make_request("GET", "/api/v1/leaseproduct/products", params={"page": 1, "size": 10})
    time.sleep(REQUEST_DELAY)
    
    # 查询产品详情
    leaseproduct_tester.make_request("GET", f"/api/v1/leaseproduct/products/{test_data['lease_product_code']}")
    time.sleep(REQUEST_DELAY)
    
    # 更新产品（保留产品，只测试更新功能）
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
    leaseproduct_tester.make_request("PUT", f"/api/v1/admin/leaseproduct/products/{test_data['lease_product_code']}", update_data)
    time.sleep(REQUEST_DELAY)
    
    # 测试库存检查（使用未来日期）
    inventory_check = {
        "product_code": test_data["lease_product_code"],
        "quantity": 2,
        "start_date": "2025-08-01",
        "end_date": "2025-08-10"
    }
    leaseproduct_tester.make_request("POST", "/api/v1/leaseproduct/products/check-inventory", inventory_check)
    time.sleep(REQUEST_DELAY)
    
    # 测试贷款产品服务
    print(f"\n{'*'*40}")
    print(f"测试贷款产品服务 (loanproduct)")
    print(f"{'*'*40}")
    
    loanproduct_tester = APITester(LOANPRODUCT_BASE_URL, "LoanProduct")
    loanproduct_tester.set_token(test_data["admin_token"])
    
    # 查询产品列表
    loanproduct_tester.make_request("GET", "/api/v1/loanproduct/products", params={"page": 1, "size": 10})
    time.sleep(REQUEST_DELAY)
    
    # 查询产品详情
    loanproduct_tester.make_request("GET", f"/api/v1/loanproduct/products/{test_data['loan_product_id']}")
    time.sleep(REQUEST_DELAY)
    
    # 更新产品（保留产品，只测试更新功能）
    update_data = {
        "name": "更新后的信用贷款",
        "type": "信用贷款",
        "max_amount": 600000.00,
        "min_amount": 5000.00,
        "max_duration": 48,
        "min_duration": 3,
        "interest_rate": 0.085,
        "description": "更新后的产品描述"
    }
    loanproduct_tester.make_request("PUT", f"/api/v1/admin/loanproduct/products/{test_data['loan_product_id']}", update_data)
    time.sleep(REQUEST_DELAY)

def test_business_services():
    """测试业务服务（使用已准备的测试数据）"""
    print(f"\n{'#'*60}")
    print(f"🌟 阶段3：测试业务流程服务")
    print(f"{'#'*60}")
    
    # 测试租赁业务服务
    print(f"\n{'*'*40}")
    print(f"测试租赁业务服务 (lease)")
    print(f"{'*'*40}")
    
    lease_tester = APITester(LEASE_BASE_URL, "Lease")
    lease_tester.set_token(test_data["user_token"])
    
    # 创建租赁申请
    lease_data = LEASE_TEST_DATA.copy()
    lease_data["product_code"] = test_data["lease_product_code"]  # 使用已创建的产品
    
    result = lease_tester.make_request("POST", "/api/v1/lease/applications", lease_data)
    if "application_id" in result:
        test_data["lease_application_id"] = result["application_id"]
        print(f"✅ 租赁申请创建成功，ID: {test_data['lease_application_id']}")
    else:
        print(f"❌ 租赁申请创建失败")
        return False
    
    time.sleep(REQUEST_DELAY)
    
    # 查询我的申请列表
    lease_tester.make_request("GET", "/api/v1/lease/applications", params={"page": 1, "size": 10})
    time.sleep(REQUEST_DELAY)
    
    # 查询申请详情
    lease_tester.make_request("GET", f"/api/v1/lease/applications/{test_data['lease_application_id']}")
    time.sleep(REQUEST_DELAY)
    
    # 更新申请
    update_data = {
        "purpose": "更新后的使用目的",
        "delivery_address": "更新后的地址",
        "contact_phone": "13900139000"
    }
    lease_tester.make_request("PUT", f"/api/v1/lease/applications/{test_data['lease_application_id']}", update_data)
    time.sleep(REQUEST_DELAY)
    
    # 管理员审批申请
    lease_tester.set_token(test_data["admin_token"])
    lease_tester.make_request("GET", "/api/v1/admin/lease/applications", params={"page": 1, "size": 10})
    
    approval_data = {
        "action": "approve",
        "suggestions": "申请已通过审核",
        "approved_duration": 10,
        "approved_amount": 8000.00,
        "approved_deposit": 10000.00
    }
    lease_tester.make_request("POST", f"/api/v1/admin/lease/applications/{test_data['lease_application_id']}/approve", approval_data)
    time.sleep(REQUEST_DELAY)
    
    # 测试贷款业务服务
    print(f"\n{'*'*40}")
    print(f"测试贷款业务服务 (loan)")
    print(f"{'*'*40}")
    
    loan_tester = APITester(LOAN_BASE_URL, "Loan")
    loan_tester.set_token(test_data["user_token"])
    
    # 创建贷款申请
    loan_data = LOAN_TEST_DATA.copy()
    loan_data["product_id"] = test_data["loan_product_id"]  # 使用已创建的产品
    
    result = loan_tester.make_request("POST", "/api/v1/loan/applications", loan_data)
    if "application_id" in result:
        test_data["loan_application_id"] = result["application_id"]
        print(f"✅ 贷款申请创建成功，ID: {test_data['loan_application_id']}")
    else:
        print(f"❌ 贷款申请创建失败")
        return False
    
    time.sleep(REQUEST_DELAY)
    
    # 查询我的申请列表
    loan_tester.make_request("GET", "/api/v1/loan/applications", params={"page": 1, "size": 10})
    time.sleep(REQUEST_DELAY)
    
    # 查询申请详情
    loan_tester.make_request("GET", f"/api/v1/loan/applications/{test_data['loan_application_id']}")
    time.sleep(REQUEST_DELAY)
    
    # 更新申请
    update_data = {
        "amount": 120000.00,
        "duration": 18,
        "purpose": "更新后的贷款用途"
    }
    loan_tester.make_request("PUT", f"/api/v1/loan/applications/{test_data['loan_application_id']}", update_data)
    time.sleep(REQUEST_DELAY)
    
    # 管理员审批申请
    loan_tester.set_token(test_data["admin_token"])
    loan_tester.make_request("GET", "/api/v1/admin/loan/applications", params={"page": 1, "size": 10})
    
    approval_data = {
        "action": "approve",
        "suggestions": "申请已通过审核",
        "approved_amount": 100000.00,
        "approved_duration": 12,
        "interest_rate": 0.08
    }
    loan_tester.make_request("POST", f"/api/v1/admin/loan/applications/{test_data['loan_application_id']}/approve", approval_data)
    time.sleep(REQUEST_DELAY)
    
    return True

def cleanup_test_data():
    """清理测试数据（可选）"""
    print(f"\n{'#'*60}")
    print(f"🧹 清理测试数据")
    print(f"{'#'*60}")
    
    # 注意：在实际环境中可能不需要清理，保留数据供后续测试使用
    # 这里只是演示如何清理
    
    if test_data["admin_token"]:
        # 清理租赁产品
        if test_data["lease_product_code"]:
            leaseproduct_tester = APITester(LEASEPRODUCT_BASE_URL, "LeaseProduct")
            leaseproduct_tester.set_token(test_data["admin_token"])
            leaseproduct_tester.make_request("DELETE", f"/api/v1/admin/leaseproduct/products/{test_data['lease_product_code']}")
            print(f"✅ 清理租赁产品：{test_data['lease_product_code']}")
        
        # 清理贷款产品
        if test_data["loan_product_id"]:
            loanproduct_tester = APITester(LOANPRODUCT_BASE_URL, "LoanProduct")
            loanproduct_tester.set_token(test_data["admin_token"])
            loanproduct_tester.make_request("DELETE", f"/api/v1/admin/loanproduct/products/{test_data['loan_product_id']}")
            print(f"✅ 清理贷款产品：{test_data['loan_product_id']}")

def main():
    """主函数"""
    print("🚀 开始优化后的完整API测试")
    print("注意：请确保相关服务正在运行")
    print(f"- appuser服务: {APPUSER_BASE_URL}")
    print(f"- oauser服务: {OAUSER_BASE_URL}")
    print(f"- leaseproduct服务: {LEASEPRODUCT_BASE_URL}")
    print(f"- lease服务: {LEASE_BASE_URL}")
    print(f"- loanproduct服务: {LOANPRODUCT_BASE_URL}")
    print(f"- loan服务: {LOAN_BASE_URL}")
    
    try:
        # 阶段0：数据准备
        if not setup_test_data():
            print("❌ 数据准备失败，终止测试")
            return
        
        # 阶段1：测试用户服务
        test_user_services()
        time.sleep(2)
        
        # 阶段2：测试产品服务
        test_product_services()
        time.sleep(2)
        
        # 阶段3：测试业务服务
        if test_business_services():
            print(f"\n✅ 所有业务流程测试成功")
        
        # 阶段4：清理数据（可选）
        # cleanup_test_data()  # 取消注释以启用清理
        
        print(f"\n{'='*60}")
        print("🎉 所有测试完成！")
        print(f"{'='*60}")
        
    except KeyboardInterrupt:
        print("\n⚠️ 测试被用户中断")
    except Exception as e:
        print(f"\n❌ 测试过程中出现错误: {str(e)}")

if __name__ == "__main__":
    main() 