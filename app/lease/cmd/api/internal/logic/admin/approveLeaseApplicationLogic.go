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
	// 获取当前管理员用户ID (从JWT中获取)
	userIdStr := l.ctx.Value("userId").(string)
	auditorId, err := strconv.ParseInt(userIdStr, 10, 64)
	if err != nil {
		logx.WithContext(l.ctx).Errorf("解析审核员ID失败: %v", err)
		return &types.ApproveLeaseApplicationResp{
			Code:    400,
			Message: "审核员ID无效",
		}, nil
	}

	// 获取审核员姓名 (简化实现，实际应该从用户服务获取)
	auditorName := "管理员" // TODO: 从用户服务获取真实姓名

	// 调用 Lease RPC 审批申请
	rpcResp, err := l.svcCtx.LeaseRpc.ApproveLeaseApplication(l.ctx, &leaseclient.ApproveLeaseApplicationReq{
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
		return &types.ApproveLeaseApplicationResp{
			Code:    500,
			Message: "服务内部错误",
		}, nil
	}

	// 转换 RPC 响应为 API 响应
	return &types.ApproveLeaseApplicationResp{
		Code:    rpcResp.Code,
		Message: rpcResp.Message,
	}, nil
}
