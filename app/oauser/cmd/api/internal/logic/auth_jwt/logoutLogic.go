package auth_jwt

import (
	"context"

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
	// 调用 RPC 服务进行注销
	logoutResp, err := l.svcCtx.OaUserRpc.Logout(l.ctx, &oauserclient.LogoutReq{
		Token: req.Token,
	})
	if err != nil {
		l.Logger.Errorf("RPC Logout failed: %v", err)
		return &types.LogoutResp{
			Code:    500,
			Message: "服务器内部错误",
		}, nil
	}

	// 转换响应格式
	return &types.LogoutResp{
		Code:    logoutResp.Code,
		Message: logoutResp.Message,
	}, nil
}
