package svc

import (
	"model"
	"rpc/internal/config"
	"rpc/internal/utils"

	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

type ServiceContext struct {
	Config        config.Config
	OaUsersModel  model.OaUsersModel
	PasswordUtil  *utils.PasswordUtil
}

func NewServiceContext(c config.Config) *ServiceContext {
	conn := sqlx.NewMysql(c.MySQL.DataSource)
	return &ServiceContext{
		Config:       c,
		OaUsersModel: model.NewOaUsersModel(conn, c.CacheConf),
		PasswordUtil: utils.NewPasswordUtil(),
	}
}
