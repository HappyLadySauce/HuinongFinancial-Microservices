package lease

import (
	"context"
	"fmt"
	"strconv"

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
	// 方法1: go-zero标准方式 - 处理json.Number类型
	if userIdVal := l.ctx.Value("user_id"); userIdVal != nil {
		// go-zero将JWT中的数字转换为json.Number类型
		if jsonUid, ok := userIdVal.(json.Number); ok {
			if int64Uid, err := jsonUid.Int64(); err == nil {
				return int64Uid, nil
			} else {
				logx.WithContext(l.ctx).Errorf("JWT user_id转换失败: %v", err)
			}
		}
		// 备用：尝试其他类型
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
		if jsonUid, ok := userIdVal.(json.Number); ok {
			if int64Uid, err := jsonUid.Int64(); err == nil {
				return int64Uid, nil
			}
		}
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
		if jsonSub, ok := subVal.(json.Number); ok {
			if int64Sub, err := jsonSub.Int64(); err == nil {
				return int64Sub, nil
			}
		}
		if subStr, ok := subVal.(string); ok {
			return strconv.ParseInt(subStr, 10, 64)
		}
	}

	// 调试信息：列出context中的所有可能的key
	logx.WithContext(l.ctx).Error("无法从JWT中获取用户ID，调试context内容:")
	// 注意：这里不能直接遍历context，但可以尝试常见的key
	debugKeys := []string{"user_id", "userId", "uid", "sub", "id"}
	for _, key := range debugKeys {
		if val := l.ctx.Value(key); val != nil {
			logx.WithContext(l.ctx).Errorf("发现key '%s': %v (类型: %T)", key, val, val)
		}
	}

	return 0, fmt.Errorf("无法从JWT中获取用户ID")
}
