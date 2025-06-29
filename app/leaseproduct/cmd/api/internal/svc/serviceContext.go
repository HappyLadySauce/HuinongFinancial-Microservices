package svc

import (
	"api/internal/config"
	"api/internal/middleware"
	"rpc/leaseproductservice"

	"github.com/zeromicro/go-zero/rest"
	"github.com/zeromicro/go-zero/zrpc"
)

type ServiceContext struct {
	Config          config.Config
	AdminAuth       rest.Middleware
	LeaseProductRpc leaseproductservice.LeaseProductService
}

func NewServiceContext(c config.Config) *ServiceContext {
	return &ServiceContext{
		Config:          c,
		AdminAuth:       middleware.NewAdminAuthMiddleware().Handle,
		LeaseProductRpc: leaseproductservice.NewLeaseProductService(zrpc.MustNewClient(c.LeaseProductRpc)),
	}
}
