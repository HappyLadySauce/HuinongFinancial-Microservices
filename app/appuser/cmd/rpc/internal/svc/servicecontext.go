package svc

import (
	"model"
	"rpc/internal/config"

	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

type ServiceContext struct {
	Config       config.Config
	AppUserModel model.AppUsersModel
}

func NewServiceContext(c config.Config) *ServiceContext {
	// 初始化数据库连接
	conn := sqlx.NewMysql(c.MySQL.DataSource)

	return &ServiceContext{
		Config:       c,
		AppUserModel: model.NewAppUsersModel(conn, c.CacheConf),
	}
}
