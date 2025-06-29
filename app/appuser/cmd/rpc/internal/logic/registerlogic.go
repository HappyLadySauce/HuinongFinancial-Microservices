package logic

import (
	"context"
	"regexp"
	"time"

	"model"
	"rpc/appuser"
	"rpc/internal/pkg/constants"
	"rpc/internal/pkg/logger"
	"rpc/internal/pkg/utils"
	"rpc/internal/svc"

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

// 用户认证管理
func (l *RegisterLogic) Register(in *appuser.RegisterReq) (*appuser.RegisterResp, error) {
	log := logger.WithContext(l.ctx).WithField("phone", in.Phone)
	log.Info("用户注册请求")

	// 参数验证
	if in.Phone == "" || in.Password == "" {
		log.Warn("注册参数不完整")
		return &appuser.RegisterResp{
			Code:    constants.CodeInvalidParams,
			Message: constants.GetMessage(constants.CodeInvalidParams),
		}, nil
	}

	// 验证手机号格式
	phoneRegex := `^1[3-9]\d{9}$`
	matched, _ := regexp.MatchString(phoneRegex, in.Phone)
	if !matched {
		log.Warn("手机号格式无效")
		return &appuser.RegisterResp{
			Code:    constants.CodePhoneInvalid,
			Message: constants.GetMessage(constants.CodePhoneInvalid),
		}, nil
	}

	// 验证密码长度
	if len(in.Password) < 6 {
		log.Warn("密码长度不足")
		return &appuser.RegisterResp{
			Code:    constants.CodeInvalidParams,
			Message: "密码长度不能少于6位",
		}, nil
	}

	// 检查用户是否已存在
	existUser, err := l.svcCtx.AppUserModel.FindOneByPhone(l.ctx, in.Phone)
	if err != nil && err != model.ErrNotFound {
		log.WithError(err).Error("查询用户失败")
		return &appuser.RegisterResp{
			Code:    constants.CodeInternalError,
			Message: constants.GetMessage(constants.CodeInternalError),
		}, nil
	}
	if existUser != nil {
		log.Warn("用户已存在")
		return &appuser.RegisterResp{
			Code:    constants.CodeUserExists,
			Message: constants.GetMessage(constants.CodeUserExists),
		}, nil
	}

	// 加密密码
	hashedPassword, err := utils.HashPassword(in.Password)
	if err != nil {
		log.WithError(err).Error("密码加密失败")
		return &appuser.RegisterResp{
			Code:    constants.CodeInternalError,
			Message: constants.GetMessage(constants.CodeInternalError),
		}, nil
	}

	// 创建用户
	now := time.Now()
	newUser := &model.AppUsers{
		Phone:     in.Phone,
		Password:  hashedPassword,
		Name:      "", // 注册时姓名为空，后续完善资料
		Nickname:  "",
		Age:       0,
		Gender:    0,
		Status:    1, // 默认正常状态
		CreatedAt: now,
		UpdatedAt: now,
	}

	result, err := l.svcCtx.AppUserModel.Insert(l.ctx, newUser)
	if err != nil {
		log.WithError(err).Error("创建用户失败")
		return &appuser.RegisterResp{
			Code:    constants.CodeInternalError,
			Message: constants.GetMessage(constants.CodeInternalError),
		}, nil
	}

	// 获取插入的用户ID
	userID, err := result.LastInsertId()
	if err != nil {
		log.WithError(err).Error("获取用户ID失败")
		return &appuser.RegisterResp{
			Code:    constants.CodeInternalError,
			Message: constants.GetMessage(constants.CodeInternalError),
		}, nil
	}

	// 生成 JWT token
	token, err := utils.GenerateToken(
		userID,
		in.Phone,
		"app",
		l.svcCtx.Config.JwtAuth.AccessSecret,
		l.svcCtx.Config.JwtAuth.AccessExpire,
	)
	if err != nil {
		log.WithError(err).Error("生成token失败")
		return &appuser.RegisterResp{
			Code:    constants.CodeInternalError,
			Message: constants.GetMessage(constants.CodeInternalError),
		}, nil
	}

	log.WithField("user_id", userID).Info("用户注册成功")
	return &appuser.RegisterResp{
		Code:    constants.CodeSuccess,
		Message: constants.GetMessage(constants.CodeSuccess),
		Token:   token,
	}, nil
}
