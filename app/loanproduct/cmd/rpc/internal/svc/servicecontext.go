package svc

import (
	"model"
	"rpc/internal/config"

	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

type ServiceContext struct {
	Config           config.Config
	LoanProductModel model.LoanProductsModel
}

func NewServiceContext(c config.Config) *ServiceContext {
	conn := sqlx.NewMysql(c.MySQL.DataSource)
	return &ServiceContext{
		Config:           c,
		LoanProductModel: model.NewLoanProductsModel(conn, c.CacheConf),
	}
}
