package auth

import (
	"context"

	"api/internal/svc"
	"api/internal/types"
	"rpc/oauserclient"

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

func (l *RegisterLogic) Register(req *types.RegisterReq) (resp *types.RegisterResp, err error) {
	// 调用 RPC 服务进行注册
	registerResp, err := l.svcCtx.OaUserRpc.Register(l.ctx, &oauserclient.RegisterReq{
		Phone:    req.Phone,
		Password: req.Password,
	})
	if err != nil {
		l.Logger.Errorf("RPC Register failed: %v", err)
		return &types.RegisterResp{
			Code:    500,
			Message: "服务器内部错误",
		}, nil
	}

	// 转换响应格式
	return &types.RegisterResp{
		Code:    registerResp.Code,
		Message: registerResp.Message,
		Token:   registerResp.Token,
	}, nil
}
