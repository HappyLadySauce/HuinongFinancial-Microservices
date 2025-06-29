package auth_jwt

import (
	"context"

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

	// 调用 RPC 修改密码服务
	changeResp, err := l.svcCtx.AppUserRpc.ChangePassword(l.ctx, &appuser.ChangePasswordReq{
		Phone:       req.Phone,
		OldPassword: req.OldPassword,
		NewPassword: req.NewPassword,
	})
	if err != nil {
		logx.WithContext(l.ctx).Errorf("RPC 修改密码调用失败: %v", err)
		return &types.ChangePasswordResp{
			Code:    500,
			Message: "服务内部错误",
		}, nil
	}

	// 转换 RPC 响应为 API 响应
	return &types.ChangePasswordResp{
		Code:    changeResp.Code,
		Message: changeResp.Message,
	}, nil
}
