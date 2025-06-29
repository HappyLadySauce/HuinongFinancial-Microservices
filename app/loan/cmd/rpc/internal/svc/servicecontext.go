package svc

import (
	"model"
	"rpc/internal/config"

	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

type ServiceContext struct {
	Config                config.Config
	LoanApplicationsModel model.LoanApplicationsModel
	LoanApprovalsModel    model.LoanApprovalsModel
}

func NewServiceContext(c config.Config) *ServiceContext {
	conn := sqlx.NewMysql(c.MySQL.DataSource)
	return &ServiceContext{
		Config:                c,
		LoanApplicationsModel: model.NewLoanApplicationsModel(conn, c.CacheConf),
		LoanApprovalsModel:    model.NewLoanApprovalsModel(conn, c.CacheConf),
	}
}
