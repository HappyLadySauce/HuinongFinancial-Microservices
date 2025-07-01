package admin

import (
	"context"

	"api/internal/breaker"
	"api/internal/svc"
	"api/internal/types"
	"rpc/loanclient"

	"github.com/zeromicro/go-zero/core/logx"
)

type ListAllLoanApplicationsLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewListAllLoanApplicationsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ListAllLoanApplicationsLogic {
	return &ListAllLoanApplicationsLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ListAllLoanApplicationsLogic) ListAllLoanApplications(req *types.ListLoanApplicationsReq) (resp *types.ListLoanApplicationsResp, err error) {
	// 调用 Loan RPC 获取申请列表 (管理员可查看所有申请，不过滤UserId) - 使用熔断器
	rpcResp, err := breaker.DoWithBreakerResultAcceptable(l.ctx, "loan-rpc", func() (*loanclient.ListLoanApplicationsResp, error) {
		return l.svcCtx.LoanRpc.ListLoanApplications(l.ctx, &loanclient.ListLoanApplicationsReq{
			Page:   req.Page,
			Size:   req.Size,
			UserId: req.UserId, // 管理员可以按用户ID过滤，也可以不过滤查看所有
			Status: req.Status,
		})
	}, breaker.IsAcceptableError)
	if err != nil {
		logx.WithContext(l.ctx).Errorf("调用Loan RPC失败: %v", err)
		return nil, err
	}

	// 转换申请列表
	var applications []types.LoanApplicationInfo
	for _, item := range rpcResp.List {
		applications = append(applications, types.LoanApplicationInfo{
			Id:            item.Id,
			ApplicationId: item.ApplicationId,
			UserId:        item.UserId,
			ApplicantName: item.ApplicantName,
			ProductId:     item.ProductId,
			Name:          item.Name,
			Type:          item.Type,
			Amount:        item.Amount,
			Duration:      item.Duration,
			Purpose:       item.Purpose,
			Status:        item.Status,
			CreatedAt:     item.CreatedAt,
			UpdatedAt:     item.UpdatedAt,
		})
	}

	return &types.ListLoanApplicationsResp{
		List:  applications,
		Total: rpcResp.Total,
	}, nil
}
