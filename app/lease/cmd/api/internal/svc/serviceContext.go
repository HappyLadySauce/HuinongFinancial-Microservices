package svc

import (
	"api/internal/breaker"
	"api/internal/config"
	"api/internal/middleware"
	"rpc/leaseclient"

	"github.com/zeromicro/go-zero/rest"
	"github.com/zeromicro/go-zero/zrpc"
)

type ServiceContext struct {
	Config    config.Config
	AdminAuth rest.Middleware
	LeaseRpc  leaseclient.Lease

	// 熔断器客户端
	LeaseRpcBreaker *breaker.RpcBreakerClient
}

func NewServiceContext(c config.Config) *ServiceContext {
	return &ServiceContext{
		Config:    c,
		AdminAuth: middleware.NewAdminAuthMiddleware().Handle,
		LeaseRpc:  leaseclient.NewLease(zrpc.MustNewClient(c.LeaseRpc)),

		// 初始化熔断器
		LeaseRpcBreaker: breaker.NewRpcBreakerClient("lease-rpc"),
	}
}
