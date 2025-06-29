# 惠农金融微服务平台 (HuinongFinancial-Microservices)

[![License](https://img.shields.io/badge/license-MIT-blue.svg)](LICENSE)
[![Go Version](https://img.shields.io/badge/go-1.21+-blue.svg)](https://golang.org/)
[![go-zero](https://img.shields.io/badge/go--zero-1.6+-brightgreen.svg)](https://go-zero.dev/)
[![Vue Version](https://img.shields.io/badge/vue-3.0+-green.svg)](https://vuejs.org/)
[![Docker](https://img.shields.io/badge/docker-ready-blue.svg)](https://www.docker.com/)
[![Kubernetes](https://img.shields.io/badge/kubernetes-ready-blue.svg)](https://kubernetes.io/)

## 📖 项目简介

惠农金融微服务平台是一个基于go-zero微服务框架的现代化金融服务系统，专为农业金融贷款服务场景设计。系统采用云原生微服务架构，提供高可用、高扩展性和高性能的金融贷款服务解决方案。

### 🎯 核心特性

- **go-zero微服务框架**: 高性能、生产级微服务框架
- **服务追踪**: 基于OpenTelemetry的分布式链路追踪
- **API自动生成**: 基于API定义自动生成代码
- **服务发现**: 内置服务注册与发现机制
- **熔断降级**: 自适应熔断和服务降级
- **负载均衡**: 多种负载均衡算法支持
- **缓存策略**: 多级缓存和缓存一致性
- **安全防护**: JWT认证，RBAC权限，API限流
- **监控告警**: 实时监控和智能告警

## 🏗️ 系统架构

### 整体架构图
```
┌─────────────────┐    ┌─────────────────┐    ┌─────────────────┐
│   用户端APP     │    │   管理端WEB     │    │   第三方支付    │
└─────────────────┘    └─────────────────┘    └─────────────────┘
         │                       │                       │
         └───────────────────────┼───────────────────────┘
                                 │
              ┌─────────────────────────────────┐
              │         go-zero Gateway         │  (API网关)
              │    (路由/鉴权/限流/监控)        │
              └─────────────────────────────────┘
                                 │
         ┌───────────────────────┼───────────────────────┐
         │                       │                       │
┌─────────────────┐    ┌─────────────────┐    ┌─────────────────┐
│   认证服务      │    │   用户服务      │    │  贷款产品服务   │
│  (Auth API)     │    │  (User API)     │    │(LoanProduct API)│
│                 │    │                 │    │                 │
│  认证RPC        │    │  用户RPC        │    │ 贷款产品RPC     │
└─────────────────┘    └─────────────────┘    └─────────────────┘
         │                       │                       │
         └───────────────────────┼───────────────────────┘
                                 │
              ┌─────────────────────────────────┐
              │       消息队列 & 服务追踪       │
              │   (NATS/Kafka + OpenTelemetry)  │
              └─────────────────────────────────┘
                                 │
         ┌───────────────────────┼───────────────────────┐
         │                       │   
┌─────────────────┐    ┌─────────────────┐
│  金融贷款服务   │    │   风控服务      │
│ (Loan API)      │    │ (Risk API)      │
│                 │    │                 │
│  贷款RPC        │    │  风控RPC        │
└─────────────────┘    └─────────────────┘
         │                       │
         └───────────────────────┘
                                 │
              ┌─────────────────────────────────┐
              │         数据存储层              │
              │  (MySQL/Redis/MongoDB/etcd)     │
              └─────────────────────────────────┘
```

### 微服务列表

| 服务名称 | API端口 | RPC端口 | 职责描述 | 技术栈 |
|---------|---------|---------|----------|--------|
| go-zero Gateway | 8080 | - | API网关、路由转发、认证鉴权、限流熔断 | go-zero gateway |
| Auth Service | 8081 | 9081 | 用户认证、Token管理、权限验证 | go-zero + MySQL + Redis |
| User Service | 8082 | 9082 | 用户信息管理、用户档案、实名认证 | go-zero + MySQL + Redis |
| LoanProduct Service | 8083 | 9083 | 贷款产品管理、利率配置、产品规则 | go-zero + MySQL + Redis |
| Loan Service | 8084 | 9084 | 贷款申请、审批流程、放款管理 | go-zero + MySQL + Redis |
| Risk Service | 8085 | 9085 | 风险评估、征信查询、反欺诈检测 | go-zero + MySQL + Redis |
| Admin Frontend | 3001 | - | 管理后台界面、运营管理 | Vue 3 + Element Plus |
| User Frontend | 3000 | - | 用户端界面、贷款申请 | Vue 3 + Vant |

## 🚀 快速开始

### 环境要求

- **开发环境**:
  - Go 1.21+
  - go-zero 1.6+
  - goctl (go-zero代码生成工具)
  - protoc (Protocol Buffer编译器)
  - air (实时热加载工具)
  - Node.js 16+
  - Docker 20.10+
  - Docker Compose 2.0+

- **生产环境**:
  - Kubernetes 1.20+
  - Helm 3.0+
  - etcd 3.5+ (服务发现)
  - OpenTelemetry Collector (链路追踪)
  - Jaeger/Zipkin (追踪可视化)

### 本地开发部署

#### 1. 克隆项目
```bash
git clone https://github.com/HappyLadySauce/HuinongFinancial-Microservices.git
cd HuinongFinancial-Microservices
```

#### 2. 安装开发工具
```bash
# 安装go-zero代码生成工具
go install github.com/zeromicro/go-zero/tools/goctl@latest

# 安装protoc工具 (macOS/Linux)
goctl env check --install --verbose --force

# 安装Air热加载工具
go install github.com/air-verse/air@latest

# 验证安装
goctl --version
protoc --version
air -v
```

#### 3. 启动基础设施
```bash
# 启动数据库、etcd、NATS等基础服务
docker-compose -f deploy/docker/infrastructure.yml up -d

# 等待服务启动完成
./scripts/wait-for-services.sh
```

#### 4. 生成代码
```bash
# 生成所有服务的API和RPC代码
make generate

# 单独生成某个服务
make generate-auth
make generate-user
make generate-loan-product
make generate-loan
```

#### 5. 启动微服务
```bash
# 方式一：Docker Compose 一键启动
docker-compose up -d

# 方式二：Air 热加载模式 (推荐开发使用)
make hot-reload-gateway
# 或分别启动
make hot-reload-auth
make hot-reload-user

# 方式三：手动启动
make start-auth-api
make start-auth-rpc
make start-user-api
make start-user-rpc
# ... 其他服务
```

#### 6. 初始化数据
```bash
# 创建数据库表结构
make migrate-up

# 导入初始数据
make seed-data
```

### 生产环境部署

#### 使用 Kubernetes

```bash
# 1. 创建命名空间
kubectl create namespace huinong-financial

# 2. 部署配置映射和密钥
kubectl apply -f k8s/configs/

# 3. 部署基础设施服务
kubectl apply -f k8s/infrastructure/

# 4. 部署微服务
kubectl apply -f k8s/services/

# 5. 部署前端应用
kubectl apply -f k8s/frontend/

# 6. 配置 Ingress
kubectl apply -f k8s/ingress/
```

#### 使用 Helm

```bash
# 添加 Helm 仓库
helm repo add huinong ./helm/charts

# 安装应用
helm install huinong-financial huinong/huinong-financial \
  --namespace huinong-financial \
  --create-namespace \
  --values helm/values/production.yaml
```

## 📁 项目结构

```
HuinongFinancial-Microservices/
├── app/                          # go-zero应用服务
│   ├── auth/                    # 认证服务
│   │   ├── api/                # API服务
│   │   └── rpc/                # RPC服务
│   ├── user/                   # 用户服务
│   │   ├── api/                # API服务
│   │   └── rpc/                # RPC服务
│   ├── loanproduct/            # 贷款产品服务
│   │   ├── api/                # API服务
│   │   └── rpc/                # RPC服务
│   ├── loan/                   # 金融贷款服务
│   │   ├── api/                # API服务
│   │   └── rpc/                # RPC服务
│   └── risk/                   # 风控服务
│       ├── api/                # API服务
│       └── rpc/                # RPC服务
├── common/                     # 公共组件
│   ├── config/                 # 配置文件
│   ├── middleware/             # 中间件
│   ├── utils/                  # 工具函数
│   ├── errorx/                 # 错误处理
│   ├── jwtx/                   # JWT工具
│   └── interceptor/            # RPC拦截器
├── frontend/                   # 前端应用
│   ├── admin/                  # 管理后台
│   └── user/                   # 用户端
├── deploy/                     # 部署配置
│   ├── docker/                 # Docker配置
│   │   ├── Dockerfile.*        # 各服务镜像
│   │   ├── docker-compose.yml  # 开发环境编排
│   │   └── infrastructure.yml  # 基础设施服务
│   ├── k8s/                    # Kubernetes配置
│   │   ├── auth/              # 认证服务部署
│   │   ├── user/              # 用户服务部署
│   │   ├── loanproduct/       # 贷款产品服务部署
│   │   ├── loan/              # 贷款服务部署
│   │   ├── gateway/           # 网关部署
│   │   └── infrastructure/    # 基础设施
│   └── helm/                   # Helm图表
│       ├── charts/            # 图表文件
│       └── values/            # 配置值
├── scripts/                    # 脚本文件
│   ├── generate.sh            # 代码生成脚本
│   ├── build.sh               # 构建脚本
│   ├── deploy.sh              # 部署脚本
│   └── migrate/               # 数据迁移
├── docs/                       # 文档
│   ├── api/                   # API文档
│   ├── architecture/          # 架构文档
│   ├── deployment/            # 部署文档
│   └── go-zero/               # go-zero使用文档
├── monitoring/                 # 监控配置
│   ├── prometheus/            # Prometheus配置
│   ├── grafana/              # Grafana面板
│   ├── jaeger/               # Jaeger链路追踪
│   └── otel/                 # OpenTelemetry配置
├── tests/                     # 测试文件
│   ├── unit/                 # 单元测试
│   ├── integration/          # 集成测试
│   └── e2e/                  # 端到端测试
├── tools/                     # 工具
│   ├── goctl/                # goctl模板
│   └── protoc/               # protoc插件
├── go.mod                     # Go模块定义
├── go.sum                     # Go依赖校验
├── Makefile                   # 构建文件
└── README.md                  # 项目说明
```

## 🔧 配置说明

### go-zero服务配置

每个微服务采用go-zero框架的配置结构：

#### API服务配置示例 (user-api.yaml)
```yaml
Name: user-api
Host: 0.0.0.0
Port: 8082
Mode: dev

# MySQL数据库配置
DataSource: root:password@tcp(localhost:3306)/huinong_financial?charset=utf8mb4&parseTime=true&loc=Asia%2FShanghai

# Redis配置
Redis:
  Host: localhost:6379
  Type: node
  Pass: ""

# JWT配置
Auth:
  AccessSecret: huinong-financial-jwt-secret
  AccessExpire: 86400  # 24小时

# 链路追踪配置
Telemetry:
  Name: user-api
  Endpoint: http://localhost:14268/api/traces
  Sampler: 1.0
  Batcher: jaeger

# RPC服务连接配置
UserRpc:
  Etcd:
    Hosts:
      - localhost:2379
    Key: user.rpc

# 日志配置
Log:
  ServiceName: user-api
  Mode: console
  Level: info
```

#### RPC服务配置示例 (user-rpc.yaml)
```yaml
Name: user.rpc
ListenOn: 0.0.0.0:9082
Mode: dev

# 服务注册配置
Etcd:
  Hosts:
    - localhost:2379
  Key: user.rpc

# MySQL数据库配置
DataSource: root:password@tcp(localhost:3306)/huinong_financial?charset=utf8mb4&parseTime=true&loc=Asia%2FShanghai

# Redis配置
Redis:
  Host: localhost:6379
  Type: node
  Pass: ""

# 链路追踪配置
Telemetry:
  Name: user-rpc
  Endpoint: http://localhost:14268/api/traces
  Sampler: 1.0
  Batcher: jaeger

# 缓存配置
Cache:
  - Host: localhost:6379
    Pass: ""
    DB: 1

# 日志配置
Log:
  ServiceName: user-rpc
  Mode: console
  Level: info
```

### 网关配置

go-zero网关配置示例：

```yaml
Name: gateway
Host: 0.0.0.0
Port: 8080
Mode: dev

# 上游服务配置
Upstreams:
  - Name: auth-api
    Uris:
      - http://localhost:8081
  - Name: user-api
    Uris:
      - http://localhost:8082
  - Name: loanproduct-api
    Uris:
      - http://localhost:8083
  - Name: loan-api
    Uris:
      - http://localhost:8084

# 路由配置
Mapping:
  - Method: post
    Path: /api/auth/**
    Upstream: auth-api
  - Method: "*"
    Path: /api/user/**
    Upstream: user-api
  - Method: "*"
    Path: /api/loanproduct/**
    Upstream: loanproduct-api
  - Method: "*"
    Path: /api/loan/**
    Upstream: loan-api

# 链路追踪配置
Telemetry:
  Name: gateway
  Endpoint: http://localhost:14268/api/traces
  Sampler: 1.0
  Batcher: jaeger

# 超时配置
Timeout: 30s

# 日志配置
Log:
  ServiceName: gateway
  Mode: console
  Level: info
```

## 📊 监控和日志

### 监控指标

系统集成了 Prometheus + Grafana 监控栈：

- **系统指标**: CPU、内存、网络、磁盘
- **业务指标**: QPS、响应时间、错误率
- **自定义指标**: 用户注册数、订单量、支付成功率

### 链路追踪

基于OpenTelemetry和Jaeger的分布式链路追踪：

#### 1. 启动Jaeger
```bash
# 使用Docker启动Jaeger
docker run -d --name jaeger \
  -e COLLECTOR_ZIPKIN_HOST_PORT=:9411 \
  -p 16686:16686 \
  -p 14268:14268 \
  -p 14250:14250 \
  -p 9411:9411 \
  jaegertracing/all-in-one:latest

# 访问 Jaeger UI
http://localhost:16686
```

#### 2. go-zero链路追踪配置
```yaml
# 在每个服务配置文件中添加
Telemetry:
  Name: service-name          # 服务名称
  Endpoint: http://localhost:14268/api/traces  # Jaeger端点
  Sampler: 1.0               # 采样率 (0.0-1.0)
  Batcher: jaeger            # 批处理器类型
```

#### 3. 自定义链路追踪
```go
// 在业务代码中添加自定义span
import (
    "go.opentelemetry.io/otel"
    "go.opentelemetry.io/otel/attribute"
    "go.opentelemetry.io/otel/trace"
)

func (l *UserLogic) CreateUser(req *types.CreateUserReq) (*types.CreateUserResp, error) {
    // 创建子span
    ctx, span := otel.Tracer("user-service").Start(l.ctx, "CreateUser")
    defer span.End()
    
    // 添加属性
    span.SetAttributes(
        attribute.String("user.mobile", req.Mobile),
        attribute.String("user.name", req.Name),
    )
    
    // 业务逻辑
    user, err := l.svcCtx.UserModel.Insert(ctx, &model.User{
        Mobile: req.Mobile,
        Name:   req.Name,
    })
    
    if err != nil {
        span.RecordError(err)
        return nil, err
    }
    
    span.SetAttributes(attribute.Int64("user.id", user.Id))
    return &types.CreateUserResp{Id: user.Id}, nil
}
```

### 日志聚合

使用 ELK Stack 进行日志收集和分析：

```bash
# 启动 ELK Stack
docker-compose -f monitoring/elk/docker-compose.yml up -d

# 访问 Kibana
http://localhost:5601
```

### go-zero 日志处理

项目采用 `go-zero` 内置的 `logx` 和 `logc` 进行日志记录，具备高性能和上下文感知能力。

- **logx**: 核心日志库，支持多种输出方式和日志级别。
- **logc**: `logx` 的封装，自动将日志与请求上下文（`context`）关联，方便链路追踪。

推荐在业务逻辑中使用 `logc`，示例：
```go
logc.Info(ctx, "用户创建成功")
logc.Errorf(ctx, "创建用户失败: %+v", err)
```

#### 日志配置

日志配置在各服务的 `yaml` 文件中定义，可以控制日志级别、输出模式（控制台、文件）、轮转策略等。

```yaml
# 日志配置
Log:
  ServiceName: user-api
  Mode: console  # console, file, volume
  Encoding: plain # json, plain
  Path: logs/user-api
  Level: info # debug, info, warn, error, severe
  KeepDays: 7
  Rotation: daily # daily, size
```

## 🔐 安全防护

### 认证授权

- **JWT Token**: 用户身份认证
- **RBAC**: 基于角色的访问控制
- **OAuth 2.0**: 第三方登录支持

### API安全

- **限流控制**: 防止API滥用
- **参数验证**: 输入数据校验
- **SQL注入防护**: 使用参数化查询
- **XSS防护**: 输出内容转义

### 网络安全

- **HTTPS**: 传输层加密
- **防火墙**: 网络访问控制
- **VPC**: 私有网络隔离

## 🧪 测试

### 单元测试

```bash
# 运行所有单元测试
make test

# 运行指定服务测试
make test-user-service

# 生成测试覆盖率报告
make test-coverage
```

### 集成测试

```bash
# 启动测试环境
make test-env-up

# 运行集成测试
make test-integration

# 清理测试环境
make test-env-down
```

### 性能测试

```bash
# 安装性能测试工具
go install github.com/rakyll/hey@latest

# 运行性能测试
make performance-test

# 生成性能报告
make performance-report
```

## 📚 API文档

### API定义文件

go-zero使用.api文件定义RESTful API：

```api
// app/auth/api/auth.api
syntax = "v1"

info(
    title: "认证服务API"
    desc: "用户认证、登录、注册相关接口"
    author: "huinong-team"
    email: "dev@huinong.com"
    version: "v1.0"
)

type (
    LoginReq {
        Mobile   string `json:"mobile" validate:"required,len=11"`   // 手机号
        Password string `json:"password" validate:"required,min=6"`  // 密码
        CaptchaId string `json:"captcha_id"`                        // 验证码ID
        Captcha   string `json:"captcha"`                           // 验证码
    }
    
    LoginResp {
        AccessToken  string `json:"access_token"`   // 访问令牌
        AccessExpire int64  `json:"access_expire"`  // 过期时间
        RefreshAfter int64  `json:"refresh_after"`  // 刷新时间
    }
    
    UserInfo {
        Id     int64  `json:"id"`
        Mobile string `json:"mobile"`
        Name   string `json:"name"`
        Avatar string `json:"avatar"`
    }
)

@server(
    prefix: /api/auth
    group: auth
)
service auth-api {
    @doc "用户登录"
    @handler login
    post /login (LoginReq) returns (LoginResp)
    
    @doc "用户注册"
    @handler register
    post /register (RegisterReq) returns (RegisterResp)
}

@server(
    prefix: /api/auth
    group: auth
    jwt: Auth
)
service auth-api {
    @doc "获取用户信息"
    @handler userInfo
    get /userinfo returns (UserInfo)
    
    @doc "刷新令牌"
    @handler refresh
    post /refresh returns (LoginResp)
}
```

### 自动生成API文档

```bash
# 使用goctl生成API文档
goctl api doc -dir app/auth/api/

# 生成的文档位于
# app/auth/api/auth.md
```

### 在线API文档

启动服务后访问 API 文档：

- 认证服务: http://localhost:8081/api/auth/doc
- 用户服务: http://localhost:8082/api/user/doc  
- 贷款产品服务: http://localhost:8083/api/loanproduct/doc
- 贷款服务: http://localhost:8084/api/loan/doc

### Postman集合

导入 `docs/postman/` 目录下的集合文件进行API测试。

## 🚢 部署指南

### Docker部署

```bash
# 构建所有服务镜像
make docker-build

# 推送镜像到仓库
make docker-push

# 部署到开发环境
make deploy-dev

# 部署到生产环境
make deploy-prod
```

### Kubernetes部署

```bash
# 创建命名空间
kubectl create namespace huinong-financial

# 部署应用
kubectl apply -f k8s/

# 检查部署状态
kubectl get pods -n huinong-financial

# 查看服务日志
kubectl logs -f deployment/user-service -n huinong-financial
```

### CI/CD流水线

项目集成了 GitHub Actions CI/CD：

```yaml
# .github/workflows/ci-cd.yml
name: CI/CD Pipeline

on:
  push:
    branches: [ main, develop ]
  pull_request:
    branches: [ main ]

jobs:
  test:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v3
      - uses: actions/setup-go@v3
        with:
          go-version: 1.21
      - run: make test

  build:
    needs: test
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v3
      - run: make docker-build
      - run: make docker-push

  deploy:
    needs: build
    runs-on: ubuntu-latest
    if: github.ref == 'refs/heads/main'
    steps:
      - run: make deploy-prod
```

## 🛠️ 常用命令

### 开发命令

```bash
# 安装依赖
make deps

# 生成所有服务代码
make generate

# 生成单个服务
make generate-auth      # 生成认证服务
make generate-user      # 生成用户服务
make generate-loan      # 生成贷款服务
make generate-risk      # 生成风控服务

# 代码格式化
make fmt

# 代码检查
make lint

# 构建所有服务
make build

# 构建单个服务
make build-auth-api     # 构建认证API服务
make build-auth-rpc     # 构建认证RPC服务
make build-user-api     # 构建用户API服务
make build-user-rpc     # 构建用户RPC服务

# 启动服务（开发模式）
make dev-auth          # 启动认证服务
make dev-user          # 启动用户服务
make dev-loan          # 启动贷款服务
make dev-gateway       # 启动网关

# 清理构建文件
make clean

# 生成API文档
make api-docs
```

### go-zero特有命令

```bash
# 生成API代码
goctl api go -api app/user/api/user.api -dir app/user/api/ --style goZero

# 生成RPC代码
goctl rpc protoc app/user/rpc/user.proto --go_out=app/user/rpc --go-grpc_out=app/user/rpc --zrpc_out=app/user/rpc --style goZero

# 生成数据库模型
goctl model mysql ddl -src common/sql/user.sql -dir app/user/model --style goZero

# 生成Docker文件
goctl docker -go app/user/api/user.go

# 生成Kubernetes部署文件
goctl kube deploy -name user-api -namespace huinong-financial -image user-api:latest -o deploy/k8s/user/user-api.yaml
```

### 运维命令

```bash
# 查看服务状态
make status

# 重启服务
make restart

# 查看日志
make logs

# 数据库迁移
make migrate-up

# 回滚数据库
make migrate-down

# 备份数据
make backup

# 恢复数据
make restore
```

## 🔗 相关链接

- [项目官网](https://huinong-financial.com)
- [API文档](https://docs.huinong-financial.com)
- [开发指南](./docs/development.md)
- [部署指南](./docs/deployment.md)
- [架构设计](./docs/architecture.md)
- [问题反馈](https://github.com/HappyLadySauce/HuinongFinancial-Microservices/issues)

## 🤝 贡献指南

我们欢迎所有形式的贡献，包括但不限于：

- 🐛 报告问题
- 💡 提出新功能
- 📝 改进文档
- 🔧 提交代码

请参阅 [CONTRIBUTING.md](./CONTRIBUTING.md) 了解详细的贡献流程。

## 📄 许可证

本项目采用 [MIT License](./LICENSE) 开源协议。

## 👥 团队成员

- **项目负责人**: HappyLadySauce
- **架构师**: AI Assistant
- **go-zero专家**: 待招募
- **前端负责人**: 待招募
- **运维负责人**: 待招募
- **开发团队**: 欢迎加入

### 技术栈总结

- **微服务框架**: go-zero (高性能微服务框架)
- **API网关**: go-zero gateway
- **服务发现**: etcd
- **链路追踪**: OpenTelemetry + Jaeger
- **数据库**: MySQL 8.0+
- **缓存**: Redis 6.0+
- **消息队列**: NATS/Kafka
- **前端框架**: Vue 3 + TypeScript
- **UI组件**: Element Plus (管理端) + Vant (移动端)
- **容器化**: Docker + Kubernetes
- **监控**: Prometheus + Grafana
- **日志**: ELK Stack

## 📞 联系我们

- **邮箱**: support@huinong-financial.com
- **微信群**: 扫描二维码加入技术交流群
- **QQ群**: 123456789
- **GitHub**: https://github.com/HappyLadySauce/HuinongFinancial-Microservices

---

**最后更新**: 2024年12月  
**版本**: v1.0.0  
**维护状态**: 🚀 积极维护中

## ⚙️ 内置中间件

go-zero 框架提供了丰富的内置中间件，在项目网关（Gateway）中统一配置和启用，用于保障服务的稳定性和安全性。

```go
// gateway.go
server.Use(server.Recover)
server.Use(server.MaxBytes(1024 * 1024 * 10)) // 10MB
server.Use(middleware.NewCorsMiddleware().Handle)
server.Use(middleware.NewLogMiddleware().Handle)
server.Use(server.Prometheus)
// ... 其他自定义中间件
```

### 常用中间件

- **Recover**: 捕获 `panic`，防止服务崩溃，并记录错误日志。
- **MaxBytes**: 限制请求体的大小，防止恶意大请求。
- **CorsMiddleware**: 处理跨域请求。
- **LogMiddleware**: 记录每个请求的详细日志。
- **Prometheus**: 暴露监控指标，供 Prometheus 采集。
- **Breaker**: 熔断器，防止服务雪崩。
- **Shedding**: 服务过载保护，主动丢弃请求。
- **Trace**: 链路追踪中间件。
