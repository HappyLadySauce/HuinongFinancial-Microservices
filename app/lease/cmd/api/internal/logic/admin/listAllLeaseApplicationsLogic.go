package admin

import (
	"context"

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
	// 调用 Lease RPC 获取申请列表 (管理员可查看所有申请，不过滤UserId)
	rpcResp, err := l.svcCtx.LeaseRpc.ListLeaseApplications(l.ctx, &leaseclient.ListLeaseApplicationsReq{
		Page:        req.Page,
		Size:        req.Size,
		UserId:      req.UserId,      // 管理员可以按用户ID过滤，也可以不过滤查看所有
		ProductCode: req.ProductCode, // 可以按产品编码过滤
		Status:      req.Status,      // 可以按状态过滤
	})
	if err != nil {
		logx.WithContext(l.ctx).Errorf("调用Lease RPC失败: %v", err)
		return &types.ListLeaseApplicationsResp{
			Code:    500,
			Message: "服务内部错误",
		}, nil
	}

	// 转换 RPC 响应为 API 响应
	resp = &types.ListLeaseApplicationsResp{
		Code:    rpcResp.Code,
		Message: rpcResp.Message,
		Total:   rpcResp.Total,
	}

	// 转换申请列表
	if len(rpcResp.List) > 0 {
		resp.List = make([]types.LeaseApplicationInfo, len(rpcResp.List))
		for i, item := range rpcResp.List {
			resp.List[i] = types.LeaseApplicationInfo{
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
			}
		}
	} else {
		resp.List = make([]types.LeaseApplicationInfo, 0)
	}

	return resp, nil
}
