package model

import (
	"context"
	"fmt"

	"github.com/zeromicro/go-zero/core/stores/cache"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

var _ LeaseApplicationsModel = (*customLeaseApplicationsModel)(nil)

type (
	// LeaseApplicationsModel is an interface to be customized, add more methods here,
	// and implement the added methods in customLeaseApplicationsModel.
	LeaseApplicationsModel interface {
		leaseApplicationsModel
		// 自定义方法
		CountWithConditions(ctx context.Context, whereClause string, args []interface{}) (int64, error)
		ListWithConditions(ctx context.Context, whereClause string, args []interface{}, limit, offset int32) ([]*LeaseApplications, error)
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

// CountWithConditions 根据条件统计申请数量
func (m *customLeaseApplicationsModel) CountWithConditions(ctx context.Context, whereClause string, args []interface{}) (int64, error) {
	var total int64
	query := fmt.Sprintf("SELECT COUNT(*) FROM %s %s", m.table, whereClause)
	err := m.QueryRowNoCacheCtx(ctx, &total, query, args...)
	return total, err
}

// ListWithConditions 根据条件查询申请列表
func (m *customLeaseApplicationsModel) ListWithConditions(ctx context.Context, whereClause string, args []interface{}, limit, offset int32) ([]*LeaseApplications, error) {
	query := fmt.Sprintf("SELECT %s FROM %s %s ORDER BY created_at DESC LIMIT ? OFFSET ?", leaseApplicationsRows, m.table, whereClause)
	queryArgs := append(args, limit, offset)

	var applications []*LeaseApplications
	err := m.QueryRowsNoCacheCtx(ctx, &applications, query, queryArgs...)
	if err != nil {
		return nil, err
	}

	return applications, nil
}
