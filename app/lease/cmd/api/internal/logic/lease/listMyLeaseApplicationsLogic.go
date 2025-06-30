package lease

import (
	"context"
	"strconv"

	"api/internal/svc"
	"api/internal/types"
	"rpc/leaseclient"

	"github.com/zeromicro/go-zero/core/logx"
)

type ListMyLeaseApplicationsLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewListMyLeaseApplicationsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ListMyLeaseApplicationsLogic {
	return &ListMyLeaseApplicationsLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ListMyLeaseApplicationsLogic) ListMyLeaseApplications(req *types.ListLeaseApplicationsReq) (resp *types.ListLeaseApplicationsResp, err error) {
	// 获取当前用户ID (从JWT中获取)
	userIdStr := l.ctx.Value("userId").(string)
	userId, err := strconv.ParseInt(userIdStr, 10, 64)
	if err != nil {
		logx.WithContext(l.ctx).Errorf("解析用户ID失败: %v", err)
		return nil, err
	}

	// 调用 Lease RPC 获取申请列表
	rpcResp, err := l.svcCtx.LeaseRpc.ListLeaseApplications(l.ctx, &leaseclient.ListLeaseApplicationsReq{
		Page:        req.Page,
		Size:        req.Size,
		UserId:      userId,
		ProductCode: req.ProductCode,
		Status:      req.Status,
	})
	if err != nil {
		logx.WithContext(l.ctx).Errorf("调用Lease RPC失败: %v", err)
		return nil, err
	}

	// 转换申请列表
	var applications []types.LeaseApplicationInfo
	for _, item := range rpcResp.List {
		applications = append(applications, types.LeaseApplicationInfo{
			Id:              item.Id,
			ApplicationId:   item.ApplicationId,
			UserId:          item.UserId,
			ApplicantName:   item.ApplicantName,
			ProductId:       item.ProductId,
			ProductCode:     item.ProductCode,
			Name:            item.Name,
			Type:            item.Type,
			Machinery:       item.Machinery,
			StartDate:       item.StartDate,
			EndDate:         item.EndDate,
			Duration:        item.Duration,
			DailyRate:       item.DailyRate,
			TotalAmount:     item.TotalAmount,
			Deposit:         item.Deposit,
			DeliveryAddress: item.DeliveryAddress,
			ContactPhone:    item.ContactPhone,
			Purpose:         item.Purpose,
			Status:          item.Status,
			CreatedAt:       item.CreatedAt,
			UpdatedAt:       item.UpdatedAt,
		})
	}

	return &types.ListLeaseApplicationsResp{
		List:  applications,
		Total: rpcResp.Total,
	}, nil
}
