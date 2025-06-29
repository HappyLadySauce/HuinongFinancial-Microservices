package svc

import (
	"api/internal/config"
	"api/internal/middleware"
	"rpc/loanproductservice"

	"github.com/zeromicro/go-zero/rest"
	"github.com/zeromicro/go-zero/zrpc"
)

type ServiceContext struct {
	Config         config.Config
	AdminAuth      rest.Middleware
	LoanProductRpc loanproductservice.LoanProductService
}

func NewServiceContext(c config.Config) *ServiceContext {
	return &ServiceContext{
		Config:         c,
		AdminAuth:      middleware.NewAdminAuthMiddleware().Handle,
		LoanProductRpc: loanproductservice.NewLoanProductService(zrpc.MustNewClient(c.LoanProductRpc)),
	}
}
