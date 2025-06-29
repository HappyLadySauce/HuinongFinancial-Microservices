package info

import (
	"context"

	"api/internal/svc"
	"api/internal/types"
	"rpc/oauserclient"

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

func (l *DeleteUserLogic) DeleteUser(req *types.DeleteUserReq) (resp *types.DeleteUserResp, err error) {
	// 调用 RPC 服务删除用户
	deleteResp, err := l.svcCtx.OaUserRpc.DeleteUser(l.ctx, &oauserclient.DeleteUserReq{
		Phone: req.Phone,
	})
	if err != nil {
		l.Logger.Errorf("RPC DeleteUser failed: %v", err)
		return &types.DeleteUserResp{
			Code:    500,
			Message: "服务器内部错误",
		}, nil
	}

	// 转换响应格式
	return &types.DeleteUserResp{
		Code:    deleteResp.Code,
		Message: deleteResp.Message,
	}, nil
}
