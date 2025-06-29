package logic

import (
	"context"
	"regexp"
	"time"

	"model"
	"rpc/internal/pkg/constants"
	"rpc/internal/pkg/logger"
	"rpc/internal/pkg/utils"
	"rpc/internal/svc"
	"rpc/oauser"

	"github.com/zeromicro/go-zero/core/logx"
)

type RegisterLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewRegisterLogic(ctx context.Context, svcCtx *svc.ServiceContext) *RegisterLogic {
	return &RegisterLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *RegisterLogic) Register(in *oauser.RegisterReq) (*oauser.RegisterResp, error) {
	log := logger.WithContext(l.ctx).WithField("phone", in.Phone)
	log.Info("后台用户注册请求")

	// 参数验证
	if in.Phone == "" || in.Password == "" {
		log.Warn("注册参数不完整")
		return &oauser.RegisterResp{
			Code:    constants.CodeInvalidParams,
			Message: constants.GetMessage(constants.CodeInvalidParams),
		}, nil
	}

	// 验证手机号格式
	phoneRegex := `^1[3-9]\d{9}$`
	matched, _ := regexp.MatchString(phoneRegex, in.Phone)
	if !matched {
		log.Warn("手机号格式无效")
		return &oauser.RegisterResp{
			Code:    constants.CodePhoneInvalid,
			Message: constants.GetMessage(constants.CodePhoneInvalid),
		}, nil
	}

	// 验证密码强度（至少6位）
	if len(in.Password) < 6 {
		log.Warn("密码长度不足")
		return &oauser.RegisterResp{
			Code:    constants.CodeInvalidParams,
			Message: "密码长度至少6位",
		}, nil
	}

	// 检查用户是否已存在
	existingUser, err := l.svcCtx.OaUserModel.FindOneByPhone(l.ctx, in.Phone)
	if err != nil && err != model.ErrNotFound {
		log.WithError(err).Error("查询用户失败")
		return &oauser.RegisterResp{
			Code:    constants.CodeInternalError,
			Message: constants.GetMessage(constants.CodeInternalError),
		}, nil
	}

	if existingUser != nil {
		log.Warn("用户已存在")
		return &oauser.RegisterResp{
			Code:    constants.CodeUserExists,
			Message: constants.GetMessage(constants.CodeUserExists),
		}, nil
	}

	// 密码哈希
	passwordHash, err := utils.HashPassword(in.Password)
	if err != nil {
		log.WithError(err).Error("密码哈希失败")
		return &oauser.RegisterResp{
			Code:    constants.CodeInternalError,
			Message: constants.GetMessage(constants.CodeInternalError),
		}, nil
	}

	// 创建新用户
	newUser := &model.OaUsers{
		Phone:        in.Phone,
		PasswordHash: passwordHash,
		Name:         "",
		Nickname:     "",
		Age:          0,
		Gender:       constants.GenderUnknown,
		Role:         constants.RoleOperator, // 默认为普通操作员
		Status:       constants.UserStatusNormal,
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}

	// 插入数据库
	result, err := l.svcCtx.OaUserModel.Insert(l.ctx, newUser)
	if err != nil {
		log.WithError(err).Error("创建用户失败")
		return &oauser.RegisterResp{
			Code:    constants.CodeInternalError,
			Message: constants.GetMessage(constants.CodeInternalError),
		}, nil
	}

	// 获取新用户ID
	userID, err := result.LastInsertId()
	if err != nil {
		log.WithError(err).Error("获取用户ID失败")
		return &oauser.RegisterResp{
			Code:    constants.CodeInternalError,
			Message: constants.GetMessage(constants.CodeInternalError),
		}, nil
	}

	// 生成 JWT token
	token, err := utils.GenerateToken(
		userID,
		in.Phone,
		"oa",
		l.svcCtx.Config.JwtAuth.AccessSecret,
		l.svcCtx.Config.JwtAuth.AccessExpire,
	)
	if err != nil {
		log.WithError(err).Error("生成token失败")
		return &oauser.RegisterResp{
			Code:    constants.CodeInternalError,
			Message: constants.GetMessage(constants.CodeInternalError),
		}, nil
	}

	log.WithField("user_id", userID).Info("后台用户注册成功")
	return &oauser.RegisterResp{
		Code:    constants.CodeSuccess,
		Message: constants.GetMessage(constants.CodeSuccess),
		Token:   token,
	}, nil
}
