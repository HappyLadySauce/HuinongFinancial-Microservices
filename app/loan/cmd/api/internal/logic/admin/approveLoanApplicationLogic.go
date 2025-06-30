package admin

import (
	"context"
	"strconv"

	"api/internal/svc"
	"api/internal/types"
	"rpc/loanclient"

	"github.com/zeromicro/go-zero/core/logx"
)

type ApproveLoanApplicationLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewApproveLoanApplicationLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ApproveLoanApplicationLogic {
	return &ApproveLoanApplicationLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ApproveLoanApplicationLogic) ApproveLoanApplication(req *types.ApproveLoanApplicationReq) (resp *types.ApproveLoanApplicationResp, err error) {
	// 获取当前管理员用户ID (从JWT中获取)
	userIdStr := l.ctx.Value("userId").(string)
	auditorId, err := strconv.ParseInt(userIdStr, 10, 64)
	if err != nil {
		logx.WithContext(l.ctx).Errorf("解析审核员ID失败: %v", err)
		return &types.ApproveLoanApplicationResp{
			Code:    400,
			Message: "审核员ID无效",
		}, nil
	}

	// 获取审核员姓名 (简化实现，实际应该从用户服务获取)
	auditorName := "管理员" // TODO: 从用户服务获取真实姓名

	// 调用 Loan RPC 审批申请
	rpcResp, err := l.svcCtx.LoanRpc.ApproveLoanApplication(l.ctx, &loanclient.ApproveLoanApplicationReq{
		ApplicationId:    req.ApplicationId,
		AuditorId:        auditorId,
		AuditorName:      auditorName,
		Action:           req.Action,
		Suggestions:      req.Suggestions,
		ApprovedAmount:   req.ApprovedAmount,
		ApprovedDuration: req.ApprovedDuration,
		InterestRate:     req.InterestRate,
	})
	if err != nil {
		logx.WithContext(l.ctx).Errorf("调用Loan RPC失败: %v", err)
		return &types.ApproveLoanApplicationResp{
			Code:    500,
			Message: "服务内部错误",
		}, nil
	}

	// 转换 RPC 响应为 API 响应
	return &types.ApproveLoanApplicationResp{
		Code:    rpcResp.Code,
		Message: rpcResp.Message,
	}, nil
}
