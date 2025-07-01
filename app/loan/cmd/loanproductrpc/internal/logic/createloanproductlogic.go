package logic

import (
	"context"

	"loanproductrpc/internal/svc"
	"loanproductrpc/loanproduct"

	"github.com/zeromicro/go-zero/core/logx"
)

type CreateLoanProductLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewCreateLoanProductLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CreateLoanProductLogic {
	return &CreateLoanProductLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 产品管理
func (l *CreateLoanProductLogic) CreateLoanProduct(in *loanproduct.CreateLoanProductReq) (*loanproduct.CreateLoanProductResp, error) {
	// todo: add your logic here and delete this line

	return &loanproduct.CreateLoanProductResp{}, nil
}
