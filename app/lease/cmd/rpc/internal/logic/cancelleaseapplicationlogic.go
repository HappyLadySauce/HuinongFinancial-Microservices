package logic

import (
	"context"

	"rpc/internal/svc"
	"rpc/lease"

	"github.com/zeromicro/go-zero/core/logx"
)

type CancelLeaseApplicationLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewCancelLeaseApplicationLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CancelLeaseApplicationLogic {
	return &CancelLeaseApplicationLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *CancelLeaseApplicationLogic) CancelLeaseApplication(in *lease.CancelLeaseApplicationReq) (*lease.CancelLeaseApplicationResp, error) {
	// todo: add your logic here and delete this line

	return &lease.CancelLeaseApplicationResp{}, nil
}
