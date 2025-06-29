package logic

import (
	"context"

	"rpc/appuser"
	"rpc/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
)

type RegisterLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewRegisterLogic(ctx context.Context, svcCtx *svc.ServiceContext) *RegisterLogic {
	return &RegisterLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *RegisterLogic) Register(in *appuser.RegisterReq) (*appuser.RegisterResp, error) {
	// todo: add your logic here and delete this line

	return &appuser.RegisterResp{}, nil
}
