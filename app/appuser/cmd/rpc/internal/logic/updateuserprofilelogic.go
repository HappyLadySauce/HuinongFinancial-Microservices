package logic

import (
	"context"
	"errors"
	"model"
	"rpc/appuser"
	"rpc/internal/svc"
	"rpc/internal/utils"

	"github.com/zeromicro/go-zero/core/logx"
)

type UpdateUserProfileLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewUpdateUserProfileLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateUserProfileLogic {
	return &UpdateUserProfileLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *UpdateUserProfileLogic) UpdateUserProfile(in *appuser.UpdateUserProfileReq) (*appuser.AppUserInfo, error) {
	// 参数验证
	if in.UserId <= 0 {
		return nil, errors.New("用户ID无效")
	}

	// 检查用户是否存在
	existUser, err := l.svcCtx.AppUsersModel.FindOne(l.ctx, uint64(in.UserId))
	if err != nil {
		if err == model.ErrNotFound {
			return nil, errors.New("用户不存在")
		}
		l.Errorf("查询用户失败: %v", err)
		return nil, errors.New("查询用户失败")
	}

	// 构建更新数据
	updateUser := &model.AppUsers{
		Id:         existUser.Id,
		Phone:      existUser.Phone,
		Password:   existUser.Password,
		Name:       existUser.Name,
		Nickname:   in.Nickname,
		Age:        uint64(in.Age),
		Gender:     uint64(in.Gender),
		Occupation: in.Occupation,
		Address:    in.Address,
		Income:     in.Income,
		Status:     existUser.Status,
		CreatedAt:  existUser.CreatedAt,
	}

	// 更新数据库
	err = l.svcCtx.AppUsersModel.Update(l.ctx, updateUser)
	if err != nil {
		l.Errorf("更新用户档案失败: %v", err)
		return nil, errors.New("更新用户档案失败")
	}

	// 查询更新后的用户信息
	updatedUser, err := l.svcCtx.AppUsersModel.FindOne(l.ctx, uint64(in.UserId))
	if err != nil {
		l.Errorf("查询更新后的用户失败: %v", err)
		return nil, errors.New("更新用户档案失败")
	}

	// 转换为Proto格式并返回
	return utils.ModelToProto(updatedUser), nil
}
