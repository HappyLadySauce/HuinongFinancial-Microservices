#!/usr/bin/env python3
# -*- coding: utf-8 -*-
"""
API测试脚本
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
    
    def make_request(self, method: str, endpoint: str, data: Optional[Dict] = None) -> Dict[Any, Any]:
        """发送HTTP请求"""
        url = f"{self.base_url}{endpoint}"
        
        print(f"\n🔄 [{self.service_name}] {method.upper()} {endpoint}")
        if data:
            print(f"📤 Request: {json.dumps(data, indent=2, ensure_ascii=False)}")
        
        try:
            if method.lower() == 'get':
                response = self.session.get(url, json=data, timeout=REQUEST_TIMEOUT)
            elif method.lower() == 'post':
                response = self.session.post(url, json=data, timeout=REQUEST_TIMEOUT)
            elif method.lower() == 'put':
                response = self.session.put(url, json=data, timeout=REQUEST_TIMEOUT)
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
    
    def test_get_user_info(self, phone: str):
        """测试获取用户信息"""
        print(f"\n{'='*50}")
        print(f"👤 测试获取用户信息 - {phone}")
        print(f"{'='*50}")
        
        data = {"phone": phone}
        result = self.make_request("GET", "/api/v1/user/info", data)
        return result
    
    def test_update_user_info(self, user_info: Dict):
        """测试更新用户信息"""
        print(f"\n{'='*50}")
        print(f"✏️ 测试更新用户信息")
        print(f"{'='*50}")
        
        data = {"user_info": user_info}
        result = self.make_request("PUT", "/api/v1/user/info", data)
        return result
    
    def test_update_user_status(self, phone: str, status: int):
        """测试更新用户状态 (仅B端服务支持)"""
        print(f"\n{'='*50}")
        print(f"🔄 测试更新用户状态 - {phone} -> {status}")
        print(f"{'='*50}")
        
        data = {
            "phone": phone,
            "status": status
        }
        
        result = self.make_request("PUT", "/api/v1/user/status", data)
        return result
    
    def test_change_password(self, phone: str, old_password: str, new_password: str):
        """测试修改密码"""
        print(f"\n{'='*50}")
        print(f"🔑 测试修改密码 - {phone}")
        print(f"{'='*50}")
        
        data = {
            "phone": phone,
            "old_password": old_password,
            "new_password": new_password
        }
        
        result = self.make_request("POST", "/api/v1/auth/password", data)
        return result
    
    def test_logout(self):
        """测试登出"""
        print(f"\n{'='*50}")
        print(f"🚪 测试登出")
        print(f"{'='*50}")
        
        result = self.make_request("POST", "/api/v1/auth/logout")
        self.clear_token()
        return result
    
    def test_delete_user(self, phone: str):
        """测试删除用户"""
        print(f"\n{'='*50}")
        print(f"🗑️ 测试删除用户 - {phone}")
        print(f"{'='*50}")
        
        data = {"phone": phone}
        result = self.make_request("POST", "/api/v1/user/delete", data)
        return result

def test_appuser_service():
    """测试C端用户服务"""
    print(f"\n{'#'*60}")
    print(f"🌟 开始测试 C端用户服务 (appuser)")
    print(f"🔍 注意：C端服务不支持status字段和状态管理")
    print(f"{'#'*60}")
    
    tester = APITester(APPUSER_BASE_URL, "AppUser")
    
    phone = APPUSER_TEST_DATA["phone"]
    password = APPUSER_TEST_DATA["password"]
    
    # 1. 尝试注册或登录
    if tester.register_or_login(phone, password):
        time.sleep(REQUEST_DELAY)
        
        # 2. 测试获取用户信息
        user_info_result = tester.test_get_user_info(phone)
        time.sleep(REQUEST_DELAY)
        
        # 3. 测试更新用户信息（如果获取成功）
        if "user_info" in user_info_result:
            updated_user_info = user_info_result["user_info"].copy()
            # C端用户信息字段：id, phone, name, nickname, age, gender, occupation, address, income, created_at, updated_at
            updated_user_info.update({
                "name": "测试用户更新",
                "nickname": "测试昵称",
                "age": 25,
                "gender": 1,
                "occupation": "软件工程师",
                "address": "北京市朝阳区",
                "income": 15000.00
            })
            # 确保不包含status字段
            if "status" in updated_user_info:
                del updated_user_info["status"]
            if "role" in updated_user_info:
                del updated_user_info["role"]
            
            tester.test_update_user_info(updated_user_info)
            time.sleep(REQUEST_DELAY)
        
        # 4. 测试修改密码
        new_password = "new_password_123"
        tester.test_change_password(phone, password, new_password)
        time.sleep(REQUEST_DELAY)
        
        # 5. 测试登出
        tester.test_logout()
        time.sleep(REQUEST_DELAY)
        
        # 6. 使用新密码登录
        if tester.test_login(phone, new_password):
            time.sleep(REQUEST_DELAY)
            
            # 7. 测试删除用户（使用新密码登录后删除）
            tester.test_delete_user(phone)
        else:
            # 如果新密码登录失败，使用原密码登录后删除
            print("⚠️ 新密码登录失败，使用原密码登录后删除用户")
            if tester.test_login(phone, password):
                time.sleep(REQUEST_DELAY)
                tester.test_delete_user(phone)

def test_oauser_service():
    """测试B端用户服务"""
    print(f"\n{'#'*60}")
    print(f"🌟 开始测试 B端用户服务 (oauser)")
    print(f"🔍 注意：B端服务支持status字段和状态管理")
    print(f"{'#'*60}")
    
    tester = APITester(OAUSER_BASE_URL, "OAUser")
    tested_users = []  # 记录成功测试的用户
    
    for user in OAUSER_TEST_DATA:
        print(f"\n{'*'*40}")
        print(f"测试用户: {user['phone']} (角色: {user['role']})")
        print(f"{'*'*40}")
        
        # 1. 尝试注册或登录
        if tester.register_or_login(user["phone"], user["password"], user["role"]):
            tested_users.append(user)
            time.sleep(REQUEST_DELAY)
            
            # 2. 测试获取用户信息
            user_info_result = tester.test_get_user_info(user["phone"])
            time.sleep(REQUEST_DELAY)
            
            # 3. 测试更新用户信息（如果获取成功）
            if "user_info" in user_info_result:
                updated_user_info = user_info_result["user_info"].copy()
                # B端用户信息字段：id, phone, name, nickname, age, gender, role, status, created_at, updated_at
                updated_user_info.update({
                    "name": f"管理员-{user['phone'][-4:]}",
                    "nickname": f"管理员昵称-{user['role']}",
                    "age": 30,
                    "gender": 1,
                    "role": user["role"]
                })
                # 确保status字段存在（B端服务需要）
                if "status" not in updated_user_info:
                    updated_user_info["status"] = 1  # 默认正常状态
                
                tester.test_update_user_info(updated_user_info)
                time.sleep(REQUEST_DELAY)
            
            # 4. 测试更新用户状态 (B端服务特有功能)
            print(f"\n🔄 测试状态管理功能")
            # 测试禁用用户
            tester.test_update_user_status(user["phone"], 2)  # 2=禁用
            time.sleep(REQUEST_DELAY)
            
            # 测试启用用户
            tester.test_update_user_status(user["phone"], 1)  # 1=正常
            time.sleep(REQUEST_DELAY)
            
            # 5. 测试修改密码
            new_password = f"new_{user['password']}"
            tester.test_change_password(user["phone"], user["password"], new_password)
            time.sleep(REQUEST_DELAY)
            
            # 6. 测试登出
            tester.test_logout()
            time.sleep(REQUEST_DELAY)
            
            # 7. 使用新密码登录并删除用户
            if tester.test_login(user["phone"], new_password):
                time.sleep(REQUEST_DELAY)
                tester.test_delete_user(user["phone"])
            else:
                # 如果新密码登录失败，使用原密码登录后删除
                print("⚠️ 新密码登录失败，使用原密码登录后删除用户")
                if tester.test_login(user["phone"], user["password"]):
                    time.sleep(REQUEST_DELAY)
                    tester.test_delete_user(user["phone"])
        
        print(f"\n完成用户 {user['phone']} 的测试")
    
    print(f"\n✅ B端服务测试完成，共测试了 {len(tested_users)} 个用户")

def main():
    """主函数"""
    print("🚀 开始API测试")
    print("注意：请确保相关服务正在运行")
    print(f"- appuser服务: {APPUSER_BASE_URL}")
    print(f"- oauser服务: {OAUSER_BASE_URL}")
    print("\n🔍 服务差异说明：")
    print("- C端服务(appuser): 无status字段，支持occupation/address/income字段")
    print("- B端服务(oauser): 有status字段和状态管理，支持role字段")
    
    try:
        # 测试C端用户服务
        test_appuser_service()
        
        # 等待一段时间
        time.sleep(2)
        
        # 测试B端用户服务
        test_oauser_service()
        
        print(f"\n{'='*60}")
        print("🎉 所有测试完成！")
        print(f"{'='*60}")
        
    except KeyboardInterrupt:
        print("\n⚠️ 测试被用户中断")
    except Exception as e:
        print(f"\n❌ 测试过程中出现错误: {str(e)}")

if __name__ == "__main__":
    main() 