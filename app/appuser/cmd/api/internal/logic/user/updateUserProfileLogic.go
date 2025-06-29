package user

import (
	"context"

	"api/internal/svc"
	"api/internal/types"
	"api/internal/utils"
	"rpc/appuser"

	"github.com/zeromicro/go-zero/core/logx"
)

type UpdateUserProfileLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewUpdateUserProfileLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateUserProfileLogic {
	return &UpdateUserProfileLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UpdateUserProfileLogic) UpdateUserProfile(req *types.UpdateProfileReq) (resp *types.DataResp, err error) {
	// 从JWT中获取用户ID
	userId, err := utils.GetUserIdFromCtx(l.ctx)
	if err != nil {
		return &types.DataResp{
			Code:    401,
			Message: err.Error(),
		}, nil
	}

	// 调用RPC更新用户档案
	updatedUser, err := l.svcCtx.AppUserRpc.UpdateUserProfile(l.ctx, &appuser.UpdateUserProfileReq{
		UserId:     userId,
		Nickname:   req.Nickname,
		Age:        int32(req.Age),
		Gender:     int32(req.Gender),
		Occupation: req.Occupation,
		Address:    req.Address,
		Income:     req.Income,
	})
	if err != nil {
		l.Errorf("调用RPC更新用户档案失败: %v", err)
		return &types.DataResp{
			Code:    500,
			Message: "更新用户档案失败",
		}, nil
	}

	// 返回成功响应
	return &types.DataResp{
		Code:    200,
		Message: "更新成功",
		Data:    updatedUser,
	}, nil
}
