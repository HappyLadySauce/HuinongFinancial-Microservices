package admin

import (
	"context"

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
	// 调用 Loan RPC 获取申请列表 (管理员可查看所有申请，不过滤UserId)
	rpcResp, err := l.svcCtx.LoanRpc.ListLoanApplications(l.ctx, &loanclient.ListLoanApplicationsReq{
		Page:   req.Page,
		Size:   req.Size,
		UserId: req.UserId, // 管理员可以按用户ID过滤，也可以不过滤查看所有
		Status: req.Status,
	})
	if err != nil {
		logx.WithContext(l.ctx).Errorf("调用Loan RPC失败: %v", err)
		return &types.ListLoanApplicationsResp{
			Code:    500,
			Message: "服务内部错误",
		}, nil
	}

	// 转换 RPC 响应为 API 响应
	resp = &types.ListLoanApplicationsResp{
		Code:    rpcResp.Code,
		Message: rpcResp.Message,
		Total:   rpcResp.Total,
	}

	// 转换申请列表
	if len(rpcResp.List) > 0 {
		resp.List = make([]types.LoanApplicationInfo, len(rpcResp.List))
		for i, item := range rpcResp.List {
			resp.List[i] = types.LoanApplicationInfo{
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
			}
		}
	} else {
		resp.List = make([]types.LoanApplicationInfo, 0)
	}

	return resp, nil
}
