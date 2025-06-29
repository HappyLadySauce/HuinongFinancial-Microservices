package logic

import (
	"context"

	"rpc/internal/svc"
	"rpc/loan"

	"github.com/zeromicro/go-zero/core/logx"
)

type ApproveLoanApplicationLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewApproveLoanApplicationLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ApproveLoanApplicationLogic {
	return &ApproveLoanApplicationLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 贷款审批管理
func (l *ApproveLoanApplicationLogic) ApproveLoanApplication(in *loan.ApproveLoanApplicationReq) (*loan.ApproveLoanApplicationResp, error) {
	// todo: add your logic here and delete this line

	return &loan.ApproveLoanApplicationResp{}, nil
}
