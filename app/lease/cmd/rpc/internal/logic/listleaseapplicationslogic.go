package logic

import (
	"context"
	"database/sql"
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
	// 参数验证
	if in.Page <= 0 {
		in.Page = 1
	}
	if in.Size <= 0 {
		in.Size = 10
	}
	if in.Size > 100 {
		in.Size = 100 // 限制最大页面大小
	}

	// 构建查询条件
	var conditions []string
	var args []interface{}

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

	// 构建WHERE子句
	whereClause := ""
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
		return &lease.ListLeaseApplicationsResp{
			Code:    500,
			Message: "查询申请失败",
		}, nil
	}

	// 查询分页数据
	offset := (in.Page - 1) * in.Size
	applications, err := l.svcCtx.LeaseApplicationsModel.ListWithConditions(l.ctx, whereClause, args, in.Size, offset)
	if err != nil && err != sql.ErrNoRows {
		l.Errorf("查询申请列表失败: %v", err)
		return &lease.ListLeaseApplicationsResp{
			Code:    500,
			Message: "查询申请列表失败",
		}, nil
	}

	// 转换为响应格式
	var applicationList []*lease.LeaseApplicationInfo
	for _, app := range applications {
		applicationInfo := &lease.LeaseApplicationInfo{
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
		}
		applicationList = append(applicationList, applicationInfo)
	}

	// 如果没有数据，返回空列表
	if applicationList == nil {
		applicationList = make([]*lease.LeaseApplicationInfo, 0)
	}

	return &lease.ListLeaseApplicationsResp{
		Code:    200,
		Message: "查询成功",
		List:    applicationList,
		Total:   total,
	}, nil
}
