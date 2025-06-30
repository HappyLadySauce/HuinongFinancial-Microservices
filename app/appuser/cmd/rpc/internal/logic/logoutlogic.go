package logic

import (
	"context"

	"rpc/appuser"
	"rpc/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
)

type LogoutLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewLogoutLogic(ctx context.Context, svcCtx *svc.ServiceContext) *LogoutLogic {
	return &LogoutLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 用户登出
func (l *LogoutLogic) Logout(in *appuser.LogoutReq) (*appuser.LogoutResp, error) {
	l.Info("用户登出请求")

	// 注：现在从 JWT 认证上下文获取用户信息，不再需要验证 token 参数
	// 在真实场景中，可以考虑以下操作：
	// 1. 将当前用户的 token 添加到黑名单 (Redis)
	// 2. 更新用户最后登出时间
	// 3. 清理用户会话相关缓存

	l.Infof("用户登出成功")
	return &appuser.LogoutResp{}, nil
}
