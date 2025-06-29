package svc

import (
	"api/internal/config"
	"rpc/authclient"

	"github.com/zeromicro/go-zero/zrpc"
)

type ServiceContext struct {
	Config  config.Config
	AuthRpc authclient.Auth
}

func NewServiceContext(c config.Config) *ServiceContext {
	return &ServiceContext{
		Config:  c,
		AuthRpc: authclient.NewAuth(zrpc.MustNewClient(c.AuthRpc)),
	}
}
