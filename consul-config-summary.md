# Consul配置和端口规划总结

## 环境配置
- **Consul地址**: `consul.huinong.internal` (测试环境)
- **生产环境地址**: `consul-master-cluster-consul-server.consul.svc.cluster.local:8500` (K8s部署时使用)

## 端口规划

### API服务端口 (10001-10999)
| 服务 | 端口 | 配置文件 | 状态 |
|------|------|----------|------|
| AppUser API | 10001 | `app/appuser/cmd/api/etc/appuser.yaml` | ✅ 已配置 |
| Auth API | 10002 | `app/auth/cmd/api/etc/auth.yaml` | ✅ 已配置 |
| OAUser API | 10003 | `app/oauser/cmd/api/etc/oauser.yaml` | ✅ 已配置 |
| Loan API | 10004 | `app/loan/cmd/api/etc/loan.yaml` | ✅ 已配置 |
| LoanProduct API | 10005 | `app/loanproduct/cmd/api/etc/loanproduct.yaml` | ✅ 已配置 |
| Lease API | 10006 | `app/lease/cmd/api/etc/leaseApi.yaml` | ✅ 已配置 |
| LeaseProduct API | 10007 | `app/leaseproduct/cmd/api/etc/leaseProductApi.yaml` | ✅ 已配置 |

### RPC服务端口 (20001-20999)
| 服务 | 端口 | 配置文件 | Consul Key | 状态 |
|------|------|----------|------------|------|
| AppUser RPC | 20001 | `app/appuser/cmd/rpc/etc/appuserrpc.yaml` | `appuserrpc.rpc` | ✅ 已配置 |
| Auth RPC | 20002 | `app/auth/cmd/rpc/etc/authrpc.yaml` | `authrpc.rpc` | ✅ 已配置 |
| OAUser RPC | 20003 | `app/oauser/cmd/rpc/etc/oauserrpc.yaml` | `oauserrpc.rpc` | ✅ 已配置 |
| Loan RPC | 20004 | `app/loan/cmd/rpc/etc/loanrpc.yaml` | `loanrpc.rpc` | ✅ 已配置 |
| LoanProduct RPC | 20005 | `app/loanproduct/cmd/rpc/etc/loanproductrpc.yaml` | `loanproductrpc.rpc` | ✅ 已配置 |
| Lease RPC | 20006 | `app/lease/cmd/rpc/etc/leaserpc.yaml` | `leaserpc.rpc` | ✅ 已配置 |
| LeaseProduct RPC | 20007 | `app/leaseproduct/cmd/rpc/etc/leaseproductrpc.yaml` | `leaseproductrpc.rpc` | ✅ 已配置 |

## RPC服务配置模式

### Config结构 (所有RPC服务统一)
```go
package config

import (
	"github.com/zeromicro/go-zero/zrpc"
	"github.com/zeromicro/zero-contrib/zrpc/registry/consul"
)

type Config struct {
	zrpc.RpcServerConf
	Consul consul.Conf
}
```

### YAML配置模式
```yaml
Name: servicename.rpc
ListenOn: 0.0.0.0:2000X

# Consul配置
Consul:
  Host: consul.huinong.internal
  Key: servicename.rpc
```

### 启动文件配置
```go
// 注册服务到consul
_ = consul.RegisterService(c.ListenOn, c.Consul)
```

## API服务RPC客户端配置示例

### Config结构 (API服务)
```go
type Config struct {
	rest.RestConf
	Auth struct {
		AccessSecret string
		AccessExpire int64
	}
	
	// RPC客户端配置
	AppUserRpc zrpc.RpcClientConf
	OAUserRpc  zrpc.RpcClientConf
}
```

### YAML配置
```yaml
# RPC客户端配置 (通过consul发现服务)
AppUserRpc:
  Target: consul://consul.huinong.internal/appuserrpc.rpc

OAUserRpc:
  Target: consul://consul.huinong.internal/oauserrpc.rpc
```

### 启动文件导入
```go
import (
    _ "github.com/zeromicro/zero-contrib/zrpc/registry/consul"
)
```

## JWT配置统一
所有API服务使用相同的JWT配置：
```yaml
Auth:
  AccessSecret: "huinong-auth-access-secret"
  AccessExpire: 3600
```

## 后续步骤
1. ✅ 端口规划完成
2. ✅ Consul配置完成
3. ✅ RPC服务注册配置完成
4. ✅ API服务RPC客户端配置(Auth API示例)
5. 🔄 其他API服务的RPC客户端配置
6. 🔄 安装consul依赖包 (`go get -u github.com/zeromicro/zero-contrib/zrpc/registry/consul`)
7. 🔄 测试服务启动和注册

## 注意事项
- RPC服务通过consul自动注册和发现
- API服务通过consul://地址调用RPC服务
- 测试环境使用 `consul.huinong.internal`
- 生产环境(K8s)使用 `consul-master-cluster-consul-server.consul.svc.cluster.local:8500` 