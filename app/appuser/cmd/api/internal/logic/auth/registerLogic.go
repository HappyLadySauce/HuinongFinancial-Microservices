package auth

import (
	"context"

	"api/internal/svc"
	"api/internal/types"
	"rpc/appuser"

	"github.com/zeromicro/go-zero/core/logx"
)

type RegisterLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewRegisterLogic(ctx context.Context, svcCtx *svc.ServiceContext) *RegisterLogic {
	return &RegisterLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

// 用户注册
func (l *RegisterLogic) Register(req *types.RegisterReq) (resp *types.RegisterResp, err error) {
	logx.WithContext(l.ctx).Infof("API: 用户注册请求, phone: %s", req.Phone)

	// 调用 RPC 注册服务
	registerResp, err := l.svcCtx.AppUserRpc.Register(l.ctx, &appuser.RegisterReq{
		Phone:    req.Phone,
		Password: req.Password,
	})
	if err != nil {
		logx.WithContext(l.ctx).Errorf("RPC 注册调用失败: %v", err)
		return nil, err
	}

	// 转换 RPC 响应为 API 响应 - 只返回 token
	return &types.RegisterResp{
		Token: registerResp.Token,
	}, nil
}
