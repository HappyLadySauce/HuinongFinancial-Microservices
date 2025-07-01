package admin

import (
	"context"
	"strconv"

	"api/internal/breaker"
	"api/internal/svc"
	"api/internal/types"
	"rpc/loanproduct"

	"github.com/zeromicro/go-zero/core/logx"
)

type DeleteLoanProductLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 删除贷款产品
func NewDeleteLoanProductLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DeleteLoanProductLogic {
	return &DeleteLoanProductLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *DeleteLoanProductLogic) DeleteLoanProduct(req *types.DeleteLoanProductReq) (resp *types.DeleteLoanProductResp, err error) {
	// 将字符串ID转换为int64
	id, err := strconv.ParseInt(req.Id, 10, 64)
	if err != nil {
		return nil, err
	}

	// 调用RPC服务 - 使用熔断器
	_, err = breaker.DoWithBreakerResultAcceptable(l.ctx, "loanproduct-rpc", func() (*loanproduct.DeleteLoanProductResp, error) {
		return l.svcCtx.LoanProductRpc.DeleteLoanProduct(l.ctx, &loanproduct.DeleteLoanProductReq{
			Id: id,
		})
	}, breaker.IsAcceptableError)
	if err != nil {
		l.Errorf("调用RPC服务失败: %v", err)
		return nil, err
	}

	// 返回响应
	return &types.DeleteLoanProductResp{}, nil
}
