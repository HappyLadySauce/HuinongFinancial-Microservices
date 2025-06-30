package product

import (
	"context"
	"strconv"

	"api/internal/svc"
	"api/internal/types"
	"rpc/loanproduct"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetLoanProductLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 获取贷款产品详情
func NewGetLoanProductLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetLoanProductLogic {
	return &GetLoanProductLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetLoanProductLogic) GetLoanProduct(idStr string) (resp *types.GetLoanProductResp, err error) {
	// 将字符串ID转换为int64
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		return nil, err
	}

	// 调用RPC服务
	rpcResp, err := l.svcCtx.LoanProductRpc.GetLoanProduct(l.ctx, &loanproduct.GetLoanProductReq{
		Id: id,
	})
	if err != nil {
		l.Errorf("调用RPC服务失败: %v", err)
		return nil, err
	}

	// 转换响应数据
	return &types.GetLoanProductResp{
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
