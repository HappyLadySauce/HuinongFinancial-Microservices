package logic

import (
	"context"
	"errors"
	"rpc/appuser"
	"rpc/internal/svc"
	"rpc/internal/utils"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetUserListLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetUserListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetUserListLogic {
	return &GetUserListLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GetUserListLogic) GetUserList(in *appuser.GetUserListReq) (*appuser.GetUserListResp, error) {
	// 参数验证和默认值设置
	page := in.Page
	if page <= 0 {
		page = 1
	}
	
	size := in.Size
	if size <= 0 {
		size = 10
	}
	if size > 100 {
		size = 100 // 限制最大每页条数
	}

	// 查询用户列表
	users, err := l.svcCtx.AppUsersModel.FindListByPage(l.ctx, page, size, in.Keyword, in.Status)
	if err != nil {
		l.Errorf("查询用户列表失败: %v", err)
		return nil, errors.New("查询用户列表失败")
	}

	// 查询总数
	total, err := l.svcCtx.AppUsersModel.CountByConditions(l.ctx, in.Keyword, in.Status)
	if err != nil {
		l.Errorf("查询用户总数失败: %v", err)
		return nil, errors.New("查询用户列表失败")
	}

	// 转换为Proto格式并返回
	return &appuser.GetUserListResp{
		Users: utils.ModelsToProtos(users),
		Total: total,
	}, nil
}
