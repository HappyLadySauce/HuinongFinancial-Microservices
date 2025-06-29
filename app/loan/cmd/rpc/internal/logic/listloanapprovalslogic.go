package logic

import (
	"context"

	"rpc/internal/svc"
	"rpc/loan"

	"github.com/zeromicro/go-zero/core/logx"
)

type ListLoanApprovalsLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewListLoanApprovalsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ListLoanApprovalsLogic {
	return &ListLoanApprovalsLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *ListLoanApprovalsLogic) ListLoanApprovals(in *loan.ListLoanApprovalsReq) (*loan.ListLoanApprovalsResp, error) {
	// todo: add your logic here and delete this line

	return &loan.ListLoanApprovalsResp{}, nil
}
