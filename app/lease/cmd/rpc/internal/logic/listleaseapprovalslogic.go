package logic

import (
	"context"

	"rpc/internal/svc"
	"rpc/lease"

	"github.com/zeromicro/go-zero/core/logx"
)

type ListLeaseApprovalsLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewListLeaseApprovalsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ListLeaseApprovalsLogic {
	return &ListLeaseApprovalsLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *ListLeaseApprovalsLogic) ListLeaseApprovals(in *lease.ListLeaseApprovalsReq) (*lease.ListLeaseApprovalsResp, error) {
	// todo: add your logic here and delete this line

	return &lease.ListLeaseApprovalsResp{}, nil
}
