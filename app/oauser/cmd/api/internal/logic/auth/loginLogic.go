package auth

import (
	"context"

	"api/internal/breaker"
	"api/internal/svc"
	"api/internal/types"
	"rpc/oauserclient"

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

func (l *LoginLogic) Login(req *types.LoginReq) (resp *types.LoginResp, err error) {
	// 调用 RPC 服务进行登录 - 使用熔断器
	loginResp, err := breaker.DoWithBreakerResultAcceptable(l.ctx, "oauser-rpc", func() (*oauserclient.LoginResp, error) {
		return l.svcCtx.OaUserRpc.Login(l.ctx, &oauserclient.LoginReq{
			Phone:    req.Phone,
			Password: req.Password,
		})
	}, breaker.IsAcceptableError)
	if err != nil {
		l.Logger.Errorf("RPC Login failed: %v", err)
		return nil, err
	}

	// 转换响应格式 - 只返回 token
	return &types.LoginResp{
		Token: loginResp.Token,
	}, nil
}
