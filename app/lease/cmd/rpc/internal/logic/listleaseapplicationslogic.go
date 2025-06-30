package logic

import (
	"context"
	"fmt"

	"rpc/internal/svc"
	"rpc/lease"

	"github.com/zeromicro/go-zero/core/logx"
)

type ListLeaseApplicationsLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewListLeaseApplicationsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ListLeaseApplicationsLogic {
	return &ListLeaseApplicationsLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *ListLeaseApplicationsLogic) ListLeaseApplications(in *lease.ListLeaseApplicationsReq) (*lease.ListLeaseApplicationsResp, error) {
	// 设置默认分页参数
	page := in.Page
	size := in.Size
	if page <= 0 {
		page = 1
	}
	if size <= 0 {
		size = 10
	}

	// 构建查询条件
	var whereClause string
	var args []interface{}

	// 构建WHERE条件
	conditions := []string{}

	if in.UserId > 0 {
		conditions = append(conditions, "user_id = ?")
		args = append(args, in.UserId)
	}

	if in.ProductCode != "" {
		conditions = append(conditions, "product_code = ?")
		args = append(args, in.ProductCode)
	}

	if in.Status != "" {
		conditions = append(conditions, "status = ?")
		args = append(args, in.Status)
	}

	if len(conditions) > 0 {
		whereClause = "WHERE " + fmt.Sprintf("%s", conditions[0])
		for i := 1; i < len(conditions); i++ {
			whereClause += " AND " + conditions[i]
		}
	}

	// 查询总数
	total, err := l.svcCtx.LeaseApplicationsModel.CountWithConditions(l.ctx, whereClause, args)
	if err != nil {
		l.Errorf("查询申请总数失败: %v", err)
		return nil, fmt.Errorf("查询申请失败")
	}

	// 查询申请列表
	offset := (page - 1) * size
	applications, err := l.svcCtx.LeaseApplicationsModel.ListWithConditions(l.ctx, whereClause, args, size, offset)
	if err != nil {
		l.Errorf("查询申请列表失败: %v", err)
		return nil, fmt.Errorf("查询申请列表失败")
	}

	// 转换为响应格式
	var applicationList []*lease.LeaseApplicationInfo
	for _, app := range applications {
		applicationList = append(applicationList, &lease.LeaseApplicationInfo{
			Id:              int64(app.Id),
			ApplicationId:   app.ApplicationId,
			UserId:          int64(app.UserId),
			ApplicantName:   app.ApplicantName,
			ProductId:       int64(app.ProductId),
			ProductCode:     app.ProductCode,
			Name:            app.Name,
			Type:            app.Type,
			Machinery:       app.Machinery,
			StartDate:       app.StartDate.Format("2006-01-02"),
			EndDate:         app.EndDate.Format("2006-01-02"),
			Duration:        int32(app.Duration),
			DailyRate:       app.DailyRate,
			TotalAmount:     app.TotalAmount,
			Deposit:         app.Deposit,
			DeliveryAddress: app.DeliveryAddress,
			ContactPhone:    app.ContactPhone,
			Purpose:         app.Purpose.String,
			Status:          app.Status,
			CreatedAt:       app.CreatedAt.Unix(),
			UpdatedAt:       app.UpdatedAt.Unix(),
		})
	}

	return &lease.ListLeaseApplicationsResp{
		List:  applicationList,
		Total: total,
	}, nil
}
