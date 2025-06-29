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

type UpdateUserStatusLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewUpdateUserStatusLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateUserStatusLogic {
	return &UpdateUserStatusLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *UpdateUserStatusLogic) UpdateUserStatus(in *appuser.UpdateUserStatusReq) (*appuser.AppUserInfo, error) {
	// 参数验证
	if in.UserId <= 0 {
		return nil, errors.New("用户ID无效")
	}

	// 验证状态值
	if in.Status < 1 || in.Status > 3 {
		return nil, errors.New("状态值无效，必须为1(正常)、2(冻结)、3(禁用)")
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

	// 检查状态是否有变化
	if int32(existUser.Status) == in.Status {
		// 状态未变化，直接返回当前用户信息
		return utils.ModelToProto(existUser), nil
	}

	// 更新用户状态
	err = l.svcCtx.AppUsersModel.UpdateStatusById(l.ctx, in.UserId, in.Status)
	if err != nil {
		l.Errorf("更新用户状态失败: %v", err)
		return nil, errors.New("更新用户状态失败")
	}

	// 记录状态变更日志
	statusName := map[int32]string{
		1: "正常",
		2: "冻结",
		3: "禁用",
	}
	oldStatusName := statusName[int32(existUser.Status)]
	newStatusName := statusName[in.Status]
	l.Infof("用户[%s]状态变更: %s -> %s, 原因: %s", existUser.Phone, oldStatusName, newStatusName, in.Reason)

	// 查询更新后的用户信息
	updatedUser, err := l.svcCtx.AppUsersModel.FindOne(l.ctx, uint64(in.UserId))
	if err != nil {
		l.Errorf("查询更新后的用户失败: %v", err)
		return nil, errors.New("更新用户状态失败")
	}

	// 转换为Proto格式并返回
	return utils.ModelToProto(updatedUser), nil
}
