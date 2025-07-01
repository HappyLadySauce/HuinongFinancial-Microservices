package svc

import (
	"api/internal/breaker"
	"api/internal/config"
	"rpc/oauserclient"

	"github.com/zeromicro/go-zero/zrpc"
)

type ServiceContext struct {
	Config    config.Config
	OaUserRpc oauserclient.OaUser

	// 熔断器客户端
	OaUserRpcBreaker *breaker.RpcBreakerClient
}

func NewServiceContext(c config.Config) *ServiceContext {
	return &ServiceContext{
		Config:    c,
		OaUserRpc: oauserclient.NewOaUser(zrpc.MustNewClient(c.OaUserRpc)),

		// 初始化熔断器
		OaUserRpcBreaker: breaker.NewRpcBreakerClient("oauser-rpc"),
	}
}
