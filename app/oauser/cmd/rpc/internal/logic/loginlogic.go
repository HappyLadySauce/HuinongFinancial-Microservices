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

type LoginLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewLoginLogic(ctx context.Context, svcCtx *svc.ServiceContext) *LoginLogic {
	return &LoginLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 用户认证管理
func (l *LoginLogic) Login(in *oauser.LoginReq) (*oauser.LoginResp, error) {
	log := logger.WithContext(l.ctx).WithField("phone", in.Phone)
	log.Info("后台用户登录请求")

	// 参数验证
	if in.Phone == "" || in.Password == "" {
		log.Warn("登录参数不完整")
		return &oauser.LoginResp{
			Code:    constants.CodeInvalidParams,
			Message: constants.GetMessage(constants.CodeInvalidParams),
		}, nil
	}

	// 验证手机号格式
	phoneRegex := `^1[3-9]\d{9}$`
	matched, _ := regexp.MatchString(phoneRegex, in.Phone)
	if !matched {
		log.Warn("手机号格式无效")
		return &oauser.LoginResp{
			Code:    constants.CodePhoneInvalid,
			Message: constants.GetMessage(constants.CodePhoneInvalid),
		}, nil
	}

	// 查找用户
	user, err := l.svcCtx.OaUserModel.FindOneByPhone(l.ctx, in.Phone)
	if err != nil {
		if err == model.ErrNotFound {
			log.Warn("用户不存在")
			return &oauser.LoginResp{
				Code:    constants.CodeUserNotFound,
				Message: constants.GetMessage(constants.CodeUserNotFound),
			}, nil
		}
		log.WithError(err).Error("查询用户失败")
		return &oauser.LoginResp{
			Code:    constants.CodeInternalError,
			Message: constants.GetMessage(constants.CodeInternalError),
		}, nil
	}

	// 检查用户状态
	if user.Status == 2 {
		log.Warn("用户账号被禁用")
		return &oauser.LoginResp{
			Code:    constants.CodeUserDisabled,
			Message: constants.GetMessage(constants.CodeUserDisabled),
		}, nil
	}

	// 验证密码
	if !utils.CheckPassword(in.Password, user.PasswordHash) {
		log.Warn("密码错误")
		return &oauser.LoginResp{
			Code:    constants.CodePasswordError,
			Message: constants.GetMessage(constants.CodePasswordError),
		}, nil
	}

	// 生成 JWT token
	token, err := utils.GenerateToken(
		int64(user.Id),
		user.Phone,
		"oa",
		l.svcCtx.Config.JwtAuth.AccessSecret,
		l.svcCtx.Config.JwtAuth.AccessExpire,
	)
	if err != nil {
		log.WithError(err).Error("生成token失败")
		return &oauser.LoginResp{
			Code:    constants.CodeInternalError,
			Message: constants.GetMessage(constants.CodeInternalError),
		}, nil
	}

	log.WithField("user_id", user.Id).WithField("role", user.Role).Info("后台用户登录成功")
	return &oauser.LoginResp{
		Code:    constants.CodeSuccess,
		Message: constants.GetMessage(constants.CodeSuccess),
		Token:   token,
	}, nil
}
