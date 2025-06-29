package logic

import (
	"context"
	"regexp"

	"model"
	"rpc/internal/pkg/constants"
	"rpc/internal/pkg/logger"
	"rpc/internal/svc"
	"rpc/oauser"

	"github.com/zeromicro/go-zero/core/logx"
)

type DeleteUserLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewDeleteUserLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DeleteUserLogic {
	return &DeleteUserLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *DeleteUserLogic) DeleteUser(in *oauser.DeleteUserReq) (*oauser.DeleteUserResp, error) {
	log := logger.WithContext(l.ctx).WithField("phone", in.Phone)
	log.Info("删除后台用户请求")

	// 参数验证
	if in.Phone == "" {
		log.Warn("手机号参数为空")
		return &oauser.DeleteUserResp{
			Code:    constants.CodeInvalidParams,
			Message: constants.GetMessage(constants.CodeInvalidParams),
		}, nil
	}

	// 权限验证：只有管理员才能删除用户
	if in.CallerToken == "" {
		log.Warn("缺少调用者身份验证信息")
		return &oauser.DeleteUserResp{
			Code:    constants.CodeUnauthorized,
			Message: "缺少身份验证信息",
		}, nil
	}

	// 验证调用者Token
	claims, err := l.svcCtx.JwtUtils.ValidateAndGetClaims(in.CallerToken)
	if err != nil {
		log.WithError(err).Warn("调用者Token验证失败")
		return &oauser.DeleteUserResp{
			Code:    constants.CodeUnauthorized,
			Message: "身份验证失败",
		}, nil
	}

	// 检查调用者是否为管理员
	if !claims.IsAdmin() {
		log.WithField("caller_role", claims.Role).Warn("权限不足，只有管理员才能删除用户")
		return &oauser.DeleteUserResp{
			Code:    constants.CodeForbidden,
			Message: "权限不足，只有管理员才能删除用户",
		}, nil
	}

	log.WithFields(map[string]interface{}{
		"caller_phone": claims.Phone,
		"caller_id":    claims.UserID,
		"caller_role":  claims.Role,
	}).Info("管理员权限验证通过")

	// 验证手机号格式
	phoneRegex := `^1[3-9]\d{9}$`
	matched, _ := regexp.MatchString(phoneRegex, in.Phone)
	if !matched {
		log.Warn("手机号格式无效")
		return &oauser.DeleteUserResp{
			Code:    constants.CodePhoneInvalid,
			Message: constants.GetMessage(constants.CodePhoneInvalid),
		}, nil
	}

	// 查找用户
	user, err := l.svcCtx.OaUserModel.FindOneByPhone(l.ctx, in.Phone)
	if err != nil {
		if err == model.ErrNotFound {
			log.Warn("用户不存在")
			return &oauser.DeleteUserResp{
				Code:    constants.CodeUserNotFound,
				Message: constants.GetMessage(constants.CodeUserNotFound),
			}, nil
		}
		log.WithError(err).Error("查询用户失败")
		return &oauser.DeleteUserResp{
			Code:    constants.CodeInternalError,
			Message: constants.GetMessage(constants.CodeInternalError),
		}, nil
	}

	// 业务逻辑：不直接删除用户，而是禁用用户账号
	// 这是为了保持数据完整性和审计需要
	if user.Status == constants.UserStatusDisabled {
		log.Warn("用户已被禁用")
		return &oauser.DeleteUserResp{
			Code:    constants.CodeUserDisabled,
			Message: "用户已被禁用",
		}, nil
	}

	// 更新用户状态为禁用
	user.Status = constants.UserStatusDisabled
	err = l.svcCtx.OaUserModel.Update(l.ctx, user)
	if err != nil {
		log.WithError(err).Error("禁用用户失败")
		return &oauser.DeleteUserResp{
			Code:    constants.CodeInternalError,
			Message: constants.GetMessage(constants.CodeInternalError),
		}, nil
	}

	log.WithField("user_id", user.Id).Info("用户账号已禁用")
	return &oauser.DeleteUserResp{
		Code:    constants.CodeSuccess,
		Message: "用户账号已禁用",
	}, nil
}
