package logic

import (
	"context"

	"rpc/internal/svc"
	"rpc/lease"

	"github.com/zeromicro/go-zero/core/logx"
)

type CreateLeaseApplicationLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewCreateLeaseApplicationLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CreateLeaseApplicationLogic {
	return &CreateLeaseApplicationLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 租赁申请管理
func (l *CreateLeaseApplicationLogic) CreateLeaseApplication(in *lease.CreateLeaseApplicationReq) (*lease.CreateLeaseApplicationResp, error) {
	// todo: add your logic here and delete this line

	return &lease.CreateLeaseApplicationResp{}, nil
}
