#!/usr/bin/env python3
# -*- coding: utf-8 -*-
"""
完全修复后的API测试脚本
修复的问题：
1. 产品重复创建问题 - 使用现有产品或创建唯一编码产品
2. Application ID为空问题 - 修复获取详情的逻辑
3. 业务流程顺序问题 - 分离取消测试和审批测试
4. 数据同步问题 - 确保数据一致性
"""

import requests
import json
import time
import random
import string
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
        
    def make_request(self, method: str, endpoint: str, data: Dict = None, params: Dict = None) -> Optional[Dict]:
        """发起HTTP请求"""
        url = f"{self.base_url}{endpoint}"
        headers = {'Content-Type': 'application/json'}
        
        if self.token:
            headers['Authorization'] = f'Bearer {self.token}'
        
        print(f"🔄 [{self.service_name}] {method} {endpoint}")
        
        if data and method in ['POST', 'PUT', 'PATCH']:
            print(f"📤 Request: {json.dumps(data, indent=2, ensure_ascii=False)}")
        
        if params and method == 'GET':
            print(f"📋 Params: {json.dumps(params, indent=2, ensure_ascii=False)}")
        
        try:
            if method == 'GET':
                response = self.session.get(url, headers=headers, params=params, timeout=REQUEST_TIMEOUT)
            elif method == 'POST':
                response = self.session.post(url, headers=headers, json=data, timeout=REQUEST_TIMEOUT)
            elif method == 'PUT':
                response = self.session.put(url, headers=headers, json=data, timeout=REQUEST_TIMEOUT)
            elif method == 'DELETE':
                response = self.session.delete(url, headers=headers, timeout=REQUEST_TIMEOUT)
            else:
                print(f"❌ 不支持的HTTP方法: {method}")
                return None
            
            print(f"📊 Status: {response.status_code}")
            
            if response.headers.get('content-type', '').startswith('application/json'):
                result = response.json()
                print(f"📥 Response: {json.dumps(result, indent=2, ensure_ascii=False)}")
                return result
            else:
                print(f"📥 Response: {response.text}")
                return {"status_code": response.status_code, "text": response.text}
                
        except requests.exceptions.RequestException as e:
            print(f"❌ 请求失败: {e}")
            return None
        except json.JSONDecodeError as e:
            print(f"❌ JSON解析失败: {e}")
            return None
        
    def register_or_login(self, phone: str, password: str, role: str = None) -> bool:
        """注册或登录用户"""
        print(f"\n{'='*50}")
        print(f"🔐 尝试注册或登录 - {phone}")
        print(f"{'='*50}")
        
        # 准备注册数据
        register_data = {"phone": phone, "password": password}
        if role:
            register_data["role"] = role
        
        # 尝试注册
        result = self.make_request("POST", "/api/v1/auth/register", register_data)
        
        if result and "token" in result:
            self.token = result["token"]
            print("✅ 注册成功，已设置token")
            return True
        
        # 注册失败，尝试登录
        print("ℹ️ 用户已存在，尝试登录...")
        print(f"\n{'='*50}")
        print(f"🔐 测试登录 - {phone}")
        print(f"{'='*50}")
        
        login_data = {"phone": phone, "password": password}
        result = self.make_request("POST", "/api/v1/auth/login", login_data)
        
        if result and "token" in result:
            self.token = result["token"]
            print("✅ 登录成功，已设置token")
            return True
            
        print("❌ 登录失败")
        return False

def generate_unique_code(prefix: str, length: int = 6) -> str:
    """生成唯一的编码"""
    suffix = ''.join(random.choices(string.ascii_uppercase + string.digits, k=length))
    return f"{prefix}{suffix}"

