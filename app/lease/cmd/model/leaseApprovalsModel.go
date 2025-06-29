package model

import (
	"github.com/zeromicro/go-zero/core/stores/cache"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

var _ LeaseApprovalsModel = (*customLeaseApprovalsModel)(nil)

type (
	// LeaseApprovalsModel is an interface to be customized, add more methods here,
	// and implement the added methods in customLeaseApprovalsModel.
	LeaseApprovalsModel interface {
		leaseApprovalsModel
	}

	customLeaseApprovalsModel struct {
		*defaultLeaseApprovalsModel
	}
)

// NewLeaseApprovalsModel returns a model for the database table.
func NewLeaseApprovalsModel(conn sqlx.SqlConn, c cache.CacheConf, opts ...cache.Option) LeaseApprovalsModel {
	return &customLeaseApprovalsModel{
		defaultLeaseApprovalsModel: newLeaseApprovalsModel(conn, c, opts...),
	}
}
