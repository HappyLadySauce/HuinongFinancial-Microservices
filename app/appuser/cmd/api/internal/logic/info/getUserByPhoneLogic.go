package info

import (
	"context"

	"api/internal/svc"
	"api/internal/types"
	"rpc/appuser"

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

// 根据手机号获取用户信息
func (l *GetUserByPhoneLogic) GetUserByPhone(req *types.GetUserInfoReq) (resp *types.GetUserInfoResp, err error) {
	logx.WithContext(l.ctx).Infof("API: 获取用户信息请求, phone: %s", req.Phone)

	// 调用 RPC 获取用户信息服务
	userResp, err := l.svcCtx.AppUserRpc.GetUserByPhone(l.ctx, &appuser.GetUserInfoReq{
		Phone: req.Phone,
	})
	if err != nil {
		logx.WithContext(l.ctx).Errorf("RPC 获取用户信息调用失败: %v", err)
		return nil, err
	}

	// 数据转换：RPC UserInfo -> API UserInfo
	var userInfo types.UserInfo
	if userResp.UserInfo != nil {
		userInfo = types.UserInfo{
			Id:         userResp.UserInfo.Id,
			Phone:      userResp.UserInfo.Phone,
			Name:       userResp.UserInfo.Name,
			Nickname:   userResp.UserInfo.Nickname,
			Age:        int(userResp.UserInfo.Age),
			Gender:     int(userResp.UserInfo.Gender),
			Occupation: userResp.UserInfo.Occupation,
			Address:    userResp.UserInfo.Address,
			Income:     userResp.UserInfo.Income,
			CreatedAt:  userResp.UserInfo.CreatedAt,
			UpdatedAt:  userResp.UserInfo.UpdatedAt,
		}
	}

	// 转换 RPC 响应为 API 响应 - 只返回 UserInfo
	return &types.GetUserInfoResp{
		UserInfo: userInfo,
	}, nil
}
