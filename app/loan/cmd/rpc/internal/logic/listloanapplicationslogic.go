package logic

import (
	"context"
	"database/sql"
	"fmt"

	"rpc/internal/svc"
	"rpc/loan"

	"github.com/zeromicro/go-zero/core/logx"
)

type ListLoanApplicationsLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewListLoanApplicationsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ListLoanApplicationsLogic {
	return &ListLoanApplicationsLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *ListLoanApplicationsLogic) ListLoanApplications(in *loan.ListLoanApplicationsReq) (*loan.ListLoanApplicationsResp, error) {
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
	total, err := l.svcCtx.LoanApplicationsModel.CountWithConditions(l.ctx, whereClause, args)
	if err != nil {
		l.Errorf("查询申请总数失败: %v", err)
		return &loan.ListLoanApplicationsResp{
			Code:    500,
			Message: "查询申请失败",
		}, nil
	}

	// 查询分页数据
	offset := (in.Page - 1) * in.Size
	applications, err := l.svcCtx.LoanApplicationsModel.ListWithConditions(l.ctx, whereClause, args, in.Size, offset)
	if err != nil && err != sql.ErrNoRows {
		l.Errorf("查询申请列表失败: %v", err)
		return &loan.ListLoanApplicationsResp{
			Code:    500,
			Message: "查询申请列表失败",
		}, nil
	}

	// 转换为响应格式
	var applicationList []*loan.LoanApplicationInfo
	for _, app := range applications {
		applicationInfo := &loan.LoanApplicationInfo{
			Id:            int64(app.Id),
			ApplicationId: app.ApplicationId,
			UserId:        int64(app.UserId),
			ApplicantName: app.ApplicantName,
			ProductId:     int64(app.ProductId),
			Name:          app.Name,
			Type:          app.Type,
			Amount:        app.Amount,
			Duration:      int32(app.Duration),
			Purpose:       app.Purpose.String,
			Status:        app.Status,
			CreatedAt:     app.CreatedAt.Unix(),
			UpdatedAt:     app.UpdatedAt.Unix(),
		}
		applicationList = append(applicationList, applicationInfo)
	}

	// 如果没有数据，返回空列表
	if applicationList == nil {
		applicationList = make([]*loan.LoanApplicationInfo, 0)
	}

	return &loan.ListLoanApplicationsResp{
		Code:    200,
		Message: "查询成功",
		List:    applicationList,
		Total:   total,
	}, nil
}
