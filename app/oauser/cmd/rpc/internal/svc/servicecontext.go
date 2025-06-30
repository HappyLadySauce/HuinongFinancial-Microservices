package svc

import (
	"model"
	"rpc/internal/config"
	"rpc/internal/pkg/utils"

	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

type ServiceContext struct {
	Config      config.Config
	OaUserModel model.OaUsersModel
	JwtUtils    *utils.JWTUtils
}

func NewServiceContext(c config.Config) *ServiceContext {
	// 初始化数据库连接
	conn := sqlx.NewMysql(c.MySQL.DataSource)

	return &ServiceContext{
		Config:      c,
		OaUserModel: model.NewOaUsersModel(conn, c.CacheConf),
		JwtUtils:    utils.NewJWTUtils(c.JwtAuth.AccessSecret),
	}
}
