package admin

import (
	"context"

	"api/internal/breaker"
	"api/internal/svc"
	"api/internal/types"
	"rpc/loanproduct"

	"github.com/zeromicro/go-zero/core/logx"
)

type ListAllLoanProductsLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 获取所有贷款产品列表
func NewListAllLoanProductsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ListAllLoanProductsLogic {
	return &ListAllLoanProductsLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ListAllLoanProductsLogic) ListAllLoanProducts(req *types.ListLoanProductsReq) (resp *types.ListLoanProductsResp, err error) {
	// 设置默认值
	if req.Page <= 0 {
		req.Page = 1
	}
	if req.Size <= 0 {
		req.Size = 10
	}

	// 调用RPC服务 - 使用熔断器
	rpcResp, err := breaker.DoWithBreakerResultAcceptable(l.ctx, "loanproduct-rpc", func() (*loanproduct.ListLoanProductsResp, error) {
		return l.svcCtx.LoanProductRpc.ListLoanProducts(l.ctx, &loanproduct.ListLoanProductsReq{
			Page:   req.Page,
			Size:   req.Size,
			Status: req.Status,
		})
	}, breaker.IsAcceptableError)
	if err != nil {
		l.Errorf("调用RPC服务失败: %v", err)
		return nil, err
	}

	// 转换响应数据
	var products []types.LoanProductInfo
	for _, item := range rpcResp.List {
		products = append(products, types.LoanProductInfo{
			Id:           item.Id,
			ProductCode:  item.ProductCode,
			Name:         item.Name,
			Type:         item.Type,
			MaxAmount:    item.MaxAmount,
			MinAmount:    item.MinAmount,
			MaxDuration:  item.MaxDuration,
			MinDuration:  item.MinDuration,
			InterestRate: item.InterestRate,
			Description:  item.Description,
			Status:       item.Status,
			CreatedAt:    item.CreatedAt,
			UpdatedAt:    item.UpdatedAt,
		})
	}

	return &types.ListLoanProductsResp{
		List:  products,
		Total: rpcResp.Total,
	}, nil
}
