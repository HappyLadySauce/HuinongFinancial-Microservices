package loan

import (
	"context"

	"api/internal/svc"
	"api/internal/types"
	"rpc/loanclient"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetMyLoanApplicationLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetMyLoanApplicationLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetMyLoanApplicationLogic {
	return &GetMyLoanApplicationLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetMyLoanApplicationLogic) GetMyLoanApplication(applicationId string) (resp *types.GetLoanApplicationResp, err error) {
	// 调用 Loan RPC 获取申请详情
	rpcResp, err := l.svcCtx.LoanRpc.GetLoanApplication(l.ctx, &loanclient.GetLoanApplicationReq{
		ApplicationId: applicationId,
	})
	if err != nil {
		logx.WithContext(l.ctx).Errorf("调用Loan RPC失败: %v", err)
		return &types.GetLoanApplicationResp{
			Code:    500,
			Message: "服务内部错误",
		}, nil
	}

	// 转换 RPC 响应为 API 响应
	resp = &types.GetLoanApplicationResp{
		Code:    rpcResp.Code,
		Message: rpcResp.Message,
	}

	// 注意：RPC响应中的字段是 ApplicationInfo 而不是 Data
	if rpcResp.ApplicationInfo != nil {
		resp.ApplicationInfo = types.LoanApplicationInfo{
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
		}
	}

	return resp, nil
}
