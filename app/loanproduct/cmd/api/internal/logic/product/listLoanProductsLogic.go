package product

import (
	"context"

	"api/internal/svc"
	"api/internal/types"
	"rpc/loanproduct"

	"github.com/zeromicro/go-zero/core/logx"
)

type ListLoanProductsLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 获取贷款产品列表
func NewListLoanProductsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ListLoanProductsLogic {
	return &ListLoanProductsLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ListLoanProductsLogic) ListLoanProducts(req *types.ListLoanProductsReq) (resp *types.ListLoanProductsResp, err error) {
	// 设置默认值
	if req.Page <= 0 {
		req.Page = 1
	}
	if req.Size <= 0 {
		req.Size = 10
	}

	// 调用RPC服务
	rpcResp, err := l.svcCtx.LoanProductRpc.ListLoanProducts(l.ctx, &loanproduct.ListLoanProductsReq{
		Page:    req.Page,
		Size:    req.Size,
		Type:    req.Type,
		Status:  req.Status,
		Keyword: req.Keyword,
	})
	if err != nil {
		l.Errorf("调用RPC服务失败: %v", err)
		return &types.ListLoanProductsResp{
			Code:    500,
			Message: "服务内部错误",
		}, nil
	}

	// 检查RPC响应
	if rpcResp.Code != 200 {
		return &types.ListLoanProductsResp{
			Code:    rpcResp.Code,
			Message: rpcResp.Message,
		}, nil
	}

	// 转换产品列表数据
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
		Code:    200,
		Message: "查询成功",
		List:    products,
		Total:   rpcResp.Total,
	}, nil
}
