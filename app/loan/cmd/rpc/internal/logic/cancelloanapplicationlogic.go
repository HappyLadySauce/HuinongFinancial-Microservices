package logic

import (
	"context"

	"rpc/internal/svc"
	"rpc/loan"

	"github.com/zeromicro/go-zero/core/logx"
)

type CancelLoanApplicationLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewCancelLoanApplicationLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CancelLoanApplicationLogic {
	return &CancelLoanApplicationLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *CancelLoanApplicationLogic) CancelLoanApplication(in *loan.CancelLoanApplicationReq) (*loan.CancelLoanApplicationResp, error) {
	// todo: add your logic here and delete this line

	return &loan.CancelLoanApplicationResp{}, nil
}
