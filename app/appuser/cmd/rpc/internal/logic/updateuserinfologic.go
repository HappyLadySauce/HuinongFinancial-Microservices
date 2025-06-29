package logic

import (
	"context"

	"rpc/appuser"
	"rpc/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
)

type UpdateUserInfoLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewUpdateUserInfoLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateUserInfoLogic {
	return &UpdateUserInfoLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *UpdateUserInfoLogic) UpdateUserInfo(in *appuser.UpdateUserInfoReq) (*appuser.UpdateUserInfoResp, error) {
	// todo: add your logic here and delete this line

	return &appuser.UpdateUserInfoResp{}, nil
}
