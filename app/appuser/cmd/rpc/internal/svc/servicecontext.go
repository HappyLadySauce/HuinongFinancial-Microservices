package svc

import (
	"model"
	"rpc/internal/config"
	"rpc/internal/pkg/logger"

	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

type ServiceContext struct {
	Config       config.Config
	AppUserModel model.AppUsersModel
}

func NewServiceContext(c config.Config) *ServiceContext {
	// 初始化 logrus 日志
	loggerConfig := logger.Config{
		ServiceName: c.Logger.ServiceName,
		Mode:        c.Logger.Mode,
		Path:        c.Logger.Path,
		Level:       c.Logger.Level,
		KeepDays:    c.Logger.KeepDays,
		Compress:    c.Logger.Compress,
		MaxSize:     c.Logger.MaxSize,
		MaxBackups:  c.Logger.MaxBackups,
	}
	logger.InitLogger(loggerConfig)

	// 初始化数据库连接
	conn := sqlx.NewMysql(c.MySQL.DataSource)

	return &ServiceContext{
		Config:       c,
		AppUserModel: model.NewAppUsersModel(conn, c.CacheConf),
	}
}
