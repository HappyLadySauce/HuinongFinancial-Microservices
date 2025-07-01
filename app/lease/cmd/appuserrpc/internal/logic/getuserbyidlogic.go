package logic

import (
	"context"

	"appuserrpc/appuser"
	"appuserrpc/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetUserByIdLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetUserByIdLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetUserByIdLogic {
	return &GetUserByIdLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GetUserByIdLogic) GetUserById(in *appuser.GetUserByIdReq) (*appuser.GetUserInfoResp, error) {
	// todo: add your logic here and delete this line

	return &appuser.GetUserInfoResp{}, nil
}
