package logic

import (
	"context"
	"fmt"
	"time"

	"rpc/internal/svc"
	"rpc/leaseproduct"

	"github.com/zeromicro/go-zero/core/logx"
)

type CheckInventoryAvailabilityLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewCheckInventoryAvailabilityLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CheckInventoryAvailabilityLogic {
	return &CheckInventoryAvailabilityLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 库存检查
func (l *CheckInventoryAvailabilityLogic) CheckInventoryAvailability(in *leaseproduct.CheckInventoryAvailabilityReq) (*leaseproduct.CheckInventoryAvailabilityResp, error) {
	// 参数验证
	if in.ProductCode == "" {
		return nil, fmt.Errorf("产品编码不能为空")
	}

	if in.Quantity <= 0 {
		return nil, fmt.Errorf("数量必须大于0")
	}

	// 验证日期格式和合理性
	startDate, err := time.Parse("2006-01-02", in.StartDate)
	if err != nil {
		return nil, fmt.Errorf("开始日期格式错误，应为YYYY-MM-DD")
	}

	endDate, err := time.Parse("2006-01-02", in.EndDate)
	if err != nil {
		return nil, fmt.Errorf("结束日期格式错误，应为YYYY-MM-DD")
	}

	if startDate.After(endDate) {
		return nil, fmt.Errorf("开始日期不能晚于结束日期")
	}

	if startDate.Before(time.Now().Truncate(24 * time.Hour)) {
		return nil, fmt.Errorf("开始日期不能早于今天")
	}

	// 查询产品信息
	product, err := l.svcCtx.LeaseProductModel.FindOneByProductCode(l.ctx, in.ProductCode)
	if err != nil {
		l.Errorf("查询产品失败: %v", err)
		return nil, fmt.Errorf("产品不存在")
	}

	// 检查产品状态
	if product.Status != 1 {
		return &leaseproduct.CheckInventoryAvailabilityResp{
			Available:      false,
			AvailableCount: 0,
		}, nil
	}

	// 计算租期天数
	duration := int(endDate.Sub(startDate).Hours()/24) + 1

	// 检查租期是否在允许范围内
	if int32(duration) < int32(product.MinDuration) {
		return nil, fmt.Errorf("租期不能少于最小租期")
	}

	if int32(duration) > int32(product.MaxDuration) {
		return nil, fmt.Errorf("租期不能超过最大租期")
	}

	// TODO: 这里应该根据已有的租赁记录计算实际可用库存
	// 需要查询在指定时间段内已被预订的数量
	// 简化处理：直接使用当前可用库存
	availableCount := int32(product.AvailableCount)
	available := availableCount >= in.Quantity

	return &leaseproduct.CheckInventoryAvailabilityResp{
		Available:      available,
		AvailableCount: availableCount,
	}, nil
}
