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
