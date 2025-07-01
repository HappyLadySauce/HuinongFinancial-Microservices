package logic

import (
	"context"

	"loanproductrpc/internal/svc"
	"loanproductrpc/loanproduct"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetLoanProductLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetLoanProductLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetLoanProductLogic {
	return &GetLoanProductLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 产品查询
func (l *GetLoanProductLogic) GetLoanProduct(in *loanproduct.GetLoanProductReq) (*loanproduct.GetLoanProductResp, error) {
	// todo: add your logic here and delete this line

	return &loanproduct.GetLoanProductResp{}, nil
}
