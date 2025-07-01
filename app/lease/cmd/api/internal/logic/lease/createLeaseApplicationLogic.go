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
			ProductCode:     req.ProductCode,
			StartDate:       req.StartDate,
			EndDate:         req.EndDate,
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
	// 从JWT上下文中获取用户ID
	// 这里需要根据实际的JWT实现来获取用户ID
	// 示例实现：从上下文中获取JWT claims

	// 方式1：从HTTP请求上下文获取JWT claims
	// 假设JWT中包含user_id字段
	claims := l.ctx.Value("jwt_claims")
	if claims == nil {
		return 0, fmt.Errorf("未找到JWT认证信息")
	}

	// 将claims转换为map
	claimsMap, ok := claims.(map[string]interface{})
	if !ok {
		return 0, fmt.Errorf("JWT claims格式错误")
	}

	// 获取user_id
	userIdInterface, exists := claimsMap["user_id"]
	if !exists {
		return 0, fmt.Errorf("JWT中缺少user_id字段")
	}

	// 类型转换
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
