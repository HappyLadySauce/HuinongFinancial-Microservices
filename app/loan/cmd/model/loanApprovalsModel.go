package model

import (
	"github.com/zeromicro/go-zero/core/stores/cache"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

var _ LoanApprovalsModel = (*customLoanApprovalsModel)(nil)

type (
	// LoanApprovalsModel is an interface to be customized, add more methods here,
	// and implement the added methods in customLoanApprovalsModel.
	LoanApprovalsModel interface {
		loanApprovalsModel
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
