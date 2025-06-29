package logic

import (
	"context"

	"rpc/internal/svc"
	"rpc/loanproduct"

	"github.com/zeromicro/go-zero/core/logx"
)

type ListLoanProductsLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewListLoanProductsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ListLoanProductsLogic {
	return &ListLoanProductsLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *ListLoanProductsLogic) ListLoanProducts(in *loanproduct.ListLoanProductsReq) (*loanproduct.ListLoanProductsResp, error) {
	// todo: add your logic here and delete this line

	return &loanproduct.ListLoanProductsResp{}, nil
}
