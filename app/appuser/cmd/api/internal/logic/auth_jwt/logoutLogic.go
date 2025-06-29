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

	// 调用 RPC 登出服务
	logoutResp, err := l.svcCtx.AppUserRpc.Logout(l.ctx, &appuser.LogoutReq{
		Token: req.Token,
	})
	if err != nil {
		logx.WithContext(l.ctx).Errorf("RPC 登出调用失败: %v", err)
		return &types.LogoutResp{
			Code:    500,
			Message: "服务内部错误",
		}, nil
	}

	// 转换 RPC 响应为 API 响应
	return &types.LogoutResp{
		Code:    logoutResp.Code,
		Message: logoutResp.Message,
	}, nil
}
