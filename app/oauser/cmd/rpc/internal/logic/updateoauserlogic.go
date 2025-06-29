package logic

import (
	"context"
	"errors"
	"model"
	"rpc/internal/svc"
	"rpc/internal/utils"
	"rpc/oauser"

	"github.com/zeromicro/go-zero/core/logx"
)

type UpdateOAUserLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewUpdateOAUserLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateOAUserLogic {
	return &UpdateOAUserLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *UpdateOAUserLogic) UpdateOAUser(in *oauser.UpdateOAUserReq) (*oauser.OAUserInfo, error) {
	// 参数验证
	if in.UserId <= 0 {
		return nil, errors.New("用户ID无效")
	}

	// 检查用户是否存在
	existUser, err := l.svcCtx.OaUsersModel.FindOne(l.ctx, uint64(in.UserId))
	if err != nil {
		if err == model.ErrNotFound {
			return nil, errors.New("用户不存在")
		}
		l.Errorf("查询用户失败: %v", err)
		return nil, errors.New("查询用户失败")
	}

	// 验证状态值
	if in.Status < 1 || in.Status > 2 {
		return nil, errors.New("状态值无效，必须为1(正常)、2(禁用)")
	}

	// 处理角色列表
	rolesStr := utils.RolesToString(in.Roles)

	// 构建更新数据
	updateUser := &model.OaUsers{
		Id:       existUser.Id,
		Username: existUser.Username,
		Password: existUser.Password,
		Name:     in.Name,
		Email:    in.Email,
		Mobile:   in.Mobile,
		Roles:    rolesStr,
		Status:   uint64(in.Status),
	}

	// 更新数据库
	err = l.svcCtx.OaUsersModel.Update(l.ctx, updateUser)
	if err != nil {
		l.Errorf("更新用户失败: %v", err)
		return nil, errors.New("更新用户失败")
	}

	// 查询更新后的用户信息
	updatedUser, err := l.svcCtx.OaUsersModel.FindOne(l.ctx, uint64(in.UserId))
	if err != nil {
		l.Errorf("查询更新后的用户失败: %v", err)
		return nil, errors.New("更新用户失败")
	}

	// 转换为Proto格式并返回
	return utils.ModelToProto(updatedUser), nil
}
