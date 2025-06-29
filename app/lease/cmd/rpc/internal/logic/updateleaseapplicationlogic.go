package logic

import (
	"context"

	"rpc/internal/svc"
	"rpc/lease"

	"github.com/zeromicro/go-zero/core/logx"
)

type UpdateLeaseApplicationLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewUpdateLeaseApplicationLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateLeaseApplicationLogic {
	return &UpdateLeaseApplicationLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *UpdateLeaseApplicationLogic) UpdateLeaseApplication(in *lease.UpdateLeaseApplicationReq) (*lease.UpdateLeaseApplicationResp, error) {
	// todo: add your logic here and delete this line

	return &lease.UpdateLeaseApplicationResp{}, nil
}
