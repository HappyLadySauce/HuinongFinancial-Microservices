package info

import (
	"context"
	"net/http"

	"api/internal/breaker"
	"api/internal/svc"
	"api/internal/types"
	"rpc/oauserclient"

	"github.com/zeromicro/go-zero/core/logx"
)

type DeleteUserLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
	r      *http.Request
}

func NewDeleteUserLogic(ctx context.Context, svcCtx *svc.ServiceContext, r *http.Request) *DeleteUserLogic {
	return &DeleteUserLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
		r:      r,
	}
}

func (l *DeleteUserLogic) DeleteUser(req *types.DeleteUserReq) (resp *types.DeleteUserResp, err error) {
	// 调用 RPC 服务删除用户 - 使用熔断器
	_, err = breaker.DoWithBreakerResultAcceptable(l.ctx, "oauser-rpc", func() (*oauserclient.DeleteUserResp, error) {
		return l.svcCtx.OaUserRpc.DeleteUser(l.ctx, &oauserclient.DeleteUserReq{
			Phone: req.Phone,
		})
	}, breaker.IsAcceptableError)
	if err != nil {
		l.Logger.Errorf("RPC DeleteUser failed: %v", err)
		return nil, err
	}

	// 返回空结构体
	return &types.DeleteUserResp{}, nil
}
