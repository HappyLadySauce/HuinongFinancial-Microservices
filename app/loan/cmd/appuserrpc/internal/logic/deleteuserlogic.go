package logic

import (
	"context"

	"appuserrpc/appuser"
	"appuserrpc/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
)

type DeleteUserLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewDeleteUserLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DeleteUserLogic {
	return &DeleteUserLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *DeleteUserLogic) DeleteUser(in *appuser.DeleteUserReq) (*appuser.DeleteUserResp, error) {
	// todo: add your logic here and delete this line

	return &appuser.DeleteUserResp{}, nil
}
