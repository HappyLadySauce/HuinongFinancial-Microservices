package model

import (
	"context"
	"fmt"

	"github.com/zeromicro/go-zero/core/stores/cache"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

var _ LeaseApprovalsModel = (*customLeaseApprovalsModel)(nil)

type (
	// LeaseApprovalsModel is an interface to be customized, add more methods here,
	// and implement the added methods in customLeaseApprovalsModel.
	LeaseApprovalsModel interface {
		leaseApprovalsModel
		// 自定义方法
		FindByApplicationId(ctx context.Context, applicationId int64) ([]*LeaseApprovals, error)
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

// FindByApplicationId 根据申请ID查询审批记录
func (m *customLeaseApprovalsModel) FindByApplicationId(ctx context.Context, applicationId int64) ([]*LeaseApprovals, error) {
	query := fmt.Sprintf("SELECT %s FROM %s WHERE `application_id` = ? ORDER BY created_at ASC", leaseApprovalsRows, m.table)

	var approvals []*LeaseApprovals
	err := m.QueryRowsNoCacheCtx(ctx, &approvals, query, applicationId)
	if err != nil {
		return nil, err
	}

	return approvals, nil
}
