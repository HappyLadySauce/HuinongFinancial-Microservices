package svc

import (
	"appuserrpc/appuserclient"
	"loanproductrpc/loanproductservice"
	"model"
	"rpc/internal/config"

	"github.com/zeromicro/go-zero/core/stores/sqlx"
	"github.com/zeromicro/go-zero/zrpc"
)

type ServiceContext struct {
	Config                config.Config
	LoanApplicationsModel model.LoanApplicationsModel
	LoanApprovalsModel    model.LoanApprovalsModel

	// RPC 客户端 - 通过consul服务发现调用其他服务
	LoanProductClient loanproductservice.LoanProductService
	AppUserClient     appuserclient.AppUser
}

func NewServiceContext(c config.Config) *ServiceContext {
	conn := sqlx.NewMysql(c.MySQL.DataSource)
	return &ServiceContext{
		Config:                c,
		LoanApplicationsModel: model.NewLoanApplicationsModel(conn, c.CacheConf),
		LoanApprovalsModel:    model.NewLoanApprovalsModel(conn, c.CacheConf),

		// 通过consul服务发现初始化RPC客户端
		LoanProductClient: loanproductservice.NewLoanProductService(zrpc.MustNewClient(c.LoanProductRpc)),
		AppUserClient:     appuserclient.NewAppUser(zrpc.MustNewClient(c.AppUserRpc)),
	}
}
