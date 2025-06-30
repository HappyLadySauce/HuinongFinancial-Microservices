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
	// 获取当前审核员ID (从JWT中获取)
	auditorIdStr := l.ctx.Value("userId").(string)
	auditorId, err := strconv.ParseInt(auditorIdStr, 10, 64)
	if err != nil {
		logx.WithContext(l.ctx).Errorf("解析审核员ID失败: %v", err)
		return nil, err
	}

	// 获取审核员姓名 (从JWT中获取或调用用户服务)
	auditorName := l.ctx.Value("username").(string)
	if auditorName == "" {
		auditorName = "系统管理员"
	}

	// 调用 Loan RPC 审批申请
	_, err = l.svcCtx.LoanRpc.ApproveLoanApplication(l.ctx, &loanclient.ApproveLoanApplicationReq{
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
		return nil, err
	}

	// 转换 RPC 响应为 API 响应
	return &types.ApproveLoanApplicationResp{}, nil
}
