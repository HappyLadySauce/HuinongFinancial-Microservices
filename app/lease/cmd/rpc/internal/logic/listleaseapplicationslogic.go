package logic

import (
	"context"

	"rpc/internal/svc"
	"rpc/lease"

	"github.com/zeromicro/go-zero/core/logx"
)

type ListLeaseApplicationsLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewListLeaseApplicationsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ListLeaseApplicationsLogic {
	return &ListLeaseApplicationsLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *ListLeaseApplicationsLogic) ListLeaseApplications(in *lease.ListLeaseApplicationsReq) (*lease.ListLeaseApplicationsResp, error) {
	// todo: add your logic here and delete this line

	return &lease.ListLeaseApplicationsResp{}, nil
}