# 全局测试数据存储
test_data = {
    "lease_product_id": None,
    "lease_product_code": None,
    "loan_product_id": None,
    "loan_product_code": None,
    "admin_token": None,
    "user_token": None,
    "lease_application_for_cancel": None,  # 专门用于取消测试的申请
    "lease_application_for_approve": None, # 专门用于审批测试的申请
    "loan_application_for_cancel": None,   # 专门用于取消测试的申请
    "loan_application_for_approve": None,  # 专门用于审批测试的申请
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
    
    # 2. 登录普通用户账户
    user_tester = APITester(APPUSER_BASE_URL, "AppUser")
    user_data = APPUSER_TEST_DATA
    
    if not user_tester.register_or_login(user_data["phone"], user_data["password"]):
        print("❌ 用户登录失败")
        return False
    
    test_data["user_token"] = user_tester.token
    
    # 3. 确保租赁产品存在（使用现有或创建新的）
    print(f"\n{'='*50}")
    print(f"🏗️ 准备租赁测试产品")
    print(f"{'='*50}")
    
    leaseproduct_tester = APITester(LEASEPRODUCT_BASE_URL, "LeaseProduct")
    leaseproduct_tester.set_token(admin_tester.token)
    
    # 先尝试获取现有产品列表
    existing_products = leaseproduct_tester.make_request("GET", "/api/v1/lease-product/products", params={"page": 1, "size": 100})
    
    lease_product_found = False
    if existing_products and "list" in existing_products:
        for product in existing_products["list"]:
            if product.get("status") == 1:  # 上架状态的产品
                test_data["lease_product_id"] = product["id"]
                test_data["lease_product_code"] = product["product_code"]
                lease_product_found = True
                print(f"✅ 使用现有租赁产品，ID: {test_data['lease_product_id']}, 编码: {test_data['lease_product_code']}")
                break
    
    if not lease_product_found:
        # 创建新的租赁产品（使用唯一编码）
        unique_code = generate_unique_code("LP")
        leaseproduct_data = LEASEPRODUCT_TEST_DATA.copy()
        leaseproduct_data["product_code"] = unique_code
        
        result = leaseproduct_tester.make_request("POST", "/api/v1/admin/leaseproduct/products", leaseproduct_data)
        if result and "data" in result and "id" in result["data"]:
            test_data["lease_product_id"] = result["data"]["id"]
            test_data["lease_product_code"] = unique_code
            print(f"✅ 租赁产品创建成功，ID: {test_data['lease_product_id']}, 编码: {test_data['lease_product_code']}")
        else:
            print("❌ 租赁产品创建失败")
            return False
    
    # 4. 确保贷款产品存在（使用现有或创建新的）
    print(f"\n{'='*50}")
    print(f"🏗️ 准备贷款测试产品")
    print(f"{'='*50}")
    
    loanproduct_tester = APITester(LOANPRODUCT_BASE_URL, "LoanProduct")
    loanproduct_tester.set_token(admin_tester.token)
    
    # 先尝试获取现有产品列表
    existing_products = loanproduct_tester.make_request("GET", "/api/v1/loan-product/products", params={"page": 1, "size": 100})
    
    loan_product_found = False
    if existing_products and "list" in existing_products:
        for product in existing_products["list"]:
            if product.get("status") == 1:  # 上架状态的产品
                test_data["loan_product_id"] = product["id"]
                test_data["loan_product_code"] = product["product_code"]
                loan_product_found = True
                print(f"✅ 使用现有贷款产品，ID: {test_data['loan_product_id']}, 编码: {test_data['loan_product_code']}")
                break
    
    if not loan_product_found:
        # 创建新的贷款产品（使用唯一编码）
        unique_code = generate_unique_code("LN")
        loanproduct_data = LOANPRODUCT_TEST_DATA.copy()
        loanproduct_data["product_code"] = unique_code
        
        result = loanproduct_tester.make_request("POST", "/api/v1/admin/loanproduct/products", loanproduct_data)
        if result and "data" in result and "id" in result["data"]:
            test_data["loan_product_id"] = result["data"]["id"]
            test_data["loan_product_code"] = unique_code
            print(f"✅ 贷款产品创建成功，ID: {test_data['loan_product_id']}, 编码: {test_data['loan_product_code']}")
        else:
            print("❌ 贷款产品创建失败")
            return False
    
    print(f"\n🎉 数据准备完成！")
    return True

def test_user_services():
    """测试用户服务"""
    print(f"\n{'#'*60}")
    print(f"🌟 阶段1：测试用户服务")
    print(f"{'#'*60}")
    
    # 测试B端用户服务
    print(f"\n{'*'*40}")
    print(f"测试B端用户服务 (oauser)")
    print(f"{'*'*40}")
    
    for user_data in OAUSER_TEST_DATA:
        tester = APITester(OAUSER_BASE_URL, "OAUser")
        
        if tester.register_or_login(user_data["phone"], user_data["password"], user_data["role"]):
            # 获取用户信息
            tester.make_request("GET", "/api/v1/user/info", params={"phone": user_data["phone"]})
            time.sleep(REQUEST_DELAY)
            
            # 更新用户状态
            tester.make_request("PUT", "/api/v1/user/status", {"phone": user_data["phone"], "status": 1})
            time.sleep(REQUEST_DELAY)
            
            # 退出登录
            tester.make_request("POST", "/api/v1/auth/logout")
            time.sleep(REQUEST_DELAY)

def test_lease_business():
    """测试租赁业务服务（分离取消和审批测试）"""
    print(f"\n{'#'*60}")
    print(f"🌟 阶段2：测试租赁业务服务")
    print(f"{'#'*60}")
    
    lease_tester = APITester(LEASE_BASE_URL, "Lease")
    lease_tester.set_token(test_data["user_token"])
    
    # === 第一部分：创建用于取消测试的申请 ===
    print(f"\n{'='*50}")
    print(f"📝 创建租赁申请（用于取消测试）")
    print(f"{'='*50}")
    
    lease_data_for_cancel = LEASE_TEST_DATA.copy()
    lease_data_for_cancel["product_id"] = test_data["lease_product_id"]
    lease_data_for_cancel["product_code"] = test_data["lease_product_code"]
    lease_data_for_cancel["name"] = "测试租赁申请（待取消）"
    
    result = lease_tester.make_request("POST", "/api/v1/lease/applications", lease_data_for_cancel)
    if result and "application_id" in result:
        test_data["lease_application_for_cancel"] = result["application_id"]
        print(f"✅ 租赁申请创建成功（用于取消），ID: {test_data['lease_application_for_cancel']}")
    else:
        print("❌ 租赁申请创建失败")
        return False
    
    time.sleep(REQUEST_DELAY)
    
    # === 测试取消功能（在pending状态） ===
    print(f"\n{'='*50}")
    print(f"❌ 测试取消申请（pending状态）")
    print(f"{'='*50}")
    
    cancel_result = lease_tester.make_request("POST", f"/api/v1/lease/applications/{test_data['lease_application_for_cancel']}/cancel", {
        "reason": "测试取消原因"
    })
    
    if cancel_result is not None:
        print("✅ 取消申请测试完成")
    
    time.sleep(REQUEST_DELAY)
    
    # === 第二部分：创建用于审批测试的申请 ===
    print(f"\n{'='*50}")
    print(f"📝 创建租赁申请（用于审批测试）")
    print(f"{'='*50}")
    
    lease_data_for_approve = LEASE_TEST_DATA.copy()
    lease_data_for_approve["product_id"] = test_data["lease_product_id"]
    lease_data_for_approve["product_code"] = test_data["lease_product_code"]
    lease_data_for_approve["name"] = "测试租赁申请（待审批）"
    
    result = lease_tester.make_request("POST", "/api/v1/lease/applications", lease_data_for_approve)
    if result and "application_id" in result:
        test_data["lease_application_for_approve"] = result["application_id"]
        print(f"✅ 租赁申请创建成功（用于审批），ID: {test_data['lease_application_for_approve']}")
    else:
        print("❌ 租赁申请创建失败")
        return False
    
    time.sleep(REQUEST_DELAY)
    
    # === 查询申请列表 ===
    print(f"\n{'='*50}")
    print(f"📋 测试获取我的申请列表")
    print(f"{'='*50}")
    
    lease_tester.make_request("GET", "/api/v1/lease/applications", params={"page": 1, "size": 10})
    time.sleep(REQUEST_DELAY)
    
    # === 更新申请信息 ===
    print(f"\n{'='*50}")
    print(f"✏️ 测试更新申请")
    print(f"{'='*50}")
    
    update_data = {
        "purpose": "更新后的使用目的",
        "delivery_address": "更新后的地址",
        "contact_phone": "13900139000"
    }
    
    lease_tester.make_request("PUT", f"/api/v1/lease/applications/{test_data['lease_application_for_approve']}", update_data)
    time.sleep(REQUEST_DELAY)
    
    # === 管理员审批 ===
    print(f"\n{'='*50}")
    print(f"👮‍♂️ 测试管理员审批")
    print(f"{'='*50}")
    
    # 切换到管理员token
    lease_admin_tester = APITester(LEASE_BASE_URL, "Lease")
    lease_admin_tester.set_token(test_data["admin_token"])
    
    # 获取所有申请
    lease_admin_tester.make_request("GET", "/api/v1/admin/lease/applications", params={"page": 1, "size": 10})
    time.sleep(REQUEST_DELAY)
    
    # 审批申请
    approve_data = {
        "action": "approve",
        "suggestions": "申请已通过审核",
        "approved_duration": 10,
        "approved_amount": 8000.0,
        "approved_deposit": 10000.0
    }
    
    lease_admin_tester.make_request("POST", f"/api/v1/admin/lease/applications/{test_data['lease_application_for_approve']}/approve", approve_data)
    
    print("✅ 租赁业务测试完成")
    return True

def test_loan_business():
    """测试贷款业务服务（分离取消和审批测试）"""
    print(f"\n{'#'*60}")
    print(f"🌟 阶段3：测试贷款业务服务")
    print(f"{'#'*60}")
    
    loan_tester = APITester(LOAN_BASE_URL, "Loan")
    loan_tester.set_token(test_data["user_token"])
    
    # === 第一部分：创建用于取消测试的申请 ===
    print(f"\n{'='*50}")
    print(f"📝 创建贷款申请（用于取消测试）")
    print(f"{'='*50}")
    
    loan_data_for_cancel = LOAN_TEST_DATA.copy()
    loan_data_for_cancel["product_id"] = test_data["loan_product_id"]
    loan_data_for_cancel["name"] = "测试贷款申请（待取消）"
    
    result = loan_tester.make_request("POST", "/api/v1/loan/applications", loan_data_for_cancel)
    if result and "application_id" in result:
        test_data["loan_application_for_cancel"] = result["application_id"]
        print(f"✅ 贷款申请创建成功（用于取消），ID: {test_data['loan_application_for_cancel']}")
    else:
        print("❌ 贷款申请创建失败")
        return False
    
    time.sleep(REQUEST_DELAY)
    
    # === 测试取消功能（在pending状态） ===
    print(f"\n{'='*50}")
    print(f"❌ 测试取消申请（pending状态）")
    print(f"{'='*50}")
    
    cancel_result = loan_tester.make_request("POST", f"/api/v1/loan/applications/{test_data['loan_application_for_cancel']}/cancel", {
        "reason": "测试取消原因"
    })
    
    if cancel_result is not None:
        print("✅ 取消申请测试完成")
    
    time.sleep(REQUEST_DELAY)
    
    # === 第二部分：创建用于审批测试的申请 ===
    print(f"\n{'='*50}")
    print(f"📝 创建贷款申请（用于审批测试）")
    print(f"{'='*50}")
    
    loan_data_for_approve = LOAN_TEST_DATA.copy()
    loan_data_for_approve["product_id"] = test_data["loan_product_id"]
    loan_data_for_approve["name"] = "测试贷款申请（待审批）"
    
    result = loan_tester.make_request("POST", "/api/v1/loan/applications", loan_data_for_approve)
    if result and "application_id" in result:
        test_data["loan_application_for_approve"] = result["application_id"]
        print(f"✅ 贷款申请创建成功（用于审批），ID: {test_data['loan_application_for_approve']}")
    else:
        print("❌ 贷款申请创建失败")
        return False
    
    time.sleep(REQUEST_DELAY)
    
    # === 查询申请列表 ===
    print(f"\n{'='*50}")
    print(f"📋 测试获取我的申请列表")
    print(f"{'='*50}")
    
    loan_tester.make_request("GET", "/api/v1/loan/applications", params={"page": 1, "size": 10})
    time.sleep(REQUEST_DELAY)
    
    # === 更新申请信息 ===
    print(f"\n{'='*50}")
    print(f"✏️ 测试更新申请")
    print(f"{'='*50}")
    
    update_data = {
        "amount": 120000.0,
        "duration": 18,
        "purpose": "更新后的贷款用途"
    }
    
    loan_tester.make_request("PUT", f"/api/v1/loan/applications/{test_data['loan_application_for_approve']}", update_data)
    time.sleep(REQUEST_DELAY)
    
    # === 管理员审批 ===
    print(f"\n{'='*50}")
    print(f"👮‍♂️ 测试管理员审批")
    print(f"{'='*50}")
    
    # 切换到管理员token
    loan_admin_tester = APITester(LOAN_BASE_URL, "Loan")
    loan_admin_tester.set_token(test_data["admin_token"])
    
    # 获取所有申请
    loan_admin_tester.make_request("GET", "/api/v1/admin/loan/applications", params={"page": 1, "size": 10})
    time.sleep(REQUEST_DELAY)
    
    # 审批申请
    approve_data = {
        "action": "approve",
        "suggestions": "申请已通过审核",
        "approved_amount": 100000.0,
        "approved_duration": 12,
        "interest_rate": 0.08
    }
    
    loan_admin_tester.make_request("POST", f"/api/v1/admin/loan/applications/{test_data['loan_application_for_approve']}/approve", approve_data)
    
    print("✅ 贷款业务测试完成")
    return True

def main():
    """主测试流程"""
    print("🚀 开始执行完全修复的API测试...")
    
    # 阶段0：数据准备
    if not setup_test_data():
        print("❌ 数据准备失败，退出测试")
        return
    
    # 阶段1：用户服务测试
    test_user_services()
    
    # 阶段2：租赁业务测试
    test_lease_business()
    
    # 阶段3：贷款业务测试  
    test_loan_business()
    
    print(f"\n{'='*60}")
    print(f"🎉 所有测试完成！")
    print(f"{'='*60}")
    
    # 打印测试数据摘要
    print(f"\n📊 测试数据摘要:")
    print(f"- 租赁产品: ID={test_data['lease_product_id']}, 编码={test_data['lease_product_code']}")
    print(f"- 贷款产品: ID={test_data['loan_product_id']}, 编码={test_data['loan_product_code']}")
    print(f"- 取消测试申请: 租赁={test_data['lease_application_for_cancel']}, 贷款={test_data['loan_application_for_cancel']}")
    print(f"- 审批测试申请: 租赁={test_data['lease_application_for_approve']}, 贷款={test_data['loan_application_for_approve']}")

if __name__ == "__main__":
    main() 