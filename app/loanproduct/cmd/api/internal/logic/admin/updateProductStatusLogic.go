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

type UpdateProductStatusLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 更新产品状态
func NewUpdateProductStatusLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateProductStatusLogic {
	return &UpdateProductStatusLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UpdateProductStatusLogic) UpdateProductStatus(idStr string, req *types.UpdateProductStatusReq) (resp *types.UpdateProductStatusResp, err error) {
	// 将字符串ID转换为int64
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		return nil, err
	}

	// 调用RPC服务 - 使用熔断器
	_, err = breaker.DoWithBreakerResultAcceptable(l.ctx, "loanproduct-rpc", func() (*loanproduct.UpdateProductStatusResp, error) {
		return l.svcCtx.LoanProductRpc.UpdateProductStatus(l.ctx, &loanproduct.UpdateProductStatusReq{
			Id:     id,
			Status: req.Status,
		})
	}, breaker.IsAcceptableError)
	if err != nil {
		l.Errorf("调用RPC服务失败: %v", err)
		return nil, err
	}

	return &types.UpdateProductStatusResp{}, nil
}
