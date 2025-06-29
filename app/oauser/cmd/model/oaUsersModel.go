package model

import (
	"github.com/zeromicro/go-zero/core/stores/cache"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

var _ OaUsersModel = (*customOaUsersModel)(nil)

type (
	// OaUsersModel is an interface to be customized, add more methods here,
	// and implement the added methods in customOaUsersModel.
	OaUsersModel interface {
		oaUsersModel
	}

	customOaUsersModel struct {
		*defaultOaUsersModel
	}
)

// NewOaUsersModel returns a model for the database table.
func NewOaUsersModel(conn sqlx.SqlConn, c cache.CacheConf, opts ...cache.Option) OaUsersModel {
	return &customOaUsersModel{
		defaultOaUsersModel: newOaUsersModel(conn, c, opts...),
	}
}
