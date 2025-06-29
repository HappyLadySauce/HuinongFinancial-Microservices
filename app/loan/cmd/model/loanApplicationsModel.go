package model

import (
	"github.com/zeromicro/go-zero/core/stores/cache"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

var _ LoanApplicationsModel = (*customLoanApplicationsModel)(nil)

type (
	// LoanApplicationsModel is an interface to be customized, add more methods here,
	// and implement the added methods in customLoanApplicationsModel.
	LoanApplicationsModel interface {
		loanApplicationsModel
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
