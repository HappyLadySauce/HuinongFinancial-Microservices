package logic

import (
	"context"

	"rpc/internal/svc"
	"rpc/oauser"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetUserByPhoneLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetUserByPhoneLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetUserByPhoneLogic {
	return &GetUserByPhoneLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 用户信息管理
func (l *GetUserByPhoneLogic) GetUserByPhone(in *oauser.GetUserInfoReq) (*oauser.GetUserInfoResp, error) {
	// todo: add your logic here and delete this line

	return &oauser.GetUserInfoResp{}, nil
}
