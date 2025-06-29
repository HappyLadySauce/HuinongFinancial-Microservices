package svc

import (
	"api/internal/config"
	"api/internal/middleware"
	"rpc/oauserclient"

	"github.com/zeromicro/go-zero/rest"
	"github.com/zeromicro/go-zero/zrpc"
)

type ServiceContext struct {
	Config    config.Config
	AdminAuth rest.Middleware
	OaUserRpc oauserclient.OAUser
}

func NewServiceContext(c config.Config) *ServiceContext {
	return &ServiceContext{
		Config:    c,
		AdminAuth: middleware.NewAdminAuthMiddleware().Handle,
		OaUserRpc: oauserclient.NewOAUser(zrpc.MustNewClient(c.OaUserRpc)),
	}
}
