package info

import (
	"context"

	"api/internal/breaker"
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

	// 调用 RPC 删除用户服务 - 使用熔断器
	_, err = breaker.DoWithBreakerResultAcceptable(l.ctx, "appuser-rpc", func() (*appuser.DeleteUserResp, error) {
		return l.svcCtx.AppUserRpc.DeleteUser(l.ctx, &appuser.DeleteUserReq{
			Phone: req.Phone,
		})
	}, breaker.IsAcceptableError)
	if err != nil {
		logx.WithContext(l.ctx).Errorf("RPC 删除用户调用失败: %v", err)
		return nil, err
	}

	// 返回空结构体
	return &types.DeleteUserResp{}, nil
}
