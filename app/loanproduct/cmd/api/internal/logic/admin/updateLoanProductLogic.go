package admin

import (
	"context"
	"strconv"

	"api/internal/svc"
	"api/internal/types"
	"rpc/loanproduct"

	"github.com/zeromicro/go-zero/core/logx"
)

type UpdateLoanProductLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 更新贷款产品
func NewUpdateLoanProductLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateLoanProductLogic {
	return &UpdateLoanProductLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UpdateLoanProductLogic) UpdateLoanProduct(req *types.UpdateLoanProductReq) (resp *types.UpdateLoanProductResp, err error) {
	// 将字符串ID转换为int64
	id, err := strconv.ParseInt(req.Id, 10, 64)
	if err != nil {
		return &types.UpdateLoanProductResp{
			Code:    400,
			Message: "无效的产品ID",
		}, nil
	}

	// 调用RPC服务更新产品基本信息
	rpcResp, err := l.svcCtx.LoanProductRpc.UpdateLoanProduct(l.ctx, &loanproduct.UpdateLoanProductReq{
		Id:           id,
		Name:         req.Name,
		Type:         req.Type,
		MaxAmount:    req.MaxAmount,
		MinAmount:    req.MinAmount,
		MaxDuration:  req.MaxDuration,
		MinDuration:  req.MinDuration,
		InterestRate: req.InterestRate,
		Description:  req.Description,
	})
	if err != nil {
		l.Errorf("调用RPC服务失败: %v", err)
		return &types.UpdateLoanProductResp{
			Code:    500,
			Message: "服务内部错误",
		}, nil
	}

	// 检查RPC响应
	if rpcResp.Code != 200 {
		return &types.UpdateLoanProductResp{
			Code:    rpcResp.Code,
			Message: rpcResp.Message,
		}, nil
	}

	// 如果需要更新状态，调用UpdateProductStatus
	if req.Status > 0 {
		statusResp, err := l.svcCtx.LoanProductRpc.UpdateProductStatus(l.ctx, &loanproduct.UpdateProductStatusReq{
			Id:     id,
			Status: req.Status,
		})
		if err != nil {
			l.Errorf("更新产品状态失败: %v", err)
			return &types.UpdateLoanProductResp{
				Code:    500,
				Message: "更新产品状态失败",
			}, nil
		}
		if statusResp.Code != 200 {
			return &types.UpdateLoanProductResp{
				Code:    statusResp.Code,
				Message: statusResp.Message,
			}, nil
		}
	}

	// 转换响应数据
	return &types.UpdateLoanProductResp{
		Code:    200,
		Message: "更新成功",
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
