package logic

import (
	"context"

	"appuserrpc/appuser"
	"appuserrpc/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
)

type ChangePasswordLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewChangePasswordLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ChangePasswordLogic {
	return &ChangePasswordLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *ChangePasswordLogic) ChangePassword(in *appuser.ChangePasswordReq) (*appuser.ChangePasswordResp, error) {
	// todo: add your logic here and delete this line

	return &appuser.ChangePasswordResp{}, nil
}
