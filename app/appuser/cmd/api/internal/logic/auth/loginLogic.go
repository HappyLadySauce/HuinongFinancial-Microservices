package auth

import (
	"context"

	"api/internal/breaker"
	"api/internal/svc"
	"api/internal/types"
	"rpc/appuser"

	"github.com/zeromicro/go-zero/core/logx"
)

type LoginLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewLoginLogic(ctx context.Context, svcCtx *svc.ServiceContext) *LoginLogic {
	return &LoginLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

// 用户登录
func (l *LoginLogic) Login(req *types.LoginReq) (resp *types.LoginResp, err error) {
	logx.WithContext(l.ctx).Infof("API: 用户登录请求, phone: %s", req.Phone)

	// 使用熔断器调用 RPC 登录服务
	loginResp, err := breaker.DoWithBreakerResultAcceptable(l.ctx, "appuser-rpc", func() (*appuser.LoginResp, error) {
		return l.svcCtx.AppUserRpc.Login(l.ctx, &appuser.LoginReq{
			Phone:    req.Phone,
			Password: req.Password,
		})
	}, breaker.IsAcceptableError)

	if err != nil {
		logx.WithContext(l.ctx).Errorf("RPC 登录调用失败: %v", err)
		return nil, err
	}

	// 转换 RPC 响应为 API 响应 - 只返回 token
	return &types.LoginResp{
		Token: loginResp.Token,
	}, nil
}
