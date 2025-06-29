package logic

import (
	"context"

	"rpc/internal/svc"
	"rpc/loan"

	"github.com/zeromicro/go-zero/core/logx"
)

type UpdateLoanApplicationLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewUpdateLoanApplicationLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateLoanApplicationLogic {
	return &UpdateLoanApplicationLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *UpdateLoanApplicationLogic) UpdateLoanApplication(in *loan.UpdateLoanApplicationReq) (*loan.UpdateLoanApplicationResp, error) {
	// todo: add your logic here and delete this line

	return &loan.UpdateLoanApplicationResp{}, nil
}
