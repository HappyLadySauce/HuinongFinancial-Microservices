#!/usr/bin/env python3
# -*- coding: utf-8 -*-
"""
完整的API测试脚本
测试所有微服务模块的接口：
- appuser (C端用户服务)
- oauser (B端用户服务) 
- leaseproduct (租赁产品服务)
- lease (租赁业务服务)
- loanproduct (贷款产品服务)
- loan (贷款业务服务)
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

    def test_register(self, phone: str, password: str, role: str = None) -> bool:
        """测试注册接口"""
        print(f"\n{'='*50}")
        print(f"🔐 测试注册 - {phone}")
        print(f"{'='*50}")
        
        data = {
            "phone": phone,
            "password": password
        }
        
        # B端服务需要role字段
        if role:
            data["role"] = role
            
        result = self.make_request("POST", "/api/v1/auth/register", data)
        
        if "token" in result:
            self.set_token(result["token"])
            print(f"✅ 注册成功，已设置token")
            return True
        else:
            print(f"❌ 注册失败")
            return False
    
    def register_or_login(self, phone: str, password: str, role: str = None) -> bool:
        """尝试注册，如果用户已存在则登录"""
        print(f"\n{'='*50}")
        print(f"🔐 尝试注册或登录 - {phone}")
        print(f"{'='*50}")
        
        # 先尝试注册
        data = {
            "phone": phone,
            "password": password
        }
        
        # B端服务需要role字段
        if role:
            data["role"] = role
            
        result = self.make_request("POST", "/api/v1/auth/register", data)
        
        if "token" in result:
            self.set_token(result["token"])
            print(f"✅ 注册成功，已设置token")
            return True
        elif "用户已存在" in str(result) or "already exists" in str(result).lower():
            print(f"ℹ️ 用户已存在，尝试登录...")
            # 用户已存在，尝试登录
            return self.test_login(phone, password)
        else:
            print(f"❌ 注册失败: {result}")
            return False
    
    def test_login(self, phone: str, password: str) -> bool:
        """测试登录接口"""
        print(f"\n{'='*50}")
        print(f"🔐 测试登录 - {phone}")
        print(f"{'='*50}")
        
        data = {
            "phone": phone,
            "password": password
        }
        
        result = self.make_request("POST", "/api/v1/auth/login", data)
        
        if "token" in result:
            self.set_token(result["token"])
            print(f"✅ 登录成功，已设置token")
            return True
        else:
            print(f"❌ 登录失败")
            return False

class LeaseProductTester(APITester):
    """租赁产品服务测试类"""
    
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

class LeaseTester(APITester):
    """租赁业务服务测试类"""
    
    def test_create_application(self, application_data: Dict) -> Optional[str]:
        """测试创建租赁申请"""
        print(f"\n{'='*50}")
        print(f"📝 测试创建租赁申请")
        print(f"{'='*50}")
        
        result = self.make_request("POST", "/api/v1/lease/applications", application_data)
        
        if "application_id" in result:
            application_id = result["application_id"]
            print(f"✅ 申请创建成功，ID: {application_id}")
            return application_id
        else:
            print(f"❌ 申请创建失败")
            return None
    
    def test_list_my_applications(self, params: Dict = None):
        """测试获取我的申请列表"""
        print(f"\n{'='*50}")
        print(f"📋 测试获取我的申请列表")
        print(f"{'='*50}")
        
        default_params = {"page": 1, "size": 10}
        if params:
            default_params.update(params)
            
        result = self.make_request("GET", "/api/v1/lease/applications", params=default_params)
        return result
    
    def test_get_application_detail(self, application_id: str):
        """测试获取申请详情"""
        print(f"\n{'='*50}")
        print(f"🔍 测试获取申请详情 - {application_id}")
        print(f"{'='*50}")
        
        result = self.make_request("GET", f"/api/v1/lease/applications/{application_id}")
        return result
    
    def test_update_application(self, application_id: str, update_data: Dict):
        """测试更新申请"""
        print(f"\n{'='*50}")
        print(f"✏️ 测试更新申请 - {application_id}")
        print(f"{'='*50}")
        
        result = self.make_request("PUT", f"/api/v1/lease/applications/{application_id}", update_data)
        return result
    
    def test_cancel_application(self, application_id: str, reason: str):
        """测试取消申请"""
        print(f"\n{'='*50}")
        print(f"❌ 测试取消申请 - {application_id}")
        print(f"{'='*50}")
        
        data = {"reason": reason}
        result = self.make_request("POST", f"/api/v1/lease/applications/{application_id}/cancel", data)
        return result
    
    def test_admin_list_all_applications(self, params: Dict = None):
        """测试管理员获取所有申请"""
        print(f"\n{'='*50}")
        print(f"👮‍♂️ 测试管理员获取所有申请")
        print(f"{'='*50}")
        
        default_params = {"page": 1, "size": 10}
        if params:
            default_params.update(params)
            
        result = self.make_request("GET", "/api/v1/admin/lease/applications", params=default_params)
        return result
    
    def test_admin_approve_application(self, application_id: str, approval_data: Dict):
        """测试管理员审批申请"""
        print(f"\n{'='*50}")
        print(f"✅ 测试管理员审批申请 - {application_id}")
        print(f"{'='*50}")
        
        result = self.make_request("POST", f"/api/v1/admin/lease/applications/{application_id}/approve", approval_data)
        return result

class LoanProductTester(APITester):
    """贷款产品服务测试类"""
    
    def test_create_product(self, product_data: Dict) -> Optional[int]:
        """测试创建贷款产品"""
        print(f"\n{'='*50}")
        print(f"🏗️ 测试创建贷款产品")
        print(f"{'='*50}")
        
        result = self.make_request("POST", "/api/v1/admin/loanproduct/products", product_data)
        
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
            
        result = self.make_request("GET", "/api/v1/loanproduct/products", params=default_params)
        return result
    
    def test_get_product_detail(self, product_id: int):
        """测试获取产品详情"""
        print(f"\n{'='*50}")
        print(f"🔍 测试获取产品详情 - {product_id}")
        print(f"{'='*50}")
        
        result = self.make_request("GET", f"/api/v1/loanproduct/products/{product_id}")
        return result
    
    def test_update_product(self, product_id: int, update_data: Dict):
        """测试更新产品"""
        print(f"\n{'='*50}")
        print(f"✏️ 测试更新产品 - {product_id}")
        print(f"{'='*50}")
        
        result = self.make_request("PUT", f"/api/v1/admin/loanproduct/products/{product_id}", update_data)
        return result
    
    def test_update_product_status(self, product_id: int, status: int):
        """测试更新产品状态"""
        print(f"\n{'='*50}")
        print(f"🔄 测试更新产品状态 - {product_id} -> {status}")
        print(f"{'='*50}")
        
        data = {"status": status}
        result = self.make_request("PUT", f"/api/v1/admin/loanproduct/products/{product_id}/status", data)
        return result
    
    def test_delete_product(self, product_id: int):
        """测试删除产品"""
        print(f"\n{'='*50}")
        print(f"🗑️ 测试删除产品 - {product_id}")
        print(f"{'='*50}")
        
        result = self.make_request("DELETE", f"/api/v1/admin/loanproduct/products/{product_id}")
        return result

class LoanTester(APITester):
    """贷款业务服务测试类"""
    
    def test_create_application(self, application_data: Dict) -> Optional[str]:
        """测试创建贷款申请"""
        print(f"\n{'='*50}")
        print(f"📝 测试创建贷款申请")
        print(f"{'='*50}")
        
        result = self.make_request("POST", "/api/v1/loan/applications", application_data)
        
        if "application_id" in result:
            application_id = result["application_id"]
            print(f"✅ 申请创建成功，ID: {application_id}")
            return application_id
        else:
            print(f"❌ 申请创建失败")
            return None
    
    def test_list_my_applications(self, params: Dict = None):
        """测试获取我的申请列表"""
        print(f"\n{'='*50}")
        print(f"📋 测试获取我的申请列表")
        print(f"{'='*50}")
        
        default_params = {"page": 1, "size": 10}
        if params:
            default_params.update(params)
            
        result = self.make_request("GET", "/api/v1/loan/applications", params=default_params)
        return result
    
    def test_get_application_detail(self, application_id: str):
        """测试获取申请详情"""
        print(f"\n{'='*50}")
        print(f"🔍 测试获取申请详情 - {application_id}")
        print(f"{'='*50}")
        
        result = self.make_request("GET", f"/api/v1/loan/applications/{application_id}")
        return result
    
    def test_update_application(self, application_id: str, update_data: Dict):
        """测试更新申请"""
        print(f"\n{'='*50}")
        print(f"✏️ 测试更新申请 - {application_id}")
        print(f"{'='*50}")
        
        result = self.make_request("PUT", f"/api/v1/loan/applications/{application_id}", update_data)
        return result
    
    def test_cancel_application(self, application_id: str, reason: str):
        """测试取消申请"""
        print(f"\n{'='*50}")
        print(f"❌ 测试取消申请 - {application_id}")
        print(f"{'='*50}")
        
        data = {"reason": reason}
        result = self.make_request("POST", f"/api/v1/loan/applications/{application_id}/cancel", data)
        return result
    
    def test_admin_list_all_applications(self, params: Dict = None):
        """测试管理员获取所有申请"""
        print(f"\n{'='*50}")
        print(f"👮‍♂️ 测试管理员获取所有申请")
        print(f"{'='*50}")
        
        default_params = {"page": 1, "size": 10}
        if params:
            default_params.update(params)
            
        result = self.make_request("GET", "/api/v1/admin/loan/applications", params=default_params)
        return result
    
    def test_admin_approve_application(self, application_id: str, approval_data: Dict):
        """测试管理员审批申请"""
        print(f"\n{'='*50}")
        print(f"✅ 测试管理员审批申请 - {application_id}")
        print(f"{'='*50}")
        
        result = self.make_request("POST", f"/api/v1/admin/loan/applications/{application_id}/approve", approval_data)
        return result

def test_appuser_service():
    """测试C端用户服务"""
    print(f"\n{'#'*60}")
    print(f"🌟 开始测试 C端用户服务 (appuser)")
    print(f"{'#'*60}")
    
    tester = APITester(APPUSER_BASE_URL, "AppUser")
    
    phone = APPUSER_TEST_DATA["phone"]
    password = APPUSER_TEST_DATA["password"]
    
    # 1. 尝试注册或登录
    if tester.register_or_login(phone, password):
        time.sleep(REQUEST_DELAY)
        
        # 2. 测试获取用户信息
        user_info_result = tester.make_request("GET", "/api/v1/user/info", {"phone": phone})
        time.sleep(REQUEST_DELAY)
        
        # 3. 测试更新用户信息
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
            if "status" in updated_user_info:
                del updated_user_info["status"]
            if "role" in updated_user_info:
                del updated_user_info["role"]
            
            tester.make_request("PUT", "/api/v1/user/info", {"user_info": updated_user_info})
            time.sleep(REQUEST_DELAY)
        
        # 4. 测试登出
        tester.make_request("POST", "/api/v1/auth/logout")
        tester.clear_token()

def test_oauser_service():
    """测试B端用户服务"""
    print(f"\n{'#'*60}")
    print(f"🌟 开始测试 B端用户服务 (oauser)")
    print(f"{'#'*60}")
    
    tester = APITester(OAUSER_BASE_URL, "OAUser")
    
    for user in OAUSER_TEST_DATA:
        print(f"\n{'*'*40}")
        print(f"测试用户: {user['phone']} (角色: {user['role']})")
        print(f"{'*'*40}")
        
        # 1. 尝试注册或登录
        if tester.register_or_login(user["phone"], user["password"], user["role"]):
            time.sleep(REQUEST_DELAY)
            
            # 2. 测试获取用户信息
            user_info_result = tester.make_request("GET", "/api/v1/user/info", {"phone": user["phone"]})
            time.sleep(REQUEST_DELAY)
            
            # 3. 测试更新用户状态
            tester.make_request("PUT", "/api/v1/user/status", {"phone": user["phone"], "status": 1})
            time.sleep(REQUEST_DELAY)
            
            # 4. 测试登出
            tester.make_request("POST", "/api/v1/auth/logout")
            tester.clear_token()

def test_leaseproduct_service():
    """测试租赁产品服务"""
    print(f"\n{'#'*60}")
    print(f"🌟 开始测试 租赁产品服务 (leaseproduct)")
    print(f"{'#'*60}")
    
    # 首先登录管理员账户
    user_tester = APITester(OAUSER_BASE_URL, "OAUser")
    admin_user = OAUSER_TEST_DATA[0]  # 使用第一个管理员账户
    
    if not user_tester.register_or_login(admin_user["phone"], admin_user["password"], admin_user["role"]):
        print("❌ 管理员登录失败，无法进行产品管理测试")
        return
    
    # 设置产品服务的token
    product_tester = LeaseProductTester(LEASEPRODUCT_BASE_URL, "LeaseProduct")
    product_tester.set_token(user_tester.token)
    
    try:
        # 1. 测试创建产品
        product_id = product_tester.test_create_product(LEASEPRODUCT_TEST_DATA)
        if not product_id:
            print("❌ 产品创建失败，跳过后续测试")
            return
        time.sleep(REQUEST_DELAY)
        
        # 2. 测试获取产品列表
        product_tester.test_list_products()
        time.sleep(REQUEST_DELAY)
        
        # 3. 测试获取产品详情
        product_tester.test_get_product_detail(LEASEPRODUCT_TEST_DATA["product_code"])
        time.sleep(REQUEST_DELAY)
        
        # 4. 测试更新产品
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
        product_tester.test_update_product(LEASEPRODUCT_TEST_DATA["product_code"], update_data)
        time.sleep(REQUEST_DELAY)
        
        # 5. 测试检查库存
        inventory_check = {
            "product_code": LEASEPRODUCT_TEST_DATA["product_code"],
            "quantity": 2,
            "start_date": "2024-01-01",
            "end_date": "2024-01-10"
        }
        product_tester.test_check_inventory(inventory_check)
        time.sleep(REQUEST_DELAY)
        
        # 6. 测试删除产品
        product_tester.test_delete_product(LEASEPRODUCT_TEST_DATA["product_code"])
        
    except Exception as e:
        print(f"❌ 租赁产品测试过程中出现错误: {str(e)}")

def test_lease_service():
    """测试租赁业务服务"""
    print(f"\n{'#'*60}")
    print(f"🌟 开始测试 租赁业务服务 (lease)")
    print(f"{'#'*60}")
    
    # 首先登录用户账户
    user_tester = APITester(APPUSER_BASE_URL, "AppUser")
    user_data = APPUSER_TEST_DATA
    
    if not user_tester.register_or_login(user_data["phone"], user_data["password"]):
        print("❌ 用户登录失败，无法进行租赁申请测试")
        return
    
    # 设置租赁服务的token
    lease_tester = LeaseTester(LEASE_BASE_URL, "Lease")
    lease_tester.set_token(user_tester.token)
    
    try:
        # 1. 测试创建租赁申请
        application_id = lease_tester.test_create_application(LEASE_TEST_DATA)
        if not application_id:
            print("❌ 申请创建失败，跳过后续测试")
            return
        time.sleep(REQUEST_DELAY)
        
        # 2. 测试获取我的申请列表
        lease_tester.test_list_my_applications()
        time.sleep(REQUEST_DELAY)
        
        # 3. 测试获取申请详情
        lease_tester.test_get_application_detail(application_id)
        time.sleep(REQUEST_DELAY)
        
        # 4. 测试更新申请
        update_data = {
            "purpose": "更新后的使用目的",
            "delivery_address": "更新后的地址",
            "contact_phone": "13900139000"
        }
        lease_tester.test_update_application(application_id, update_data)
        time.sleep(REQUEST_DELAY)
        
        # 5. 登录管理员进行审批测试
        admin_tester = APITester(OAUSER_BASE_URL, "OAUser")
        admin_user = OAUSER_TEST_DATA[0]
        
        if admin_tester.register_or_login(admin_user["phone"], admin_user["password"], admin_user["role"]):
            lease_tester.set_token(admin_tester.token)
            
            # 测试管理员获取所有申请
            lease_tester.test_admin_list_all_applications()
            time.sleep(REQUEST_DELAY)
            
            # 测试审批申请
            approval_data = {
                "action": "approve",
                "suggestions": "申请已通过审核",
                "approved_duration": 10,
                "approved_amount": 8000.00,
                "approved_deposit": 10000.00
            }
            lease_tester.test_admin_approve_application(application_id, approval_data)
            time.sleep(REQUEST_DELAY)
        
        # 切回用户身份，测试取消申请
        lease_tester.set_token(user_tester.token)
        lease_tester.test_cancel_application(application_id, "测试取消原因")
        
    except Exception as e:
        print(f"❌ 租赁业务测试过程中出现错误: {str(e)}")

def test_loanproduct_service():
    """测试贷款产品服务"""
    print(f"\n{'#'*60}")
    print(f"🌟 开始测试 贷款产品服务 (loanproduct)")
    print(f"{'#'*60}")
    
    # 首先登录管理员账户
    user_tester = APITester(OAUSER_BASE_URL, "OAUser")
    admin_user = OAUSER_TEST_DATA[0]
    
    if not user_tester.register_or_login(admin_user["phone"], admin_user["password"], admin_user["role"]):
        print("❌ 管理员登录失败，无法进行产品管理测试")
        return
    
    # 设置产品服务的token
    product_tester = LoanProductTester(LOANPRODUCT_BASE_URL, "LoanProduct")
    product_tester.set_token(user_tester.token)
    
    try:
        # 1. 测试创建产品
        product_id = product_tester.test_create_product(LOANPRODUCT_TEST_DATA)
        if not product_id:
            print("❌ 产品创建失败，跳过后续测试")
            return
        time.sleep(REQUEST_DELAY)
        
        # 2. 测试获取产品列表
        product_tester.test_list_products()
        time.sleep(REQUEST_DELAY)
        
        # 3. 测试获取产品详情
        product_tester.test_get_product_detail(product_id)
        time.sleep(REQUEST_DELAY)
        
        # 4. 测试更新产品
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
        product_tester.test_update_product(product_id, update_data)
        time.sleep(REQUEST_DELAY)
        
        # 5. 测试更新产品状态
        product_tester.test_update_product_status(product_id, 2)  # 禁用
        time.sleep(REQUEST_DELAY)
        product_tester.test_update_product_status(product_id, 1)  # 启用
        time.sleep(REQUEST_DELAY)
        
        # 6. 测试删除产品
        product_tester.test_delete_product(product_id)
        
    except Exception as e:
        print(f"❌ 贷款产品测试过程中出现错误: {str(e)}")

def test_loan_service():
    """测试贷款业务服务"""
    print(f"\n{'#'*60}")
    print(f"🌟 开始测试 贷款业务服务 (loan)")
    print(f"{'#'*60}")
    
    # 首先登录用户账户
    user_tester = APITester(APPUSER_BASE_URL, "AppUser")
    user_data = APPUSER_TEST_DATA
    
    if not user_tester.register_or_login(user_data["phone"], user_data["password"]):
        print("❌ 用户登录失败，无法进行贷款申请测试")
        return
    
    # 设置贷款服务的token
    loan_tester = LoanTester(LOAN_BASE_URL, "Loan")
    loan_tester.set_token(user_tester.token)
    
    try:
        # 1. 测试创建贷款申请
        application_id = loan_tester.test_create_application(LOAN_TEST_DATA)
        if not application_id:
            print("❌ 申请创建失败，跳过后续测试")
            return
        time.sleep(REQUEST_DELAY)
        
        # 2. 测试获取我的申请列表
        loan_tester.test_list_my_applications()
        time.sleep(REQUEST_DELAY)
        
        # 3. 测试获取申请详情
        loan_tester.test_get_application_detail(application_id)
        time.sleep(REQUEST_DELAY)
        
        # 4. 测试更新申请
        update_data = {
            "amount": 120000.00,
            "duration": 18,
            "purpose": "更新后的贷款用途"
        }
        loan_tester.test_update_application(application_id, update_data)
        time.sleep(REQUEST_DELAY)
        
        # 5. 登录管理员进行审批测试
        admin_tester = APITester(OAUSER_BASE_URL, "OAUser")
        admin_user = OAUSER_TEST_DATA[0]
        
        if admin_tester.register_or_login(admin_user["phone"], admin_user["password"], admin_user["role"]):
            loan_tester.set_token(admin_tester.token)
            
            # 测试管理员获取所有申请
            loan_tester.test_admin_list_all_applications()
            time.sleep(REQUEST_DELAY)
            
            # 测试审批申请
            approval_data = {
                "action": "approve",
                "suggestions": "申请已通过审核",
                "approved_amount": 100000.00,
                "approved_duration": 12,
                "interest_rate": 0.08
            }
            loan_tester.test_admin_approve_application(application_id, approval_data)
            time.sleep(REQUEST_DELAY)
        
        # 切回用户身份，测试取消申请
        loan_tester.set_token(user_tester.token)
        loan_tester.test_cancel_application(application_id, "测试取消原因")
        
    except Exception as e:
        print(f"❌ 贷款业务测试过程中出现错误: {str(e)}")

def main():
    """主函数"""
    print("🚀 开始完整的API测试")
    print("注意：请确保相关服务正在运行")
    print(f"- appuser服务: {APPUSER_BASE_URL}")
    print(f"- oauser服务: {OAUSER_BASE_URL}")
    print(f"- leaseproduct服务: {LEASEPRODUCT_BASE_URL}")
    print(f"- lease服务: {LEASE_BASE_URL}")
    print(f"- loanproduct服务: {LOANPRODUCT_BASE_URL}")
    print(f"- loan服务: {LOAN_BASE_URL}")
    
    try:
        # 测试用户服务
        test_appuser_service()
        time.sleep(2)
        test_oauser_service()
        time.sleep(2)
        
        # 测试产品服务
        test_leaseproduct_service()
        time.sleep(2)
        test_loanproduct_service()
        time.sleep(2)
        
        # 测试业务服务
        test_lease_service()
        time.sleep(2)
        test_loan_service()
        
        print(f"\n{'='*60}")
        print("🎉 所有测试完成！")
        print(f"{'='*60}")
        
    except KeyboardInterrupt:
        print("\n⚠️ 测试被用户中断")
    except Exception as e:
        print(f"\n❌ 测试过程中出现错误: {str(e)}")

if __name__ == "__main__":
    main() 