package auth_jwt

import (
	"context"

	"api/internal/breaker"
	"api/internal/svc"
	"api/internal/types"
	"rpc/oauserclient"

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

func (l *LogoutLogic) Logout(req *types.LogoutReq) (resp *types.LogoutResp, err error) {
	// 调用 RPC 服务进行登出 - 使用熔断器
	_, err = breaker.DoWithBreakerResultAcceptable(l.ctx, "oauser-rpc", func() (*oauserclient.LogoutResp, error) {
		return l.svcCtx.OaUserRpc.Logout(l.ctx, &oauserclient.LogoutReq{})
	}, breaker.IsAcceptableError)
	if err != nil {
		l.Logger.Errorf("RPC Logout failed: %v", err)
		return nil, err
	}

	// 返回空结构体
	return &types.LogoutResp{}, nil
}
