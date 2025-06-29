package info

import (
	"context"

	"api/internal/svc"
	"api/internal/types"
	"rpc/appuser"

	"github.com/zeromicro/go-zero/core/logx"
)

type DeleteUserLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewDeleteUserLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DeleteUserLogic {
	return &DeleteUserLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

// 删除用户（软删除）
func (l *DeleteUserLogic) DeleteUser(req *types.DeleteUserReq) (resp *types.DeleteUserResp, err error) {
	logx.WithContext(l.ctx).Infof("API: 删除用户请求, phone: %s", req.Phone)

	// 调用 RPC 删除用户服务
	deleteResp, err := l.svcCtx.AppUserRpc.DeleteUser(l.ctx, &appuser.DeleteUserReq{
		Phone: req.Phone,
	})
	if err != nil {
		logx.WithContext(l.ctx).Errorf("RPC 删除用户调用失败: %v", err)
		return &types.DeleteUserResp{
			Code:    500,
			Message: "服务内部错误",
		}, nil
	}

	// 转换 RPC 响应为 API 响应
	return &types.DeleteUserResp{
		Code:    deleteResp.Code,
		Message: deleteResp.Message,
	}, nil
}
