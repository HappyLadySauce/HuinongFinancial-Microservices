package auth_jwt

import (
	"context"

	"api/internal/svc"
	"api/internal/types"
	"rpc/oauserclient"

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

func (l *ChangePasswordLogic) ChangePassword(req *types.ChangePasswordReq) (resp *types.ChangePasswordResp, err error) {
	// 调用 RPC 服务修改密码
	_, err = l.svcCtx.OaUserRpc.ChangePassword(l.ctx, &oauserclient.ChangePasswordReq{
		Phone:       req.Phone,
		OldPassword: req.OldPassword,
		NewPassword: req.NewPassword,
	})
	if err != nil {
		l.Logger.Errorf("RPC ChangePassword failed: %v", err)
		return nil, err
	}

	// 返回空结构体
	return &types.ChangePasswordResp{}, nil
}
