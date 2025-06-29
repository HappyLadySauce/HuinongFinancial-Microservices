package model

import (
	"github.com/zeromicro/go-zero/core/stores/cache"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

var _ LoanProductsModel = (*customLoanProductsModel)(nil)

type (
	// LoanProductsModel is an interface to be customized, add more methods here,
	// and implement the added methods in customLoanProductsModel.
	LoanProductsModel interface {
		loanProductsModel
	}

	customLoanProductsModel struct {
		*defaultLoanProductsModel
	}
)

// NewLoanProductsModel returns a model for the database table.
func NewLoanProductsModel(conn sqlx.SqlConn, c cache.CacheConf, opts ...cache.Option) LoanProductsModel {
	return &customLoanProductsModel{
		defaultLoanProductsModel: newLoanProductsModel(conn, c, opts...),
	}
}
