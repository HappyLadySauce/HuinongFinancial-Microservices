package logic

import (
	"context"

	"rpc/appuser"
	"rpc/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
)

type LogoutLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewLogoutLogic(ctx context.Context, svcCtx *svc.ServiceContext) *LogoutLogic {
	return &LogoutLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *LogoutLogic) Logout(in *appuser.LogoutReq) (*appuser.LogoutResp, error) {
	// todo: add your logic here and delete this line

	return &appuser.LogoutResp{}, nil
}
