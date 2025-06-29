package logic

import (
	"context"

	"rpc/internal/svc"
	"rpc/loanproduct"

	"github.com/zeromicro/go-zero/core/logx"
)

type UpdateLoanProductLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewUpdateLoanProductLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateLoanProductLogic {
	return &UpdateLoanProductLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *UpdateLoanProductLogic) UpdateLoanProduct(in *loanproduct.UpdateLoanProductReq) (*loanproduct.UpdateLoanProductResp, error) {
	// todo: add your logic here and delete this line

	return &loanproduct.UpdateLoanProductResp{}, nil
}
