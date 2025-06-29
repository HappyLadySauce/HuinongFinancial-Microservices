package logic

import (
	"context"
	"time"

	"model"
	"rpc/internal/pkg/constants"
	"rpc/internal/pkg/logger"
	"rpc/internal/svc"
	"rpc/oauser"

	"github.com/zeromicro/go-zero/core/logx"
)

type UpdateUserInfoLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewUpdateUserInfoLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateUserInfoLogic {
	return &UpdateUserInfoLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *UpdateUserInfoLogic) UpdateUserInfo(in *oauser.UpdateUserInfoReq) (*oauser.UpdateUserInfoResp, error) {
	log := logger.WithContext(l.ctx).WithField("phone", in.UserInfo.Phone)
	log.Info("更新后台用户信息请求")

	// 参数验证
	if in.UserInfo == nil || in.UserInfo.Phone == "" {
		log.Warn("用户信息参数无效")
		return &oauser.UpdateUserInfoResp{
			Code:    constants.CodeInvalidParams,
			Message: constants.GetMessage(constants.CodeInvalidParams),
		}, nil
	}

	userInfo := in.UserInfo

	// 查找现有用户
	existingUser, err := l.svcCtx.OaUserModel.FindOneByPhone(l.ctx, userInfo.Phone)
	if err != nil {
		if err == model.ErrNotFound {
			log.Warn("用户不存在")
			return &oauser.UpdateUserInfoResp{
				Code:    constants.CodeUserNotFound,
				Message: constants.GetMessage(constants.CodeUserNotFound),
			}, nil
		}
		log.WithError(err).Error("查询用户失败")
		return &oauser.UpdateUserInfoResp{
			Code:    constants.CodeInternalError,
			Message: constants.GetMessage(constants.CodeInternalError),
		}, nil
	}

	// 检查用户状态
	if existingUser.Status == constants.UserStatusDisabled {
		log.Warn("用户账号被禁用")
		return &oauser.UpdateUserInfoResp{
			Code:    constants.CodeUserDisabled,
			Message: constants.GetMessage(constants.CodeUserDisabled),
		}, nil
	}

	// 更新用户信息
	updateUser := &model.OaUsers{
		Id:           existingUser.Id,
		Phone:        existingUser.Phone,        // 手机号不允许修改
		PasswordHash: existingUser.PasswordHash, // 密码不在此处修改
		Name:         userInfo.Name,
		Nickname:     userInfo.Nickname,
		Age:          uint64(userInfo.Age),
		Gender:       uint64(userInfo.Gender),
		Role:         userInfo.Role,
		Status:       uint64(userInfo.Status),
		CreatedAt:    existingUser.CreatedAt,
		UpdatedAt:    time.Now(),
	}

	// 验证角色有效性
	if updateUser.Role != constants.RoleAdmin && updateUser.Role != constants.RoleOperator {
		log.Warn("无效的用户角色")
		return &oauser.UpdateUserInfoResp{
			Code:    constants.CodeInvalidParams,
			Message: "无效的用户角色",
		}, nil
	}

	// 验证状态有效性
	if updateUser.Status != constants.UserStatusNormal && updateUser.Status != constants.UserStatusDisabled {
		log.Warn("无效的用户状态")
		return &oauser.UpdateUserInfoResp{
			Code:    constants.CodeInvalidParams,
			Message: "无效的用户状态",
		}, nil
	}

	// 执行更新
	err = l.svcCtx.OaUserModel.Update(l.ctx, updateUser)
	if err != nil {
		log.WithError(err).Error("更新用户失败")
		return &oauser.UpdateUserInfoResp{
			Code:    constants.CodeInternalError,
			Message: constants.GetMessage(constants.CodeInternalError),
		}, nil
	}

	// 构建响应用户信息
	responseUserInfo := &oauser.UserInfo{
		Id:        int64(updateUser.Id),
		Phone:     updateUser.Phone,
		Name:      updateUser.Name,
		Nickname:  updateUser.Nickname,
		Age:       int32(updateUser.Age),
		Gender:    int32(updateUser.Gender),
		Role:      updateUser.Role,
		Status:    int32(updateUser.Status),
		CreatedAt: updateUser.CreatedAt.Unix(),
		UpdatedAt: updateUser.UpdatedAt.Unix(),
	}

	log.WithField("user_id", updateUser.Id).Info("更新用户信息成功")
	return &oauser.UpdateUserInfoResp{
		Code:     constants.CodeSuccess,
		Message:  constants.GetMessage(constants.CodeSuccess),
		UserInfo: responseUserInfo,
	}, nil
}
