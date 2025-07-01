package admin

import (
	"context"

	"api/internal/breaker"
	"api/internal/svc"
	"api/internal/types"
	"rpc/leaseproductservice"

	"github.com/zeromicro/go-zero/core/logx"
)

type ListAllLeaseProductsLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 获取所有租赁产品列表
func NewListAllLeaseProductsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ListAllLeaseProductsLogic {
	return &ListAllLeaseProductsLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ListAllLeaseProductsLogic) ListAllLeaseProducts(req *types.ListLeaseProductsReq) (resp *types.ListLeaseProductsResp, err error) {
	// 设置默认值
	if req.Page <= 0 {
		req.Page = 1
	}
	if req.Size <= 0 {
		req.Size = 10
	}

	// 调用RPC服务 - 使用熔断器
	rpcResp, err := breaker.DoWithBreakerResultAcceptable(l.ctx, "leaseproduct-rpc", func() (*leaseproductservice.ListLeaseProductsResp, error) {
		return l.svcCtx.LeaseProductRpc.ListLeaseProducts(l.ctx, &leaseproductservice.ListLeaseProductsReq{
			Page:   req.Page,
			Size:   req.Size,
			Type:   req.Type,
			Status: req.Status,
		})
	}, breaker.IsAcceptableError)
	if err != nil {
		l.Errorf("调用RPC服务失败: %v", err)
		return nil, err
	}

	// 转换产品列表数据
	var products []types.LeaseProductInfo
	for _, item := range rpcResp.List {
		products = append(products, types.LeaseProductInfo{
			Id:             item.Id,
			ProductCode:    item.ProductCode,
			Name:           item.Name,
			Type:           item.Type,
			Machinery:      item.Machinery,
			Brand:          item.Brand,
			Model:          item.Model,
			DailyRate:      item.DailyRate,
			Deposit:        item.Deposit,
			MaxDuration:    item.MaxDuration,
			MinDuration:    item.MinDuration,
			Description:    item.Description,
			InventoryCount: item.InventoryCount,
			AvailableCount: item.AvailableCount,
			Status:         item.Status,
			CreatedAt:      item.CreatedAt,
			UpdatedAt:      item.UpdatedAt,
		})
	}

	return &types.ListLeaseProductsResp{
		List:  products,
		Total: rpcResp.Total,
	}, nil
}
