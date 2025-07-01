package admin

import (
	"context"

	"api/internal/breaker"
	"api/internal/svc"
	"api/internal/types"
	"rpc/loanproduct"

	"github.com/zeromicro/go-zero/core/logx"
)

type CreateLoanProductLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 创建贷款产品
func NewCreateLoanProductLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CreateLoanProductLogic {
	return &CreateLoanProductLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *CreateLoanProductLogic) CreateLoanProduct(req *types.CreateLoanProductReq) (resp *types.CreateLoanProductResp, err error) {
	// 使用熔断器调用RPC服务
	rpcResp, err := breaker.DoWithBreakerResultAcceptable(l.ctx, "loanproduct-rpc", func() (*loanproduct.CreateLoanProductResp, error) {
		return l.svcCtx.LoanProductRpc.CreateLoanProduct(l.ctx, &loanproduct.CreateLoanProductReq{
			ProductCode:  req.ProductCode,
			Name:         req.Name,
			Type:         req.Type,
			MaxAmount:    req.MaxAmount,
			MinAmount:    req.MinAmount,
			MaxDuration:  req.MaxDuration,
			MinDuration:  req.MinDuration,
			InterestRate: req.InterestRate,
			Description:  req.Description,
		})
	}, breaker.IsAcceptableError)

	if err != nil {
		l.Errorf("调用RPC服务失败: %v", err)
		return nil, err
	}

	// 转换响应数据
	return &types.CreateLoanProductResp{
		Data: types.LoanProductInfo{
			Id:           rpcResp.Data.Id,
			ProductCode:  rpcResp.Data.ProductCode,
			Name:         rpcResp.Data.Name,
			Type:         rpcResp.Data.Type,
			MaxAmount:    rpcResp.Data.MaxAmount,
			MinAmount:    rpcResp.Data.MinAmount,
			MaxDuration:  rpcResp.Data.MaxDuration,
			MinDuration:  rpcResp.Data.MinDuration,
			InterestRate: rpcResp.Data.InterestRate,
			Description:  rpcResp.Data.Description,
			Status:       rpcResp.Data.Status,
			CreatedAt:    rpcResp.Data.CreatedAt,
			UpdatedAt:    rpcResp.Data.UpdatedAt,
		},
	}, nil
}
