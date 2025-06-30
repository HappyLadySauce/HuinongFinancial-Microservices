package admin

import (
	"context"

	"api/internal/svc"
	"api/internal/types"
	"rpc/leaseproductservice"

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
	// 调用 LeaseProduct RPC 删除产品
	_, err = l.svcCtx.LeaseProductRpc.DeleteLeaseProduct(l.ctx, &leaseproductservice.DeleteLeaseProductReq{
		ProductCode: req.ProductCode,
	})
	if err != nil {
		logx.WithContext(l.ctx).Errorf("调用LeaseProduct RPC失败: %v", err)
		return nil, err
	}

	// 返回响应
	return &types.DeleteLeaseProductResp{}, nil
}
