package svc

import (
	"rpc/internal/config"
	"model"

	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

type ServiceContext struct {
	Config config.Config
	AppUserModel model.AppUsersModel
}

func NewServiceContext(c config.Config) *ServiceContext {
	conn := sqlx.NewMysql(c.MySQL.DataSource)
	return &ServiceContext{
		Config: c,
		AppUserModel: model.NewAppUsersModel(conn, c.CacheConf),
	}
}
