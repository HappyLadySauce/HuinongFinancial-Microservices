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
		return &types.ListLeaseApplicationsResp{
			Code:    400,
			Message: "用户ID无效",
		}, nil
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
