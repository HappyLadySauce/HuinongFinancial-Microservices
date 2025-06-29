package svc

import (
	"model"
	"rpc/internal/config"

	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

type ServiceContext struct {
	Config                 config.Config
	LeaseApplicationsModel model.LeaseApplicationsModel
	LeaseApprovalsModel    model.LeaseApprovalsModel
}

func NewServiceContext(c config.Config) *ServiceContext {
	conn := sqlx.NewMysql(c.MySQL.DataSource)
	return &ServiceContext{
		Config:                 c,
		LeaseApplicationsModel: model.NewLeaseApplicationsModel(conn, c.CacheConf),
		LeaseApprovalsModel:    model.NewLeaseApprovalsModel(conn, c.CacheConf),
	}
}
