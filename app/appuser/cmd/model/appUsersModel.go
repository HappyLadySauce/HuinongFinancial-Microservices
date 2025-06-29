package model

import (
	"github.com/zeromicro/go-zero/core/stores/cache"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

var _ AppUsersModel = (*customAppUsersModel)(nil)

type (
	// AppUsersModel is an interface to be customized, add more methods here,
	// and implement the added methods in customAppUsersModel.
	AppUsersModel interface {
		appUsersModel
	}

	customAppUsersModel struct {
		*defaultAppUsersModel
	}
)

// NewAppUsersModel returns a model for the database table.
func NewAppUsersModel(conn sqlx.SqlConn, c cache.CacheConf, opts ...cache.Option) AppUsersModel {
	return &customAppUsersModel{
		defaultAppUsersModel: newAppUsersModel(conn, c, opts...),
	}
}
