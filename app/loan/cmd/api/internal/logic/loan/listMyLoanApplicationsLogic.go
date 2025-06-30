package loan

import (
	"context"
	"strconv"

	"api/internal/svc"
	"api/internal/types"
	"rpc/loanclient"

	"github.com/zeromicro/go-zero/core/logx"
)

type ListMyLoanApplicationsLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewListMyLoanApplicationsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ListMyLoanApplicationsLogic {
	return &ListMyLoanApplicationsLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ListMyLoanApplicationsLogic) ListMyLoanApplications(req *types.ListLoanApplicationsReq) (resp *types.ListLoanApplicationsResp, err error) {
	// 获取当前用户ID (从JWT中获取)
	userIdStr := l.ctx.Value("userId").(string)
	userId, err := strconv.ParseInt(userIdStr, 10, 64)
	if err != nil {
		logx.WithContext(l.ctx).Errorf("解析用户ID失败: %v", err)
		return nil, err
	}

	// 调用 Loan RPC 获取申请列表
	rpcResp, err := l.svcCtx.LoanRpc.ListLoanApplications(l.ctx, &loanclient.ListLoanApplicationsReq{
		Page:   req.Page,
		Size:   req.Size,
		UserId: userId,
		Status: req.Status,
	})
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
