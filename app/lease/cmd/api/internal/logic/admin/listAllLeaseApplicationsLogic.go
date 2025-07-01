package admin

import (
	"context"

	"api/internal/breaker"
	"api/internal/svc"
	"api/internal/types"
	"rpc/leaseclient"

	"github.com/zeromicro/go-zero/core/logx"
)

type ListAllLeaseApplicationsLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewListAllLeaseApplicationsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ListAllLeaseApplicationsLogic {
	return &ListAllLeaseApplicationsLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ListAllLeaseApplicationsLogic) ListAllLeaseApplications(req *types.ListLeaseApplicationsReq) (resp *types.ListLeaseApplicationsResp, err error) {
	// 调用 Lease RPC 获取所有申请列表 (不传UserId，获取所有用户的申请) - 使用熔断器
	rpcResp, err := breaker.DoWithBreakerResultAcceptable(l.ctx, "lease-rpc", func() (*leaseclient.ListLeaseApplicationsResp, error) {
		return l.svcCtx.LeaseRpc.ListLeaseApplications(l.ctx, &leaseclient.ListLeaseApplicationsReq{
			Page:        req.Page,
			Size:        req.Size,
			UserId:      req.UserId, // 管理员可以指定用户ID查询
			ProductCode: req.ProductCode,
			Status:      req.Status,
		})
	}, breaker.IsAcceptableError)
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
