package lease

import (
	"context"
	"fmt"
	"strconv"

	"api/internal/breaker"
	"api/internal/svc"
	"api/internal/types"
	"rpc/leaseclient"

	"encoding/json"

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

	// 调用 Lease RPC 创建申请 - 使用熔断器
	rpcResp, err := breaker.DoWithBreakerResultAcceptable(l.ctx, "lease-rpc", func() (*leaseclient.CreateLeaseApplicationResp, error) {
		return l.svcCtx.LeaseRpc.CreateLeaseApplication(l.ctx, &leaseclient.CreateLeaseApplicationReq{
			UserId:          userId,
			ProductId:       req.ProductId, // 添加缺失的ProductId
			ProductCode:     req.ProductCode,
			Name:            req.Name,      // 添加缺失的Name
			Type:            req.Type,      // 添加缺失的Type
			Machinery:       req.Machinery, // 添加缺失的Machinery
			StartDate:       req.StartDate,
			EndDate:         req.EndDate,
			Duration:        req.Duration,    // 添加缺失的Duration
			DailyRate:       req.DailyRate,   // 添加缺失的DailyRate
			TotalAmount:     req.TotalAmount, // 添加缺失的TotalAmount
			Deposit:         req.Deposit,     // 添加缺失的Deposit
			DeliveryAddress: req.DeliveryAddress,
			ContactPhone:    req.ContactPhone,
			Purpose:         req.Purpose,
		})
	}, breaker.IsAcceptableError)
	if err != nil {
		logx.WithContext(l.ctx).Errorf("调用Lease RPC失败: %v", err)
		return nil, err
	}

	// 转换响应信息 - CreateLeaseApplicationResp只包含ApplicationId
	return &types.CreateLeaseApplicationResp{
		ApplicationId: rpcResp.ApplicationId,
	}, nil
}

// 从JWT中获取用户ID的辅助方法
func (l *CreateLeaseApplicationLogic) getUserIdFromJWT() (int64, error) {
	// 在go-zero中，JWT claims直接存储在context中
	// 尝试获取JWT claims
	if claims := l.ctx.Value("user_id"); claims != nil {
		switch v := claims.(type) {
		case int64:
			return v, nil
		case float64:
			return int64(v), nil
		case string:
			return strconv.ParseInt(v, 10, 64)
		case json.Number:
			return v.Int64()
		default:
			return 0, fmt.Errorf("user_id类型错误: %T", v)
		}
	}

	// 如果直接获取user_id失败，尝试获取完整的JWT claims
	if rawClaims := l.ctx.Value("claims"); rawClaims != nil {
		if claimsMap, ok := rawClaims.(map[string]interface{}); ok {
			if userIdInterface, exists := claimsMap["user_id"]; exists {
				switch v := userIdInterface.(type) {
				case float64:
					return int64(v), nil
				case string:
					return strconv.ParseInt(v, 10, 64)
				case json.Number:
					return v.Int64()
				default:
					return 0, fmt.Errorf("user_id类型错误: %T", v)
				}
			}
		}
	}

	return 0, fmt.Errorf("未找到JWT认证信息")
}
