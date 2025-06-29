package logic

import (
	"context"

	"rpc/internal/svc"
	"rpc/loan"

	"github.com/zeromicro/go-zero/core/logx"
)

type ListLoanApplicationsLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewListLoanApplicationsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ListLoanApplicationsLogic {
	return &ListLoanApplicationsLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *ListLoanApplicationsLogic) ListLoanApplications(in *loan.ListLoanApplicationsReq) (*loan.ListLoanApplicationsResp, error) {
	// todo: add your logic here and delete this line

	return &loan.ListLoanApplicationsResp{}, nil
}
