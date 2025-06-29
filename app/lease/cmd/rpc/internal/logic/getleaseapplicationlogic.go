package logic

import (
	"context"

	"rpc/internal/svc"
	"rpc/lease"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetLeaseApplicationLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetLeaseApplicationLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetLeaseApplicationLogic {
	return &GetLeaseApplicationLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GetLeaseApplicationLogic) GetLeaseApplication(in *lease.GetLeaseApplicationReq) (*lease.GetLeaseApplicationResp, error) {
	// todo: add your logic here and delete this line

	return &lease.GetLeaseApplicationResp{}, nil
}
