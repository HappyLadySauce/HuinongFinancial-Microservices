package auth

import (
	"context"
	"fmt"
	"rpc/authclient"

	"api/internal/svc"
	"api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type AppLoginLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewAppLoginLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AppLoginLogic {
	return &AppLoginLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *AppLoginLogic) AppLogin(req *types.AppLoginReq) (resp *types.DataResp, err error) {
	// 1. 参数验证
	if req.Account == "" || req.Password == "" {
		return nil, fmt.Errorf("账号和密码不能为空")
	}

	// 2. 调用Auth RPC进行登录验证
	rpcResp, err := l.svcCtx.AuthRpc.AppLogin(l.ctx, &authclient.AppLoginReq{
		Account:  req.Account,
		Password: req.Password,
	})
	if err != nil {
		logx.Errorf("Auth RPC登录失败: %v", err)
		return nil, fmt.Errorf("登录失败")
	}

	// 3. 组装响应数据
	loginData := map[string]interface{}{
		"accessToken":   rpcResp.AccessToken,
		"accessExpire":  rpcResp.AccessExpire,
		"refreshToken":  rpcResp.RefreshToken,
		"refreshExpire": rpcResp.RefreshExpire,
		"userId":        rpcResp.UserId,
		"userType":      rpcResp.UserType,
	}

	return &types.DataResp{
		Code:    200,
		Message: "登录成功",
		Data:    loginData,
	}, nil
}
