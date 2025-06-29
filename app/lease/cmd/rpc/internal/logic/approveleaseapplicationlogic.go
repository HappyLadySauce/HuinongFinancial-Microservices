package logic

import (
	"context"

	"rpc/internal/svc"
	"rpc/lease"

	"github.com/zeromicro/go-zero/core/logx"
)

type ApproveLeaseApplicationLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewApproveLeaseApplicationLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ApproveLeaseApplicationLogic {
	return &ApproveLeaseApplicationLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 租赁审批管理
func (l *ApproveLeaseApplicationLogic) ApproveLeaseApplication(in *lease.ApproveLeaseApplicationReq) (*lease.ApproveLeaseApplicationResp, error) {
	// todo: add your logic here and delete this line

	return &lease.ApproveLeaseApplicationResp{}, nil
}
