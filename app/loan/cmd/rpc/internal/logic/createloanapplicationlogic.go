package logic

import (
	"context"

	"rpc/internal/svc"
	"rpc/loan"

	"github.com/zeromicro/go-zero/core/logx"
)

type CreateLoanApplicationLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewCreateLoanApplicationLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CreateLoanApplicationLogic {
	return &CreateLoanApplicationLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 贷款申请管理
func (l *CreateLoanApplicationLogic) CreateLoanApplication(in *loan.CreateLoanApplicationReq) (*loan.CreateLoanApplicationResp, error) {
	// todo: add your logic here and delete this line

	return &loan.CreateLoanApplicationResp{}, nil
}
