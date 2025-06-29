package model

import (
	"context"
	"database/sql"
	"fmt"
	"strings"

	"github.com/zeromicro/go-zero/core/stores/cache"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

var _ OaUsersModel = (*customOaUsersModel)(nil)

type (
	// OaUsersModel is an interface to be customized, add more methods here,
	// and implement the added methods in customOaUsersModel.
	OaUsersModel interface {
		oaUsersModel
		FindOneByUsername(ctx context.Context, username string) (*OaUsers, error)
		FindListByPage(ctx context.Context, page, size int32, keyword string, status int32) ([]*OaUsers, error)
		CountByConditions(ctx context.Context, keyword string, status int32) (int64, error)
		UpdateStatusById(ctx context.Context, id int64, status int32) error
	}

	customOaUsersModel struct {
		*defaultOaUsersModel
	}
)

// NewOaUsersModel returns a model for the database table.
func NewOaUsersModel(conn sqlx.SqlConn, c cache.CacheConf, opts ...cache.Option) OaUsersModel {
	return &customOaUsersModel{
		defaultOaUsersModel: newOaUsersModel(conn, c, opts...),
	}
}

// FindOneByUsername 根据用户名查找用户
func (m *customOaUsersModel) FindOneByUsername(ctx context.Context, username string) (*OaUsers, error) {
	var resp OaUsers
	query := fmt.Sprintf("select %s from %s where `username` = ? limit 1", oaUsersRows, m.table)
	err := m.QueryRowNoCacheCtx(ctx, &resp, query, username)
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
func (m *customOaUsersModel) FindListByPage(ctx context.Context, page, size int32, keyword string, status int32) ([]*OaUsers, error) {
	offset := (page - 1) * size
	
	var conditions []string
	var args []interface{}
	
	// 构建查询条件
	if status > 0 {
		conditions = append(conditions, "`status` = ?")
		args = append(args, status)
	}
	
	if keyword != "" {
		conditions = append(conditions, "(`name` LIKE ? OR `username` LIKE ? OR `email` LIKE ?)")
		likeKeyword := "%" + keyword + "%"
		args = append(args, likeKeyword, likeKeyword, likeKeyword)
	}
	
	whereClause := ""
	if len(conditions) > 0 {
		whereClause = " WHERE " + strings.Join(conditions, " AND ")
	}
	
	query := fmt.Sprintf("select %s from %s%s ORDER BY `created_at` DESC limit ?, ?", 
		oaUsersRows, m.table, whereClause)
	args = append(args, offset, size)
	
	var resp []*OaUsers
	err := m.QueryRowsNoCacheCtx(ctx, &resp, query, args...)
	return resp, err
}

// CountByConditions 根据条件统计用户数量
func (m *customOaUsersModel) CountByConditions(ctx context.Context, keyword string, status int32) (int64, error) {
	var conditions []string
	var args []interface{}
	
	if status > 0 {
		conditions = append(conditions, "`status` = ?")
		args = append(args, status)
	}
	
	if keyword != "" {
		conditions = append(conditions, "(`name` LIKE ? OR `username` LIKE ? OR `email` LIKE ?)")
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
func (m *customOaUsersModel) UpdateStatusById(ctx context.Context, id int64, status int32) error {
	oaUsersIdKey := fmt.Sprintf("%s%v", cacheOaUsersIdPrefix, id)
	_, err := m.ExecCtx(ctx, func(ctx context.Context, conn sqlx.SqlConn) (result sql.Result, err error) {
		query := fmt.Sprintf("update %s set `status` = ?, `updated_at` = CURRENT_TIMESTAMP where `id` = ?", m.table)
		return conn.ExecCtx(ctx, query, status, id)
	}, oaUsersIdKey)
	return err
}
