package logic

import (
	"context"

	"appuserrpc/appuser"
	"appuserrpc/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
)

type LoginLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewLoginLogic(ctx context.Context, svcCtx *svc.ServiceContext) *LoginLogic {
	return &LoginLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 用户认证管理
func (l *LoginLogic) Login(in *appuser.LoginReq) (*appuser.LoginResp, error) {
	// todo: add your logic here and delete this line

	return &appuser.LoginResp{}, nil
}
