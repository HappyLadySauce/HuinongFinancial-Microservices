package auth_jwt

import (
	"context"

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
	logx.WithContext(l.ctx).Info("API: 用户登出请求")

	// 调用 RPC 登出服务（不需要传递 token，从 JWT 认证上下文获取）
	_, err = l.svcCtx.AppUserRpc.Logout(l.ctx, &appuser.LogoutReq{})
	if err != nil {
		logx.WithContext(l.ctx).Errorf("RPC 登出调用失败: %v", err)
		return nil, err
	}

	// 返回空结构体
	return &types.LogoutResp{}, nil
}
