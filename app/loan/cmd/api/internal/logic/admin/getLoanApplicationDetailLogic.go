package admin

import (
	"context"

	"api/internal/svc"
	"api/internal/types"
	"rpc/loanclient"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetLoanApplicationDetailLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetLoanApplicationDetailLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetLoanApplicationDetailLogic {
	return &GetLoanApplicationDetailLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetLoanApplicationDetailLogic) GetLoanApplicationDetail(applicationId string) (resp *types.GetLoanApplicationResp, err error) {
	// 调用 Loan RPC 获取申请详情
	rpcResp, err := l.svcCtx.LoanRpc.GetLoanApplication(l.ctx, &loanclient.GetLoanApplicationReq{
		ApplicationId: applicationId,
	})
	if err != nil {
		logx.WithContext(l.ctx).Errorf("调用Loan RPC失败: %v", err)
		return nil, err
	}

	// 转换申请信息
	return &types.GetLoanApplicationResp{
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
