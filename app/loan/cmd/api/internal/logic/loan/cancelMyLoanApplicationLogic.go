package loan

import (
	"context"

	"api/internal/svc"
	"api/internal/types"
	"rpc/loanclient"

	"github.com/zeromicro/go-zero/core/logx"
)

type CancelMyLoanApplicationLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewCancelMyLoanApplicationLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CancelMyLoanApplicationLogic {
	return &CancelMyLoanApplicationLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *CancelMyLoanApplicationLogic) CancelMyLoanApplication(req *types.CancelLoanApplicationReq) (resp *types.CancelLoanApplicationResp, err error) {
	// 调用 Loan RPC 撤销申请
	rpcResp, err := l.svcCtx.LoanRpc.CancelLoanApplication(l.ctx, &loanclient.CancelLoanApplicationReq{
		ApplicationId: req.ApplicationId,
		Reason:        req.Reason,
	})
	if err != nil {
		logx.WithContext(l.ctx).Errorf("调用Loan RPC失败: %v", err)
		return &types.CancelLoanApplicationResp{
			Code:    500,
			Message: "服务内部错误",
		}, nil
	}

	// 转换 RPC 响应为 API 响应
	return &types.CancelLoanApplicationResp{
		Code:    rpcResp.Code,
		Message: rpcResp.Message,
	}, nil
}
