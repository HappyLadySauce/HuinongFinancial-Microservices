package svc

import (
	"api/internal/breaker"
	"api/internal/config"
	"rpc/appuserclient"

	"github.com/zeromicro/go-zero/zrpc"
)

type ServiceContext struct {
	Config     config.Config
	AppUserRpc appuserclient.AppUser

	// 熔断器客户端
	AppUserRpcBreaker *breaker.RpcBreakerClient
}

func NewServiceContext(c config.Config) *ServiceContext {
	return &ServiceContext{
		Config:     c,
		AppUserRpc: appuserclient.NewAppUser(zrpc.MustNewClient(c.AppUserRpc)),

		// 初始化熔断器
		AppUserRpcBreaker: breaker.NewRpcBreakerClient("appuser-rpc"),
	}
}
