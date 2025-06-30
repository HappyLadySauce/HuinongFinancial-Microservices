package model

import (
	"context"
	"fmt"

	"github.com/zeromicro/go-zero/core/stores/cache"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

var _ LoanApplicationsModel = (*customLoanApplicationsModel)(nil)

type (
	// LoanApplicationsModel is an interface to be customized, add more methods here,
	// and implement the added methods in customLoanApplicationsModel.
	LoanApplicationsModel interface {
		loanApplicationsModel
		// 自定义方法
		CountWithConditions(ctx context.Context, whereClause string, args []interface{}) (int64, error)
		ListWithConditions(ctx context.Context, whereClause string, args []interface{}, limit, offset int32) ([]*LoanApplications, error)
	}

	customLoanApplicationsModel struct {
		*defaultLoanApplicationsModel
	}
)

// NewLoanApplicationsModel returns a model for the database table.
func NewLoanApplicationsModel(conn sqlx.SqlConn, c cache.CacheConf, opts ...cache.Option) LoanApplicationsModel {
	return &customLoanApplicationsModel{
		defaultLoanApplicationsModel: newLoanApplicationsModel(conn, c, opts...),
	}
}

// CountWithConditions 根据条件统计申请数量
func (m *customLoanApplicationsModel) CountWithConditions(ctx context.Context, whereClause string, args []interface{}) (int64, error) {
	var total int64
	query := fmt.Sprintf("SELECT COUNT(*) FROM %s %s", m.table, whereClause)
	err := m.QueryRowNoCacheCtx(ctx, &total, query, args...)
	return total, err
}

// ListWithConditions 根据条件查询申请列表
func (m *customLoanApplicationsModel) ListWithConditions(ctx context.Context, whereClause string, args []interface{}, limit, offset int32) ([]*LoanApplications, error) {
	query := fmt.Sprintf("SELECT %s FROM %s %s ORDER BY created_at DESC LIMIT ? OFFSET ?", loanApplicationsRows, m.table, whereClause)
	queryArgs := append(args, limit, offset)
	
	var applications []*LoanApplications
	err := m.QueryRowsNoCacheCtx(ctx, &applications, query, queryArgs...)
	if err != nil {
		return nil, err
	}
	
	return applications, nil
}
