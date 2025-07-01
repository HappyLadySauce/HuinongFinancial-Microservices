package auth_jwt

import (
	"context"

	"api/internal/breaker"
	"api/internal/svc"
	"api/internal/types"
	"rpc/appuser"

	"github.com/zeromicro/go-zero/core/logx"
)

type ChangePasswordLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewChangePasswordLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ChangePasswordLogic {
	return &ChangePasswordLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

// 修改密码
func (l *ChangePasswordLogic) ChangePassword(req *types.ChangePasswordReq) (resp *types.ChangePasswordResp, err error) {
	logx.WithContext(l.ctx).Infof("API: 修改密码请求, phone: %s", req.Phone)

	// 调用 RPC 修改密码服务 - 使用熔断器
	_, err = breaker.DoWithBreakerResultAcceptable(l.ctx, "appuser-rpc", func() (*appuser.ChangePasswordResp, error) {
		return l.svcCtx.AppUserRpc.ChangePassword(l.ctx, &appuser.ChangePasswordReq{
			Phone:       req.Phone,
			OldPassword: req.OldPassword,
			NewPassword: req.NewPassword,
		})
	}, breaker.IsAcceptableError)
	if err != nil {
		logx.WithContext(l.ctx).Errorf("RPC 修改密码调用失败: %v", err)
		return nil, err
	}

	// 返回空结构体
	return &types.ChangePasswordResp{}, nil
}
