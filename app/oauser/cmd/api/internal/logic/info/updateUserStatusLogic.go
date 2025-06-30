package info

import (
	"context"

	"api/internal/svc"
	"api/internal/types"
	"rpc/oauserclient"

	"github.com/zeromicro/go-zero/core/logx"
)

type UpdateUserStatusLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewUpdateUserStatusLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateUserStatusLogic {
	return &UpdateUserStatusLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UpdateUserStatusLogic) UpdateUserStatus(req *types.UpdateUserStatusReq) (resp *types.UpdateUserStatusResp, err error) {
	// 调用 RPC 服务更新用户状态
	updateResp, err := l.svcCtx.OaUserRpc.UpdateUserStatus(l.ctx, &oauserclient.UpdateUserStatusReq{
		Phone:  req.Phone,
		Status: int32(req.Status),
	})
	if err != nil {
		l.Logger.Errorf("RPC UpdateUserStatus failed: %v", err)
		return nil, err
	}

	// 返回响应
	return &types.UpdateUserStatusResp{
		Status: int(updateResp.Status),
	}, nil
}
