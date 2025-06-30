package model

import (
	"context"
	"fmt"

	"github.com/zeromicro/go-zero/core/stores/cache"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

var _ LoanApprovalsModel = (*customLoanApprovalsModel)(nil)

type (
	// LoanApprovalsModel is an interface to be customized, add more methods here,
	// and implement the added methods in customLoanApprovalsModel.
	LoanApprovalsModel interface {
		loanApprovalsModel
		// 自定义方法
		FindByApplicationId(ctx context.Context, applicationId int64) ([]*LoanApprovals, error)
	}

	customLoanApprovalsModel struct {
		*defaultLoanApprovalsModel
	}
)

// NewLoanApprovalsModel returns a model for the database table.
func NewLoanApprovalsModel(conn sqlx.SqlConn, c cache.CacheConf, opts ...cache.Option) LoanApprovalsModel {
	return &customLoanApprovalsModel{
		defaultLoanApprovalsModel: newLoanApprovalsModel(conn, c, opts...),
	}
}

// FindByApplicationId 根据申请ID查询审批记录
func (m *customLoanApprovalsModel) FindByApplicationId(ctx context.Context, applicationId int64) ([]*LoanApprovals, error) {
	query := fmt.Sprintf("SELECT %s FROM %s WHERE `application_id` = ? ORDER BY created_at ASC", loanApprovalsRows, m.table)

	var approvals []*LoanApprovals
	err := m.QueryRowsNoCacheCtx(ctx, &approvals, query, applicationId)
	if err != nil {
		return nil, err
	}

	return approvals, nil
}
