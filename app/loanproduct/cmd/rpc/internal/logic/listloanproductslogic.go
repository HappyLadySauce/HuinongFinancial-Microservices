package logic

import (
	"context"
	"database/sql"
	"fmt"

	"rpc/internal/svc"
	"rpc/loanproduct"

	"github.com/zeromicro/go-zero/core/logx"
)

type ListLoanProductsLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewListLoanProductsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ListLoanProductsLogic {
	return &ListLoanProductsLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *ListLoanProductsLogic) ListLoanProducts(in *loanproduct.ListLoanProductsReq) (*loanproduct.ListLoanProductsResp, error) {
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

	// 只查询上架的产品(status=1)，如果没有指定状态的话
	if in.Status == 0 {
		conditions = append(conditions, "status = ?")
		args = append(args, 1)
	} else {
		conditions = append(conditions, "status = ?")
		args = append(args, in.Status)
	}

	if in.Type != "" {
		conditions = append(conditions, "type = ?")
		args = append(args, in.Type)
	}

	if in.Keyword != "" {
		conditions = append(conditions, "(name LIKE ? OR description LIKE ?)")
		keyword := "%" + in.Keyword + "%"
		args = append(args, keyword, keyword)
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
	total, err := l.svcCtx.LoanProductModel.CountWithConditions(l.ctx, whereClause, args)
	if err != nil {
		l.Errorf("查询产品总数失败: %v", err)
		return &loanproduct.ListLoanProductsResp{
			Code:    500,
			Message: "查询产品失败",
		}, nil
	}

	// 查询分页数据
	offset := (in.Page - 1) * in.Size
	productRows, err := l.svcCtx.LoanProductModel.ListWithConditions(l.ctx, whereClause, args, in.Size, offset)
	if err != nil && err != sql.ErrNoRows {
		l.Errorf("查询产品列表失败: %v", err)
		return &loanproduct.ListLoanProductsResp{
			Code:    500,
			Message: "查询产品列表失败",
		}, nil
	}

	// 转换为响应格式
	var products []*loanproduct.LoanProductInfo
	for _, row := range productRows {
		products = append(products, &loanproduct.LoanProductInfo{
			Id:           int64(row.Id),
			ProductCode:  row.ProductCode,
			Name:         row.Name,
			Type:         row.Type,
			MaxAmount:    row.MaxAmount,
			MinAmount:    row.MinAmount,
			MaxDuration:  int32(row.MaxDuration),
			MinDuration:  int32(row.MinDuration),
			InterestRate: row.InterestRate,
			Description:  row.Description,
			Status:       int32(row.Status),
			CreatedAt:    row.CreatedAt.Unix(),
			UpdatedAt:    row.UpdatedAt.Unix(),
		})
	}

	// 如果没有数据，返回空列表
	if products == nil {
		products = make([]*loanproduct.LoanProductInfo, 0)
	}

	return &loanproduct.ListLoanProductsResp{
		Code:    200,
		Message: "查询成功",
		List:    products,
		Total:   total,
	}, nil
}
