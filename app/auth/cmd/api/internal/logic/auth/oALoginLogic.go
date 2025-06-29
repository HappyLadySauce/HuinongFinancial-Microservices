package auth

import (
	"context"
	"fmt"
	"rpc/authclient"

	"api/internal/svc"
	"api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type OALoginLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewOALoginLogic(ctx context.Context, svcCtx *svc.ServiceContext) *OALoginLogic {
	return &OALoginLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *OALoginLogic) OALogin(req *types.OALoginReq) (resp *types.DataResp, err error) {
	// 1. 参数验证
	if req.Username == "" || req.Password == "" {
		return nil, fmt.Errorf("用户名和密码不能为空")
	}

	// 2. 调用Auth RPC进行管理员登录验证
	rpcResp, err := l.svcCtx.AuthRpc.OALogin(l.ctx, &authclient.OALoginReq{
		Username: req.Username,
		Password: req.Password,
	})
	if err != nil {
		logx.Errorf("Auth RPC管理员登录失败: %v", err)
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
