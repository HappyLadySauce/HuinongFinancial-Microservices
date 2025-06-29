package logic

import (
	"context"

	"rpc/internal/svc"
	"rpc/loanproduct"

	"github.com/zeromicro/go-zero/core/logx"
)

type DeleteLoanProductLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewDeleteLoanProductLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DeleteLoanProductLogic {
	return &DeleteLoanProductLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *DeleteLoanProductLogic) DeleteLoanProduct(in *loanproduct.DeleteLoanProductReq) (*loanproduct.DeleteLoanProductResp, error) {
	// todo: add your logic here and delete this line

	return &loanproduct.DeleteLoanProductResp{}, nil
}
