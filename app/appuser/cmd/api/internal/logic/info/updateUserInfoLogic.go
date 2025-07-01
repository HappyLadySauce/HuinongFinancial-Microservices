package info

import (
	"context"

	"api/internal/breaker"
	"api/internal/svc"
	"api/internal/types"
	"rpc/appuser"

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

// 更新用户信息
func (l *UpdateUserInfoLogic) UpdateUserInfo(req *types.UpdateUserInfoReq) (resp *types.UpdateUserInfoResp, err error) {
	logx.WithContext(l.ctx).Infof("API: 更新用户信息请求, user_id: %d", req.UserInfo.Id)

	// 数据转换：API UserInfo -> RPC UserInfo
	rpcUserInfo := &appuser.UserInfo{
		Id:         req.UserInfo.Id,
		Phone:      req.UserInfo.Phone,
		Name:       req.UserInfo.Name,
		Nickname:   req.UserInfo.Nickname,
		Age:        int32(req.UserInfo.Age),
		Gender:     int32(req.UserInfo.Gender),
		Occupation: req.UserInfo.Occupation,
		Address:    req.UserInfo.Address,
		Income:     req.UserInfo.Income,
		CreatedAt:  req.UserInfo.CreatedAt,
		UpdatedAt:  req.UserInfo.UpdatedAt,
	}

	// 调用 RPC 更新用户信息服务 - 使用熔断器
	updateResp, err := breaker.DoWithBreakerResultAcceptable(l.ctx, "appuser-rpc", func() (*appuser.UpdateUserInfoResp, error) {
		return l.svcCtx.AppUserRpc.UpdateUserInfo(l.ctx, &appuser.UpdateUserInfoReq{
			UserInfo: rpcUserInfo,
		})
	}, breaker.IsAcceptableError)
	if err != nil {
		logx.WithContext(l.ctx).Errorf("RPC 更新用户信息调用失败: %v", err)
		return nil, err
	}

	// 数据转换：RPC UserInfo -> API UserInfo
	var userInfo types.UserInfo
	if updateResp.UserInfo != nil {
		userInfo = types.UserInfo{
			Id:         updateResp.UserInfo.Id,
			Phone:      updateResp.UserInfo.Phone,
			Name:       updateResp.UserInfo.Name,
			Nickname:   updateResp.UserInfo.Nickname,
			Age:        int(updateResp.UserInfo.Age),
			Gender:     int(updateResp.UserInfo.Gender),
			Occupation: updateResp.UserInfo.Occupation,
			Address:    updateResp.UserInfo.Address,
			Income:     updateResp.UserInfo.Income,
			CreatedAt:  updateResp.UserInfo.CreatedAt,
			UpdatedAt:  updateResp.UserInfo.UpdatedAt,
		}
	}

	// 转换 RPC 响应为 API 响应 - 只返回 UserInfo
	return &types.UpdateUserInfoResp{
		UserInfo: userInfo,
	}, nil
}
