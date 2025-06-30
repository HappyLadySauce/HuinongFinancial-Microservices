package loan

import (
	"context"

	"api/internal/svc"
	"api/internal/types"
	"rpc/loanclient"

	"github.com/zeromicro/go-zero/core/logx"
)

type UpdateMyLoanApplicationLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewUpdateMyLoanApplicationLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateMyLoanApplicationLogic {
	return &UpdateMyLoanApplicationLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UpdateMyLoanApplicationLogic) UpdateMyLoanApplication(req *types.UpdateLoanApplicationReq) (resp *types.UpdateLoanApplicationResp, err error) {
	// 调用 Loan RPC 更新申请
	rpcResp, err := l.svcCtx.LoanRpc.UpdateLoanApplication(l.ctx, &loanclient.UpdateLoanApplicationReq{
		ApplicationId: req.ApplicationId,
		Amount:        req.Amount,
		Duration:      req.Duration,
		Purpose:       req.Purpose,
	})
	if err != nil {
		logx.WithContext(l.ctx).Errorf("调用Loan RPC失败: %v", err)
		return nil, err
	}

	// 转换申请信息
	return &types.UpdateLoanApplicationResp{
		ApplicationInfo: types.LoanApplicationInfo{
			Id:            rpcResp.ApplicationInfo.Id,
			ApplicationId: rpcResp.ApplicationInfo.ApplicationId,
			UserId:        rpcResp.ApplicationInfo.UserId,
			ApplicantName: rpcResp.ApplicationInfo.ApplicantName,
			ProductId:     rpcResp.ApplicationInfo.ProductId,
			Name:          rpcResp.ApplicationInfo.Name,
			Type:          rpcResp.ApplicationInfo.Type,
			Amount:        rpcResp.ApplicationInfo.Amount,
			Duration:      rpcResp.ApplicationInfo.Duration,
			Purpose:       rpcResp.ApplicationInfo.Purpose,
			Status:        rpcResp.ApplicationInfo.Status,
			CreatedAt:     rpcResp.ApplicationInfo.CreatedAt,
			UpdatedAt:     rpcResp.ApplicationInfo.UpdatedAt,
		},
	}, nil
}
