#!/usr/bin/env python3
# -*- coding: utf-8 -*-
"""
服务健康检查脚本
"""

import requests
import json
from config import APPUSER_BASE_URL, OAUSER_BASE_URL, REQUEST_TIMEOUT

def check_service_health(service_name: str, base_url: str) -> bool:
    """检查服务健康状态"""
    print(f"\n🔍 检查 {service_name} 服务状态...")
    print(f"📍 服务地址: {base_url}")
    
    try:
        # 尝试访问一个简单的接口来检查服务状态
        test_url = f"{base_url}/api/v1/auth/login"
        
        # 发送一个简单的POST请求（预期会返回错误，但说明服务在运行）
        response = requests.post(
            test_url, 
            json={"phone": "", "password": ""}, 
            timeout=REQUEST_TIMEOUT
        )
        
        print(f"✅ {service_name} 服务正在运行 (状态码: {response.status_code})")
        return True
        
    except requests.exceptions.ConnectionError:
        print(f"❌ {service_name} 服务连接失败 - 服务可能未启动")
        return False
    except requests.exceptions.Timeout:
        print(f"⏰ {service_name} 服务响应超时")
        return False
    except Exception as e:
        print(f"⚠️  {service_name} 服务检查出现异常: {str(e)}")
        return False

def main():
    """主函数"""
    print("🚀 开始服务健康检查")
    print("=" * 50)
    
    services_status = {}
    
    # 检查 appuser 服务
    services_status['appuser'] = check_service_health("AppUser (C端用户服务)", APPUSER_BASE_URL)
    
    # 检查 oauser 服务
    services_status['oauser'] = check_service_health("OAUser (B端用户服务)", OAUSER_BASE_URL)
    
    # 总结
    print(f"\n{'='*50}")
    print("📊 健康检查结果:")
    print(f"{'='*50}")
    
    all_healthy = True
    for service, status in services_status.items():
        status_icon = "✅" if status else "❌"
        print(f"{status_icon} {service}: {'正常' if status else '异常'}")
        if not status:
            all_healthy = False
    
    if all_healthy:
        print(f"\n🎉 所有服务都在正常运行！可以开始API测试。")
        print(f"💡 运行测试: python test_api.py")
    else:
        print(f"\n⚠️  有服务未正常运行，请检查并启动服务后再进行测试。")
        print(f"💡 启动命令:")
        print(f"   ./scripts/start.sh start appuser-rpc")
        print(f"   ./scripts/start.sh start appuser-api")
        print(f"   ./scripts/start.sh start oauser-rpc")
        print(f"   ./scripts/start.sh start oauser-api")

if __name__ == "__main__":
    main() 