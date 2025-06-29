package auth

import (
	"context"

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
	// 调用 RPC 服务进行登录
	loginResp, err := l.svcCtx.OaUserRpc.Login(l.ctx, &oauserclient.LoginReq{
		Phone:    req.Phone,
		Password: req.Password,
	})
	if err != nil {
		l.Logger.Errorf("RPC Login failed: %v", err)
		return &types.LoginResp{
			Code:    500,
			Message: "服务器内部错误",
		}, nil
	}

	// 转换响应格式
	return &types.LoginResp{
		Code:    loginResp.Code,
		Message: loginResp.Message,
		Token:   loginResp.Token,
	}, nil
}
