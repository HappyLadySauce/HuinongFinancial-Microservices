package admin

import (
	"context"
	"strconv"

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

func (l *UpdateProductStatusLogic) UpdateProductStatus(req *types.UpdateProductStatusReq) (resp *types.UpdateProductStatusResp, err error) {
	// 将字符串ID转换为int64
	id, err := strconv.ParseInt(req.Id, 10, 64)
	if err != nil {
		return &types.UpdateProductStatusResp{
			Code:    400,
			Message: "无效的产品ID",
		}, nil
	}

	// 调用RPC服务
	rpcResp, err := l.svcCtx.LoanProductRpc.UpdateProductStatus(l.ctx, &loanproduct.UpdateProductStatusReq{
		Id:     id,
		Status: req.Status,
	})
	if err != nil {
		l.Errorf("调用RPC服务失败: %v", err)
		return &types.UpdateProductStatusResp{
			Code:    500,
			Message: "服务内部错误",
		}, nil
	}

	// 返回响应
	return &types.UpdateProductStatusResp{
		Code:    rpcResp.Code,
		Message: rpcResp.Message,
	}, nil
}
