package svc

import (
	"api/internal/breaker"
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

	// 熔断器客户端
	LoanProductRpcBreaker *breaker.RpcBreakerClient
}

func NewServiceContext(c config.Config) *ServiceContext {
	return &ServiceContext{
		Config:         c,
		AdminAuth:      middleware.NewAdminAuthMiddleware().Handle,
		LoanProductRpc: loanproductservice.NewLoanProductService(zrpc.MustNewClient(c.LoanProductRpc)),

		// 初始化熔断器
		LoanProductRpcBreaker: breaker.NewRpcBreakerClient("loanproduct-rpc"),
	}
}
