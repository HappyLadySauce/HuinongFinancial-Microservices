package svc

import (
	"api/internal/config"
	"rpc/appuserclient"

	"github.com/zeromicro/go-zero/zrpc"
)

type ServiceContext struct {
	Config     config.Config
	AppUserRpc appuserclient.AppUser
}

func NewServiceContext(c config.Config) *ServiceContext {
	return &ServiceContext{
		Config:     c,
		AppUserRpc: appuserclient.NewAppUser(zrpc.MustNewClient(c.AppUserRpc)),
	}
}
