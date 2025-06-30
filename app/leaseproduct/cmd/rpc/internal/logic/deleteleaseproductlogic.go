package logic

import (
	"context"

	"rpc/internal/svc"
	"rpc/leaseproduct"

	"github.com/zeromicro/go-zero/core/logx"
)

type DeleteLeaseProductLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewDeleteLeaseProductLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DeleteLeaseProductLogic {
	return &DeleteLeaseProductLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *DeleteLeaseProductLogic) DeleteLeaseProduct(in *leaseproduct.DeleteLeaseProductReq) (*leaseproduct.DeleteLeaseProductResp, error) {
	// 参数验证
	if in.ProductCode == "" {
		return &leaseproduct.DeleteLeaseProductResp{
			Code:    400,
			Message: "产品编码不能为空",
		}, nil
	}

	// 检查产品是否存在
	product, err := l.svcCtx.LeaseProductModel.FindOneByProductCode(l.ctx, in.ProductCode)
	if err != nil {
		l.Errorf("查询产品失败: %v", err)
		return &leaseproduct.DeleteLeaseProductResp{
			Code:    404,
			Message: "产品不存在",
		}, nil
	}

	// TODO: 检查是否有正在进行的租赁申请
	// 这里应该调用lease.rpc检查是否有未完成的租赁申请
	// 如果有正在进行的租赁，应该禁止删除

	// 删除产品
	err = l.svcCtx.LeaseProductModel.Delete(l.ctx, product.Id)
	if err != nil {
		l.Errorf("删除产品失败: %v", err)
		return &leaseproduct.DeleteLeaseProductResp{
			Code:    500,
			Message: "删除产品失败",
		}, nil
	}

	return &leaseproduct.DeleteLeaseProductResp{
		Code:    200,
		Message: "删除成功",
	}, nil
}
