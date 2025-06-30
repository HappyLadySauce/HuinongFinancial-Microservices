package admin

import (
	"context"
	"strconv"

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
		return &types.DeleteLoanProductResp{
			Code:    400,
			Message: "无效的产品ID",
		}, nil
	}

	// 调用RPC服务
	rpcResp, err := l.svcCtx.LoanProductRpc.DeleteLoanProduct(l.ctx, &loanproduct.DeleteLoanProductReq{
		Id: id,
	})
	if err != nil {
		l.Errorf("调用RPC服务失败: %v", err)
		return &types.DeleteLoanProductResp{
			Code:    500,
			Message: "服务内部错误",
		}, nil
	}

	// 返回响应
	return &types.DeleteLoanProductResp{
		Code:    rpcResp.Code,
		Message: rpcResp.Message,
	}, nil
}
