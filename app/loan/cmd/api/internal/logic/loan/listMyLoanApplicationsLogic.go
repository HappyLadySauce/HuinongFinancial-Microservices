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
		return &types.ListLoanApplicationsResp{
			Code:    400,
			Message: "用户ID无效",
		}, nil
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
