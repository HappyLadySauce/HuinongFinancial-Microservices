package auth_jwt

import (
	"context"

	"api/internal/breaker"
	"api/internal/svc"
	"api/internal/types"
	"rpc/appuser"

	"github.com/zeromicro/go-zero/core/logx"
)

type LogoutLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewLogoutLogic(ctx context.Context, svcCtx *svc.ServiceContext) *LogoutLogic {
	return &LogoutLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

// 用户登出
func (l *LogoutLogic) Logout(req *types.LogoutReq) (resp *types.LogoutResp, err error) {
	logx.WithContext(l.ctx).Infof("API: 用户登出请求")

	// 调用 RPC 登出服务 - 使用熔断器
	_, err = breaker.DoWithBreakerResultAcceptable(l.ctx, "appuser-rpc", func() (*appuser.LogoutResp, error) {
		return l.svcCtx.AppUserRpc.Logout(l.ctx, &appuser.LogoutReq{})
	}, breaker.IsAcceptableError)
	if err != nil {
		logx.WithContext(l.ctx).Errorf("RPC 用户登出调用失败: %v", err)
		return nil, err
	}

	// 返回空结构体
	return &types.LogoutResp{}, nil
}
