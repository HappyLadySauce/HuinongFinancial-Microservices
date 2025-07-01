package info

import (
	"context"

	"api/internal/breaker"
	"api/internal/svc"
	"api/internal/types"
	"rpc/oauserclient"

	"github.com/zeromicro/go-zero/core/logx"
)

type UpdateUserInfoLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewUpdateUserInfoLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateUserInfoLogic {
	return &UpdateUserInfoLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UpdateUserInfoLogic) UpdateUserInfo(req *types.UpdateUserInfoReq) (resp *types.UpdateUserInfoResp, err error) {
	// 数据转换：API UserInfo -> RPC UserInfo
	rpcUserInfo := &oauserclient.UserInfo{
		Id:        req.UserInfo.Id,
		Phone:     req.UserInfo.Phone,
		Name:      req.UserInfo.Name,
		Nickname:  req.UserInfo.Nickname,
		Age:       int32(req.UserInfo.Age),
		Gender:    int32(req.UserInfo.Gender),
		Role:      req.UserInfo.Role,
		Status:    int32(req.UserInfo.Status),
		CreatedAt: req.UserInfo.CreatedAt,
		UpdatedAt: req.UserInfo.UpdatedAt,
	}

	// 调用 RPC 服务更新用户信息 - 使用熔断器
	updateResp, err := breaker.DoWithBreakerResultAcceptable(l.ctx, "oauser-rpc", func() (*oauserclient.UpdateUserInfoResp, error) {
		return l.svcCtx.OaUserRpc.UpdateUserInfo(l.ctx, &oauserclient.UpdateUserInfoReq{
			UserInfo: rpcUserInfo,
		})
	}, breaker.IsAcceptableError)
	if err != nil {
		l.Logger.Errorf("RPC UpdateUserInfo failed: %v", err)
		return nil, err
	}

	// 数据转换：RPC UserInfo -> API UserInfo
	var userInfo types.UserInfo
	if updateResp.UserInfo != nil {
		userInfo = types.UserInfo{
			Id:        updateResp.UserInfo.Id,
			Phone:     updateResp.UserInfo.Phone,
			Name:      updateResp.UserInfo.Name,
			Nickname:  updateResp.UserInfo.Nickname,
			Age:       int(updateResp.UserInfo.Age),
			Gender:    int(updateResp.UserInfo.Gender),
			Role:      updateResp.UserInfo.Role,
			Status:    int(updateResp.UserInfo.Status),
			CreatedAt: updateResp.UserInfo.CreatedAt,
			UpdatedAt: updateResp.UserInfo.UpdatedAt,
		}
	}

	// 转换响应格式
	return &types.UpdateUserInfoResp{
		UserInfo: userInfo,
	}, nil
}
