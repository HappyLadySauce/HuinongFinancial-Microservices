package model

import (
	"context"
	"fmt"

	"github.com/zeromicro/go-zero/core/stores/cache"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

var _ LeaseProductsModel = (*customLeaseProductsModel)(nil)

type (
	// LeaseProductsModel is an interface to be customized, add more methods here,
	// and implement the added methods in customLeaseProductsModel.
	LeaseProductsModel interface {
		leaseProductsModel
		// 自定义方法
		CountWithConditions(ctx context.Context, whereClause string, args []interface{}) (int64, error)
		ListWithConditions(ctx context.Context, whereClause string, args []interface{}, limit, offset int32) ([]*LeaseProducts, error)
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

// CountWithConditions 根据条件统计产品数量
func (m *customLeaseProductsModel) CountWithConditions(ctx context.Context, whereClause string, args []interface{}) (int64, error) {
	var total int64
	query := fmt.Sprintf("SELECT COUNT(*) FROM %s %s", m.table, whereClause)
	err := m.QueryRowNoCacheCtx(ctx, &total, query, args...)
	return total, err
}

// ListWithConditions 根据条件查询产品列表
func (m *customLeaseProductsModel) ListWithConditions(ctx context.Context, whereClause string, args []interface{}, limit, offset int32) ([]*LeaseProducts, error) {
	query := fmt.Sprintf("SELECT %s FROM %s %s ORDER BY created_at DESC LIMIT ? OFFSET ?", leaseProductsRows, m.table, whereClause)
	queryArgs := append(args, limit, offset)

	var products []*LeaseProducts
	err := m.QueryRowsNoCacheCtx(ctx, &products, query, queryArgs...)
	if err != nil {
		return nil, err
	}

	return products, nil
}
