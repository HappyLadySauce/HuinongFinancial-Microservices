# 🧪 惠农金服微服务 API 测试指南

## 📋 概述

本文档提供了惠农金服微服务系统的完整API测试方案，包含认证服务(Auth)、C端用户服务(AppUser)、B端管理服务(OAUser)的全面测试。

## 🗂️ 测试文件清单

### 测试脚本
- `test_api_comprehensive.sh` - **完整测试脚本**（推荐）
- `test_api_quick.sh` - **快速测试脚本**（基础功能）
- `test_auth.sh` - **认证专项测试**

### 数据文件
- `init_test_data.sql` - **测试数据初始化脚本**
- `API_TEST_README.md` - **本说明文档**

## 🚀 快速开始

### 步骤 1: 初始化测试数据

```bash
# 初始化数据库和测试数据
mysql -u root -p < init_test_data.sql
```

### 步骤 2: 启动微服务

确保以下服务正常运行：
- **Auth API**: http://127.0.0.1:10003
- **AppUser API**: http://127.0.0.1:10001  
- **OAUser API**: http://127.0.0.1:10002

### 步骤 3: 运行测试脚本

```bash
# 方式一：完整测试（推荐）
chmod +x test_api_comprehensive.sh
./test_api_comprehensive.sh

# 方式二：快速测试
chmod +x test_api_quick.sh
./test_api_quick.sh

# 方式三：认证专项测试
chmod +x test_auth.sh
./test_auth.sh
```

## 🔐 测试账号信息

### C端用户账号
| 账号 | 密码 | 姓名 | 状态 | 用途 |
|------|------|------|------|------|
| 13800138000 | 123456 | 张三 | 正常 | 主要测试账号 |
| 13800138001 | 123456 | 李四 | 正常 | 业务流程测试 |
| 13800138002 | 123456 | 王五 | 正常 | 农机师测试 |
| 13800138003 | 123456 | 赵六 | 正常 | 畜牧员测试 |
| 13800138008 | 123456 | 郑十一 | 冻结 | 冻结状态测试 |
| 13800138009 | 123456 | 王十二 | 禁用 | 禁用状态测试 |

### B端管理员账号  
| 用户名 | 密码 | 姓名 | 角色 | 用途 |
|--------|------|------|------|------|
| admin | admin123 | 系统管理员 | 超级管理员 | 主要测试账号 |
| editor | editor123 | 内容编辑 | 编辑员 | 权限测试 |
| auditor | auditor123 | 内容审核 | 审核员 | 权限测试 |
| viewer | viewer123 | 只读用户 | 查看员 | 权限测试 |

## 📊 测试覆盖范围

### 🔐 认证服务 (Auth Service)
- ✅ C端用户登录/登出
- ✅ B端管理员登录/登出  
- ✅ Token刷新
- ✅ 错误凭据处理
- ✅ Token验证

### 👤 C端用户服务 (AppUser Service)
- ✅ 获取用户信息
- ✅ 更新用户档案
- ✅ JWT权限验证
- ✅ 参数验证

### 🛡️ B端管理服务 (OAUser Service)  
- ✅ 用户列表查询（分页/搜索）
- ✅ 创建用户
- ✅ 更新用户信息
- ✅ 删除用户
- ✅ 获取用户详情
- ✅ 权限级别验证

### 🧪 高级测试场景
- ✅ 数据边界测试
- ✅ 并发请求测试  
- ✅ SQL注入防护测试
- ✅ 性能基准测试
- ✅ 完整业务流程测试

## 📈 测试结果说明

### 成功标识
- ✅ **绿色** - 测试通过
- 🔄 **紫色** - 正在执行
- ℹ️ **蓝色** - 信息提示

### 错误标识  
- ❌ **红色** - 测试失败
- ⚠️ **黄色** - 警告信息

### HTTP状态码说明
- **200-299**: 成功
- **400-499**: 客户端错误（参数错误、权限不足等）
- **500-599**: 服务器错误

## 🔧 故障排查

### 常见问题

#### 1. 服务连接失败
```bash
# 检查服务是否启动
curl -I http://127.0.0.1:10003
curl -I http://127.0.0.1:10001  
curl -I http://127.0.0.1:10002
```

#### 2. 数据库连接错误
```bash
# 检查数据库连接
mysql -u root -p -e "SHOW DATABASES;"
```

#### 3. Redis连接错误
```bash
# 检查Redis连接
redis-cli ping
```

#### 4. Token获取失败
- 检查认证服务是否正常
- 确认测试账号密码正确
- 查看服务日志排查错误

### 日志查看
```bash
# 查看服务日志
tail -f app/auth/cmd/api/logs/*.log
tail -f app/auth/cmd/rpc/logs/*.log
tail -f app/appuser/cmd/api/logs/*.log
tail -f app/appuser/cmd/rpc/logs/*.log
```

## 🎯 自定义测试

### 修改测试配置
编辑测试脚本中的服务地址：
```bash
# 在脚本中修改这些变量
AUTH_API="http://127.0.0.1:10003/api/v1/auth"
APPUSER_API="http://127.0.0.1:10001/api/v1/appuser"  
OAUSER_API="http://127.0.0.1:10002/api/v1/oa"
```

### 添加新测试用例
参考现有测试函数编写新的测试：
```bash
test_custom_scenario() {
    print_header "🧪 自定义测试场景"
    
    # 你的测试逻辑
    make_request "POST" "$API_URL/custom" \
        '{"key":"value"}' \
        "-H \"Authorization: Bearer $TOKEN\"" \
        "自定义测试描述"
}
```

## 📋 测试报告

### 性能基准
- **登录QPS**: ~100+ (视硬件配置)
- **查询QPS**: ~200+ (视硬件配置)
- **平均响应时间**: <100ms (本地环境)

### 测试通过率目标
- **核心功能**: 100%
- **错误处理**: 100%  
- **性能测试**: 95%+
- **安全测试**: 100%

## 🚦 CI/CD 集成

### Jenkins集成示例
```groovy
pipeline {
    stages {
        stage('API Test') {
            steps {
                sh './test_api_comprehensive.sh'
            }
        }
    }
}
```

### GitHub Actions集成示例
```yaml
- name: Run API Tests
  run: |
    chmod +x test_api_comprehensive.sh
    ./test_api_comprehensive.sh
```

## 🤝 贡献指南

### 添加新测试
1. 在对应的测试函数中添加测试用例
2. 更新测试账号信息（如需要）
3. 更新本文档的测试覆盖范围

### 报告问题
- 描述测试环境
- 提供错误日志
- 说明复现步骤

## 📞 技术支持

如果您在使用测试脚本时遇到问题，请：

1. 检查[故障排查](#-故障排查)章节
2. 查看服务日志文件
3. 确认网络连接和服务状态
4. 联系开发团队获取支持

---

**🎉 祝您测试愉快！** 

*该测试套件将帮助确保惠农金服微服务系统的稳定性和可靠性。* 