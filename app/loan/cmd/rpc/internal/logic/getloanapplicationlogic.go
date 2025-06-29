package logic

import (
	"context"

	"rpc/internal/svc"
	"rpc/loan"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetLoanApplicationLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetLoanApplicationLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetLoanApplicationLogic {
	return &GetLoanApplicationLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GetLoanApplicationLogic) GetLoanApplication(in *loan.GetLoanApplicationReq) (*loan.GetLoanApplicationResp, error) {
	// todo: add your logic here and delete this line

	return &loan.GetLoanApplicationResp{}, nil
}
