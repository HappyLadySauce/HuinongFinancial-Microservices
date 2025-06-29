package user

import (
	"context"

	"api/internal/svc"
	"api/internal/types"
	"api/internal/utils"
	"rpc/appuser"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetUserInfoLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetUserInfoLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetUserInfoLogic {
	return &GetUserInfoLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetUserInfoLogic) GetUserInfo() (resp *types.DataResp, err error) {
	// 从JWT中获取用户ID
	userId, err := utils.GetUserIdFromCtx(l.ctx)
	if err != nil {
		return &types.DataResp{
			Code:    401,
			Message: err.Error(),
		}, nil
	}

	// 调用RPC获取用户信息
	userInfo, err := l.svcCtx.AppUserRpc.GetUserInfo(l.ctx, &appuser.GetUserInfoReq{
		UserId: userId,
	})
	if err != nil {
		l.Errorf("调用RPC获取用户信息失败: %v", err)
		return &types.DataResp{
			Code:    500,
			Message: "获取用户信息失败",
		}, nil
	}

	// 返回成功响应
	return &types.DataResp{
		Code:    200,
		Message: "success",
		Data:    userInfo,
	}, nil
}
