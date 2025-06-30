package svc

import (
	"model"
	"rpc/internal/clients"
	"rpc/internal/config"

	"github.com/zeromicro/go-zero/core/stores/sqlx"
	"github.com/zeromicro/go-zero/zrpc"
)

type ServiceContext struct {
	Config                config.Config
	LoanApplicationsModel model.LoanApplicationsModel
	LoanApprovalsModel    model.LoanApprovalsModel

	// RPC 客户端 - 通过consul服务发现调用其他服务
	LoanProductClient clients.LoanProductClient
	AppUserClient     clients.AppUserClient
}

func NewServiceContext(c config.Config) *ServiceContext {
	conn := sqlx.NewMysql(c.MySQL.DataSource)
	return &ServiceContext{
		Config:                c,
		LoanApplicationsModel: model.NewLoanApplicationsModel(conn, c.CacheConf),
		LoanApprovalsModel:    model.NewLoanApprovalsModel(conn, c.CacheConf),

		// 通过consul服务发现初始化RPC客户端
		LoanProductClient: clients.NewLoanProductClient(zrpc.MustNewClient(c.LoanProductRpc)),
		AppUserClient:     clients.NewAppUserClient(zrpc.MustNewClient(c.AppUserRpc)),
	}
}
