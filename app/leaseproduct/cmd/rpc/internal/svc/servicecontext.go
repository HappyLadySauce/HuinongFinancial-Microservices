package svc

import (
	"model"
	"rpc/internal/config"

	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

type ServiceContext struct {
	Config            config.Config
	LeaseProductModel model.LeaseProductsModel
}

func NewServiceContext(c config.Config) *ServiceContext {
	conn := sqlx.NewMysql(c.MySQL.DataSource)
	return &ServiceContext{
		Config:            c,
		LeaseProductModel: model.NewLeaseProductsModel(conn, c.CacheConf),
	}
}
