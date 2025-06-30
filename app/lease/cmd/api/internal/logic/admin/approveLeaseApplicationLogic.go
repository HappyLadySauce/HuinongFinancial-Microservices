package admin

import (
	"context"
	"fmt"
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
	auditorId, err := l.getUserIdFromJWT()
	if err != nil {
		logx.WithContext(l.ctx).Errorf("获取审核员ID失败: %v", err)
		return nil, err
	}

	// 获取审核员姓名 (从JWT中获取或设置默认值)
	auditorName := l.getUserNameFromJWT()
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

// 从JWT中获取用户ID的辅助方法
func (l *ApproveLeaseApplicationLogic) getUserIdFromJWT() (int64, error) {
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

// 从JWT中获取用户名的辅助方法
func (l *ApproveLeaseApplicationLogic) getUserNameFromJWT() string {
	if nameVal := l.ctx.Value("username"); nameVal != nil {
		if name, ok := nameVal.(string); ok {
			return name
		}
	}
	if nameVal := l.ctx.Value("name"); nameVal != nil {
		if name, ok := nameVal.(string); ok {
			return name
		}
	}
	if phoneVal := l.ctx.Value("phone"); phoneVal != nil {
		if phone, ok := phoneVal.(string); ok {
			return phone // 如果没有用户名，使用手机号
		}
	}
	return ""
}
