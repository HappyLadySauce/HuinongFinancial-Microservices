package logic

import (
	"context"
	"fmt"

	"rpc/internal/svc"
	"rpc/loanproduct"

	"github.com/zeromicro/go-zero/core/logx"
)

type DeleteLoanProductLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewDeleteLoanProductLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DeleteLoanProductLogic {
	return &DeleteLoanProductLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *DeleteLoanProductLogic) DeleteLoanProduct(in *loanproduct.DeleteLoanProductReq) (*loanproduct.DeleteLoanProductResp, error) {
	// 参数验证
	if in.Id <= 0 {
		return nil, fmt.Errorf("产品ID不能为空")
	}

	// 检查产品是否存在
	product, err := l.svcCtx.LoanProductModel.FindOne(l.ctx, uint64(in.Id))
	if err != nil {
		l.Errorf("查询产品失败: %v", err)
		return nil, fmt.Errorf("产品不存在")
	}

	// TODO: 检查是否有正在进行的贷款申请
	// 这里应该调用loan.rpc检查是否有未完成的贷款申请
	// 如果有正在进行的贷款，应该禁止删除

	// 删除产品
	err = l.svcCtx.LoanProductModel.Delete(l.ctx, product.Id)
	if err != nil {
		l.Errorf("删除产品失败: %v", err)
		return nil, fmt.Errorf("删除产品失败")
	}

	return &loanproduct.DeleteLoanProductResp{}, nil
}
