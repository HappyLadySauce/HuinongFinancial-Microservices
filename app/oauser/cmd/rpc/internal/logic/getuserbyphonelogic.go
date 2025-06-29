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

type GetUserByPhoneLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetUserByPhoneLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetUserByPhoneLogic {
	return &GetUserByPhoneLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 用户信息管理
func (l *GetUserByPhoneLogic) GetUserByPhone(in *oauser.GetUserInfoReq) (*oauser.GetUserInfoResp, error) {
	log := logger.WithContext(l.ctx).WithField("phone", in.Phone)
	log.Info("获取后台用户信息请求")

	// 参数验证
	if in.Phone == "" {
		log.Warn("手机号参数为空")
		return &oauser.GetUserInfoResp{
			Code:    constants.CodeInvalidParams,
			Message: constants.GetMessage(constants.CodeInvalidParams),
		}, nil
	}

	// 验证手机号格式
	phoneRegex := `^1[3-9]\d{9}$`
	matched, _ := regexp.MatchString(phoneRegex, in.Phone)
	if !matched {
		log.Warn("手机号格式无效")
		return &oauser.GetUserInfoResp{
			Code:    constants.CodePhoneInvalid,
			Message: constants.GetMessage(constants.CodePhoneInvalid),
		}, nil
	}

	// 查找用户
	user, err := l.svcCtx.OaUserModel.FindOneByPhone(l.ctx, in.Phone)
	if err != nil {
		if err == model.ErrNotFound {
			log.Warn("用户不存在")
			return &oauser.GetUserInfoResp{
				Code:    constants.CodeUserNotFound,
				Message: constants.GetMessage(constants.CodeUserNotFound),
			}, nil
		}
		log.WithError(err).Error("查询用户失败")
		return &oauser.GetUserInfoResp{
			Code:    constants.CodeInternalError,
			Message: constants.GetMessage(constants.CodeInternalError),
		}, nil
	}

	// 构建用户信息响应
	userInfo := &oauser.UserInfo{
		Id:        int64(user.Id),
		Phone:     user.Phone,
		Name:      user.Name,
		Nickname:  user.Nickname,
		Age:       int32(user.Age),
		Gender:    int32(user.Gender),
		Role:      user.Role,
		Status:    int32(user.Status),
		CreatedAt: user.CreatedAt.Unix(),
		UpdatedAt: user.UpdatedAt.Unix(),
	}

	log.WithField("user_id", user.Id).WithField("role", user.Role).Info("获取用户信息成功")
	return &oauser.GetUserInfoResp{
		Code:     constants.CodeSuccess,
		Message:  constants.GetMessage(constants.CodeSuccess),
		UserInfo: userInfo,
	}, nil
}
