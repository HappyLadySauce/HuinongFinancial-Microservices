package svc

import (
	"api/internal/config"
	"rpc/oauserclient"

	"github.com/zeromicro/go-zero/zrpc"
)

type ServiceContext struct {
	Config    config.Config
	OaUserRpc oauserclient.OaUser
}

func NewServiceContext(c config.Config) *ServiceContext {
	return &ServiceContext{
		Config:    c,
		OaUserRpc: oauserclient.NewOaUser(zrpc.MustNewClient(c.OaUserRpc)),
	}
}
