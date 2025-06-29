package model

import (
	"context"
	"database/sql"
	"fmt"
	"strings"

	"github.com/zeromicro/go-zero/core/stores/cache"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

var _ AppUsersModel = (*customAppUsersModel)(nil)

type (
	// AppUsersModel is an interface to be customized, add more methods here,
	// and implement the added methods in customAppUsersModel.
	AppUsersModel interface {
		appUsersModel
		FindOneByPhone(ctx context.Context, phone string) (*AppUsers, error)
		FindListByPage(ctx context.Context, page, size int32, keyword string, status int32) ([]*AppUsers, error)
		CountByConditions(ctx context.Context, keyword string, status int32) (int64, error)
		UpdateStatusById(ctx context.Context, id int64, status int32) error
	}

	customAppUsersModel struct {
		*defaultAppUsersModel
	}
)

// NewAppUsersModel returns a model for the database table.
func NewAppUsersModel(conn sqlx.SqlConn, c cache.CacheConf, opts ...cache.Option) AppUsersModel {
	return &customAppUsersModel{
		defaultAppUsersModel: newAppUsersModel(conn, c, opts...),
	}
}

// FindOneByPhone 根据手机号查找用户
func (m *customAppUsersModel) FindOneByPhone(ctx context.Context, phone string) (*AppUsers, error) {
	var resp AppUsers
	query := fmt.Sprintf("select %s from %s where `phone` = ? limit 1", appUsersRows, m.table)
	err := m.QueryRowNoCacheCtx(ctx, &resp, query, phone)
	switch err {
	case nil:
		return &resp, nil
	case sqlx.ErrNotFound:
		return nil, ErrNotFound
	default:
		return nil, err
	}
}

// FindListByPage 分页查询用户列表
func (m *customAppUsersModel) FindListByPage(ctx context.Context, page, size int32, keyword string, status int32) ([]*AppUsers, error) {
	offset := (page - 1) * size

	var conditions []string
	var args []interface{}

	// 构建查询条件
	if status > 0 {
		conditions = append(conditions, "`status` = ?")
		args = append(args, status)
	}

	if keyword != "" {
		conditions = append(conditions, "(`name` LIKE ? OR `phone` LIKE ? OR `nickname` LIKE ?)")
		likeKeyword := "%" + keyword + "%"
		args = append(args, likeKeyword, likeKeyword, likeKeyword)
	}

	whereClause := ""
	if len(conditions) > 0 {
		whereClause = " WHERE " + strings.Join(conditions, " AND ")
	}

	query := fmt.Sprintf("select %s from %s%s ORDER BY `created_at` DESC limit ?, ?",
		appUsersRows, m.table, whereClause)
	args = append(args, offset, size)

	var resp []*AppUsers
	err := m.QueryRowsNoCacheCtx(ctx, &resp, query, args...)
	return resp, err
}

// CountByConditions 根据条件统计用户数量
func (m *customAppUsersModel) CountByConditions(ctx context.Context, keyword string, status int32) (int64, error) {
	var conditions []string
	var args []interface{}

	if status > 0 {
		conditions = append(conditions, "`status` = ?")
		args = append(args, status)
	}

	if keyword != "" {
		conditions = append(conditions, "(`name` LIKE ? OR `phone` LIKE ? OR `nickname` LIKE ?)")
		likeKeyword := "%" + keyword + "%"
		args = append(args, likeKeyword, likeKeyword, likeKeyword)
	}

	whereClause := ""
	if len(conditions) > 0 {
		whereClause = " WHERE " + strings.Join(conditions, " AND ")
	}

	query := fmt.Sprintf("select count(*) from %s%s", m.table, whereClause)

	var count int64
	err := m.QueryRowNoCacheCtx(ctx, &count, query, args...)
	return count, err
}

// UpdateStatusById 更新用户状态
func (m *customAppUsersModel) UpdateStatusById(ctx context.Context, id int64, status int32) error {
	appUsersIdKey := fmt.Sprintf("%s%v", cacheAppUsersIdPrefix, id)
	_, err := m.ExecCtx(ctx, func(ctx context.Context, conn sqlx.SqlConn) (result sql.Result, err error) {
		query := fmt.Sprintf("update %s set `status` = ?, `updated_at` = CURRENT_TIMESTAMP where `id` = ?", m.table)
		return conn.ExecCtx(ctx, query, status, id)
	}, appUsersIdKey)
	return err
}
