package admin

import (
	"context"

	"api/internal/svc"
	"api/internal/types"
	"rpc/leaseproduct"

	"github.com/zeromicro/go-zero/core/logx"
)

type DeleteLeaseProductLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 删除租赁产品
func NewDeleteLeaseProductLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DeleteLeaseProductLogic {
	return &DeleteLeaseProductLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *DeleteLeaseProductLogic) DeleteLeaseProduct(req *types.DeleteLeaseProductReq) (resp *types.DeleteLeaseProductResp, err error) {
	// 调用RPC服务
	rpcResp, err := l.svcCtx.LeaseProductRpc.DeleteLeaseProduct(l.ctx, &leaseproduct.DeleteLeaseProductReq{
		ProductCode: req.ProductCode,
	})
	if err != nil {
		l.Errorf("调用RPC服务失败: %v", err)
		return &types.DeleteLeaseProductResp{
			Code:    500,
			Message: "服务内部错误",
		}, nil
	}

	// 返回响应
	return &types.DeleteLeaseProductResp{
		Code:    rpcResp.Code,
		Message: rpcResp.Message,
	}, nil
}
