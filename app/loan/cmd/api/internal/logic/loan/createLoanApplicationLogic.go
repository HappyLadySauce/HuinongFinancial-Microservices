package loan

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"

	"api/internal/breaker"
	"api/internal/svc"
	"api/internal/types"
	"rpc/loanclient"

	"github.com/zeromicro/go-zero/core/logx"
)

type CreateLoanApplicationLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewCreateLoanApplicationLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CreateLoanApplicationLogic {
	return &CreateLoanApplicationLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *CreateLoanApplicationLogic) CreateLoanApplication(req *types.CreateLoanApplicationReq) (resp *types.CreateLoanApplicationResp, err error) {
	// 从JWT上下文中获取用户ID
	userId, err := l.getUserIdFromJWT()
	if err != nil {
		logx.WithContext(l.ctx).Errorf("获取用户ID失败: %v", err)
		return nil, err
	}

	// 使用熔断器调用 Loan RPC 创建申请
	rpcResp, err := breaker.DoWithBreakerResultAcceptable(l.ctx, "loan-rpc", func() (*loanclient.CreateLoanApplicationResp, error) {
		return l.svcCtx.LoanRpc.CreateLoanApplication(l.ctx, &loanclient.CreateLoanApplicationReq{
			UserId:    userId,
			ProductId: req.ProductId,
			Name:      req.Name,
			Type:      req.Type,
			Amount:    req.Amount,
			Duration:  req.Duration,
			Purpose:   req.Purpose,
		})
	}, breaker.IsAcceptableError)

	if err != nil {
		logx.WithContext(l.ctx).Errorf("调用Loan RPC失败: %v", err)
		return nil, err
	}

	// 转换 RPC 响应为 API 响应
	return &types.CreateLoanApplicationResp{
		ApplicationId: rpcResp.ApplicationId,
	}, nil
}

// 从JWT中获取用户ID的辅助方法
func (l *CreateLoanApplicationLogic) getUserIdFromJWT() (int64, error) {
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
