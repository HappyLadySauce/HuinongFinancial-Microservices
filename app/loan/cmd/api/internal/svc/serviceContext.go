package svc

import (
	"api/internal/breaker"
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

	// 熔断器客户端
	LoanRpcBreaker *breaker.RpcBreakerClient
}

func NewServiceContext(c config.Config) *ServiceContext {
	return &ServiceContext{
		Config:    c,
		AdminAuth: middleware.NewAdminAuthMiddleware().Handle,
		LoanRpc:   loanclient.NewLoan(zrpc.MustNewClient(c.LoanRpc)),

		// 初始化熔断器
		LoanRpcBreaker: breaker.NewRpcBreakerClient("loan-rpc"),
	}
}
