package info

import (
	"context"

	"api/internal/svc"
	"api/internal/types"
	"rpc/oauserclient"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetUserByPhoneLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetUserByPhoneLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetUserByPhoneLogic {
	return &GetUserByPhoneLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetUserByPhoneLogic) GetUserByPhone(req *types.GetUserInfoReq) (resp *types.GetUserInfoResp, err error) {
	// 调用 RPC 服务获取用户信息
	getUserResp, err := l.svcCtx.OaUserRpc.GetUserByPhone(l.ctx, &oauserclient.GetUserInfoReq{
		Phone: req.Phone,
	})
	if err != nil {
		l.Logger.Errorf("RPC GetUserByPhone failed: %v", err)
		return &types.GetUserInfoResp{
			Code:    500,
			Message: "服务器内部错误",
		}, nil
	}

	// 转换用户信息格式
	var userInfo types.UserInfo
	if getUserResp.UserInfo != nil {
		userInfo = types.UserInfo{
			Id:        getUserResp.UserInfo.Id,
			Phone:     getUserResp.UserInfo.Phone,
			Name:      getUserResp.UserInfo.Name,
			Nickname:  getUserResp.UserInfo.Nickname,
			Age:       int(getUserResp.UserInfo.Age),
			Gender:    int(getUserResp.UserInfo.Gender),
			Role:      getUserResp.UserInfo.Role,
			Status:    int(getUserResp.UserInfo.Status),
			CreatedAt: getUserResp.UserInfo.CreatedAt,
			UpdatedAt: getUserResp.UserInfo.UpdatedAt,
		}
	}

	// 转换响应格式
	return &types.GetUserInfoResp{
		Code:     getUserResp.Code,
		Message:  getUserResp.Message,
		UserInfo: userInfo,
	}, nil
}
