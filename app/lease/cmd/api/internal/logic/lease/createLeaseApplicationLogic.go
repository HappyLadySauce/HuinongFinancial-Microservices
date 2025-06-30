package lease

import (
	"context"
	"fmt"
	"strconv"

	"api/internal/svc"
	"api/internal/types"
	"rpc/leaseclient"

	"github.com/zeromicro/go-zero/core/logx"
)

type CreateLeaseApplicationLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewCreateLeaseApplicationLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CreateLeaseApplicationLogic {
	return &CreateLeaseApplicationLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *CreateLeaseApplicationLogic) CreateLeaseApplication(req *types.CreateLeaseApplicationReq) (resp *types.CreateLeaseApplicationResp, err error) {
	// 获取当前用户ID (从JWT中获取)
	// 在go-zero中，JWT claims通过httpx.ParseJsonBody解析后存储在特定的context key中
	// 根据JWT配置，user_id字段应该从JWT的标准claims中获取
	userId, err := l.getUserIdFromJWT()
	if err != nil {
		logx.WithContext(l.ctx).Errorf("获取用户ID失败: %v", err)
		return nil, err
	}

	// 调用 Lease RPC 创建申请
	rpcResp, err := l.svcCtx.LeaseRpc.CreateLeaseApplication(l.ctx, &leaseclient.CreateLeaseApplicationReq{
		UserId:          userId,
		ProductId:       req.ProductId,
		ProductCode:     req.ProductCode,
		Name:            req.Name,
		Type:            req.Type,
		Machinery:       req.Machinery,
		StartDate:       req.StartDate,
		EndDate:         req.EndDate,
		Duration:        req.Duration,
		DailyRate:       req.DailyRate,
		TotalAmount:     req.TotalAmount,
		Deposit:         req.Deposit,
		DeliveryAddress: req.DeliveryAddress,
		ContactPhone:    req.ContactPhone,
		Purpose:         req.Purpose,
	})
	if err != nil {
		logx.WithContext(l.ctx).Errorf("调用Lease RPC失败: %v", err)
		return nil, err
	}

	// 转换 RPC 响应为 API 响应
	return &types.CreateLeaseApplicationResp{
		ApplicationId: rpcResp.ApplicationId,
	}, nil
}

// 从JWT中获取用户ID的辅助方法
func (l *CreateLeaseApplicationLogic) getUserIdFromJWT() (int64, error) {
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
