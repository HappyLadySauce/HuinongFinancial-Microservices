package logic

import (
	"context"
	"regexp"

	"model"
	"rpc/internal/pkg/constants"
	"rpc/internal/pkg/logger"
	"rpc/internal/pkg/utils"
	"rpc/internal/svc"
	"rpc/oauser"

	"github.com/zeromicro/go-zero/core/logx"
)

type ChangePasswordLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewChangePasswordLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ChangePasswordLogic {
	return &ChangePasswordLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *ChangePasswordLogic) ChangePassword(in *oauser.ChangePasswordReq) (*oauser.ChangePasswordResp, error) {
	log := logger.WithContext(l.ctx).WithField("phone", in.Phone)
	log.Info("修改密码请求")

	// 参数验证
	if in.Phone == "" || in.OldPassword == "" || in.NewPassword == "" {
		log.Warn("修改密码参数不完整")
		return &oauser.ChangePasswordResp{
			Code:    constants.CodeInvalidParams,
			Message: constants.GetMessage(constants.CodeInvalidParams),
		}, nil
	}

	// 验证手机号格式
	phoneRegex := `^1[3-9]\d{9}$`
	matched, _ := regexp.MatchString(phoneRegex, in.Phone)
	if !matched {
		log.Warn("手机号格式无效")
		return &oauser.ChangePasswordResp{
			Code:    constants.CodePhoneInvalid,
			Message: constants.GetMessage(constants.CodePhoneInvalid),
		}, nil
	}

	// 验证新密码强度（至少6位）
	if len(in.NewPassword) < 6 {
		log.Warn("新密码长度不足")
		return &oauser.ChangePasswordResp{
			Code:    constants.CodeInvalidParams,
			Message: "新密码长度至少6位",
		}, nil
	}

	// 检查新旧密码不能相同
	if in.OldPassword == in.NewPassword {
		log.Warn("新旧密码相同")
		return &oauser.ChangePasswordResp{
			Code:    constants.CodeInvalidParams,
			Message: "新旧密码不能相同",
		}, nil
	}

	// 查找用户
	user, err := l.svcCtx.OaUserModel.FindOneByPhone(l.ctx, in.Phone)
	if err != nil {
		if err == model.ErrNotFound {
			log.Warn("用户不存在")
			return &oauser.ChangePasswordResp{
				Code:    constants.CodeUserNotFound,
				Message: constants.GetMessage(constants.CodeUserNotFound),
			}, nil
		}
		log.WithError(err).Error("查询用户失败")
		return &oauser.ChangePasswordResp{
			Code:    constants.CodeInternalError,
			Message: constants.GetMessage(constants.CodeInternalError),
		}, nil
	}

	// 检查用户状态
	if user.Status == constants.UserStatusDisabled {
		log.Warn("用户账号被禁用")
		return &oauser.ChangePasswordResp{
			Code:    constants.CodeUserDisabled,
			Message: constants.GetMessage(constants.CodeUserDisabled),
		}, nil
	}

	// 验证旧密码
	if !utils.CheckPassword(in.OldPassword, user.PasswordHash) {
		log.Warn("原密码错误")
		return &oauser.ChangePasswordResp{
			Code:    constants.CodePasswordError,
			Message: "原密码错误",
		}, nil
	}

	// 生成新密码哈希
	newPasswordHash, err := utils.HashPassword(in.NewPassword)
	if err != nil {
		log.WithError(err).Error("新密码哈希失败")
		return &oauser.ChangePasswordResp{
			Code:    constants.CodeInternalError,
			Message: constants.GetMessage(constants.CodeInternalError),
		}, nil
	}

	// 更新密码
	user.PasswordHash = newPasswordHash
	err = l.svcCtx.OaUserModel.Update(l.ctx, user)
	if err != nil {
		log.WithError(err).Error("更新密码失败")
		return &oauser.ChangePasswordResp{
			Code:    constants.CodeInternalError,
			Message: constants.GetMessage(constants.CodeInternalError),
		}, nil
	}

	log.WithField("user_id", user.Id).Info("密码修改成功")
	return &oauser.ChangePasswordResp{
		Code:    constants.CodeSuccess,
		Message: constants.GetMessage(constants.CodeSuccess),
	}, nil
}
