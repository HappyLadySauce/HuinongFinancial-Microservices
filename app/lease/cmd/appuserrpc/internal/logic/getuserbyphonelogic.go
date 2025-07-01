package logic

import (
	"context"

	"appuserrpc/appuser"
	"appuserrpc/internal/svc"

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
func (l *GetUserByPhoneLogic) GetUserByPhone(in *appuser.GetUserInfoReq) (*appuser.GetUserInfoResp, error) {
	// todo: add your logic here and delete this line

	return &appuser.GetUserInfoResp{}, nil
}
