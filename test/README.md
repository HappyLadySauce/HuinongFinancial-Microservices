# API测试脚本

这是一个用于测试 `appuser`(C端用户服务) 和 `oauser`(B端用户服务) API接口的自动化测试脚本。

## 文件说明

- `test_api.py` - 主要的测试脚本
- `config.py` - 配置文件，包含服务地址和测试数据
- `requirements.txt` - Python依赖包
- `README.md` - 使用说明文档

## 安装依赖

```bash
pip install -r requirements.txt
```

## 配置说明

在 `config.py` 文件中可以修改以下配置：

### 服务地址配置
```python
APPUSER_BASE_URL = "http://localhost:8080"  # C端用户服务地址
OAUSER_BASE_URL = "http://localhost:8081"   # B端用户服务地址
```

### 测试数据配置

**C端用户测试数据:**
```python
APPUSER_TEST_DATA = {
    "phone": "13452552349",
    "password": "13452552349"
}
```

**B端用户测试数据:**
```python
OAUSER_TEST_DATA = [
    {
        "phone": "13452552349",
        "password": "13452552349", 
        "role": "admin"
    },
    {
        "phone": "13452552348", 
        "password": "13452552348",
        "role": "manager"
    }
]
```

## 运行测试

```bash
python test_api.py
```

## 测试流程

### C端用户服务 (appuser) 测试流程

1. **注册** - 使用配置的手机号和密码注册新用户
2. **获取用户信息** - 获取刚注册用户的详细信息
3. **更新用户信息** - 更新用户的个人信息（姓名、昵称、年龄、职业、地址、收入等）
4. **修改密码** - 修改用户密码
5. **登出** - 用户登出，清除token
6. **登录** - 使用新密码重新登录
7. **删除用户** - 删除测试用户

### B端用户服务 (oauser) 测试流程

对每个配置的管理用户执行以下测试：

1. **注册** - 使用配置的手机号、密码和角色注册新管理用户
2. **获取用户信息** - 获取刚注册用户的详细信息
3. **更新用户信息** - 更新用户的个人信息（姓名、昵称、年龄、角色等）
4. **修改密码** - 修改用户密码
5. **登出** - 用户登出，清除token
6. **登录** - 使用新密码重新登录

最后统一删除所有测试用户。

## API接口测试覆盖

### 通用接口 (两个服务都包含)

- `POST /api/v1/auth/register` - 用户注册
- `POST /api/v1/auth/login` - 用户登录
- `POST /api/v1/auth/logout` - 用户登出
- `POST /api/v1/auth/password` - 修改密码
- `GET /api/v1/user/info` - 获取用户信息
- `PUT /api/v1/user/info` - 更新用户信息
- `POST /api/v1/user/delete` - 删除用户

### 主要区别

- **B端服务注册**：需要额外的 `role` 字段
- **C端用户信息**：包含更多字段如 `address`、`income`、`occupation`
- **B端用户信息**：包含 `role` 字段

## 输出格式

测试脚本会输出详细的日志信息，包括：

- 🔄 请求信息（方法、URL）
- 📤 请求数据（JSON格式）
- 📊 响应状态码
- 📥 响应数据（JSON格式）
- ✅ 成功信息
- ❌ 错误信息

## 注意事项

1. **确保服务运行**：在运行测试前，请确保对应的服务正在运行
2. **端口配置**：默认配置为 appuser(8080端口) 和 oauser(8081端口)
3. **测试数据**：测试完成后会清理所有测试用户数据
4. **网络连接**：确保测试环境能正常访问配置的服务地址
5. **JWT Token**：脚本会自动处理JWT token的获取和使用

## 自定义扩展

如需添加更多测试用例，可以：

1. 在 `APITester` 类中添加新的测试方法
2. 在 `test_appuser_service()` 或 `test_oauser_service()` 函数中调用新方法
3. 根据需要修改 `config.py` 中的配置参数

## 错误处理

脚本包含完善的错误处理机制：

- 网络请求超时处理
- JSON解析错误处理
- 用户中断处理 (Ctrl+C)
- 异常信息详细输出 