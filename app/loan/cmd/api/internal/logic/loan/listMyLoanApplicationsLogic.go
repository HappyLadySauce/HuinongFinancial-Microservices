package loan

import (
	"context"
	"fmt"
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
	userId, err := l.getUserIdFromJWT()
	if err != nil {
		logx.WithContext(l.ctx).Errorf("获取用户ID失败: %v", err)
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
	for _, app := range rpcResp.List {
		applications = append(applications, types.LoanApplicationInfo{
			Id:        app.Id,
			UserId:    app.UserId,
			ProductId: app.ProductId,
			Name:      app.Name,
			Type:      app.Type,
			Amount:    app.Amount,
			Duration:  app.Duration,
			Purpose:   app.Purpose,
			Status:    app.Status,
			CreatedAt: app.CreatedAt,
			UpdatedAt: app.UpdatedAt,
		})
	}

	return &types.ListLoanApplicationsResp{
		List:  applications,
		Total: rpcResp.Total,
	}, nil
}

// 从JWT中获取用户ID的辅助方法
func (l *ListMyLoanApplicationsLogic) getUserIdFromJWT() (int64, error) {
	// 方法1: 尝试从context的标准JWT字段获取
	if userIdVal := l.ctx.Value("user_id"); userIdVal != nil {
		if userId, ok := userIdVal.(float64); ok {
			return int64(userId), nil
		}
		if userId, ok := userIdVal.(int64); ok {
			return userId, nil
		}
		if userIdStr, ok := userIdVal.(string); ok {
			return strconv.ParseInt(userIdStr, 10, 64)
		}
	}

	// 方法2: 尝试从context的其他可能字段获取
	if userIdVal := l.ctx.Value("userId"); userIdVal != nil {
		if userId, ok := userIdVal.(float64); ok {
			return int64(userId), nil
		}
		if userId, ok := userIdVal.(int64); ok {
			return userId, nil
		}
		if userIdStr, ok := userIdVal.(string); ok {
			return strconv.ParseInt(userIdStr, 10, 64)
		}
	}

	// 方法3: 尝试从JWT标准字段获取 (sub字段通常包含用户ID)
	if subVal := l.ctx.Value("sub"); subVal != nil {
		if subStr, ok := subVal.(string); ok {
			return strconv.ParseInt(subStr, 10, 64)
		}
	}

	return 0, fmt.Errorf("无法从JWT中获取用户ID")
}
