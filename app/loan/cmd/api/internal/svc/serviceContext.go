package svc

import (
	"api/internal/config"
	"api/internal/middleware"
	"rpc/loanclient"

	"github.com/zeromicro/go-zero/rest"
	"github.com/zeromicro/go-zero/zrpc"
)

type ServiceContext struct {
	Config    config.Config
	AdminAuth rest.Middleware
	LoanRpc   loanclient.Loan
}

func NewServiceContext(c config.Config) *ServiceContext {
	return &ServiceContext{
		Config:    c,
		AdminAuth: middleware.NewAdminAuthMiddleware().Handle,
		LoanRpc:   loanclient.NewLoan(zrpc.MustNewClient(c.LoanRpc)),
	}
}
