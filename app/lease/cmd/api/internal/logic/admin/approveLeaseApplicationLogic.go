package admin

import (
	"context"
	"strconv"

	"api/internal/svc"
	"api/internal/types"
	"rpc/leaseclient"

	"github.com/zeromicro/go-zero/core/logx"
)

type ApproveLeaseApplicationLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewApproveLeaseApplicationLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ApproveLeaseApplicationLogic {
	return &ApproveLeaseApplicationLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ApproveLeaseApplicationLogic) ApproveLeaseApplication(req *types.ApproveLeaseApplicationReq) (resp *types.ApproveLeaseApplicationResp, err error) {
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

	// 调用 Lease RPC 审批申请
	_, err = l.svcCtx.LeaseRpc.ApproveLeaseApplication(l.ctx, &leaseclient.ApproveLeaseApplicationReq{
		ApplicationId:    req.ApplicationId,
		AuditorId:        auditorId,
		AuditorName:      auditorName,
		Action:           req.Action,
		Suggestions:      req.Suggestions,
		ApprovedDuration: req.ApprovedDuration,
		ApprovedAmount:   req.ApprovedAmount,
		ApprovedDeposit:  req.ApprovedDeposit,
	})
	if err != nil {
		logx.WithContext(l.ctx).Errorf("调用Lease RPC失败: %v", err)
		return nil, err
	}

	// 转换 RPC 响应为 API 响应
	return &types.ApproveLeaseApplicationResp{}, nil
}
