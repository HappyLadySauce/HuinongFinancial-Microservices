package model

import (
	"github.com/zeromicro/go-zero/core/stores/cache"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

var _ LeaseProductsModel = (*customLeaseProductsModel)(nil)

type (
	// LeaseProductsModel is an interface to be customized, add more methods here,
	// and implement the added methods in customLeaseProductsModel.
	LeaseProductsModel interface {
		leaseProductsModel
	}

	customLeaseProductsModel struct {
		*defaultLeaseProductsModel
	}
)

// NewLeaseProductsModel returns a model for the database table.
func NewLeaseProductsModel(conn sqlx.SqlConn, c cache.CacheConf, opts ...cache.Option) LeaseProductsModel {
	return &customLeaseProductsModel{
		defaultLeaseProductsModel: newLeaseProductsModel(conn, c, opts...),
	}
}
