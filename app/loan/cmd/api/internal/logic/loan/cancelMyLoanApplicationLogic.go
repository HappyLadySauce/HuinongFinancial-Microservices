package loan

import (
	"context"

	"api/internal/breaker"
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
	// 调用 Loan RPC 取消申请 - 使用熔断器
	_, err = breaker.DoWithBreakerResultAcceptable(l.ctx, "loan-rpc", func() (*loanclient.CancelLoanApplicationResp, error) {
		return l.svcCtx.LoanRpc.CancelLoanApplication(l.ctx, &loanclient.CancelLoanApplicationReq{
			ApplicationId: req.ApplicationId,
			Reason:        req.Reason,
		})
	}, breaker.IsAcceptableError)
	if err != nil {
		logx.WithContext(l.ctx).Errorf("调用Loan RPC失败: %v", err)
		return nil, err
	}

	// 返回空结构体
	return &types.CancelLoanApplicationResp{}, nil
}
