package model

import (
	"github.com/zeromicro/go-zero/core/stores/cache"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

var _ LeaseApplicationsModel = (*customLeaseApplicationsModel)(nil)

type (
	// LeaseApplicationsModel is an interface to be customized, add more methods here,
	// and implement the added methods in customLeaseApplicationsModel.
	LeaseApplicationsModel interface {
		leaseApplicationsModel
	}

	customLeaseApplicationsModel struct {
		*defaultLeaseApplicationsModel
	}
)

// NewLeaseApplicationsModel returns a model for the database table.
func NewLeaseApplicationsModel(conn sqlx.SqlConn, c cache.CacheConf, opts ...cache.Option) LeaseApplicationsModel {
	return &customLeaseApplicationsModel{
		defaultLeaseApplicationsModel: newLeaseApplicationsModel(conn, c, opts...),
	}
}
