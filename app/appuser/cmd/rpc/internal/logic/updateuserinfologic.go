package logic

import (
	"context"
	"time"

	"model"
	"rpc/appuser"
	"rpc/internal/pkg/constants"
	"rpc/internal/pkg/logger"
	"rpc/internal/svc"

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

// 用户信息管理
func (l *UpdateUserInfoLogic) UpdateUserInfo(in *appuser.UpdateUserInfoReq) (*appuser.UpdateUserInfoResp, error) {
	log := logger.WithContext(l.ctx).WithField("user_id", in.UserInfo.Id)
	log.Info("更新用户信息请求")

	// 参数验证
	if in.UserInfo == nil || in.UserInfo.Id <= 0 {
		log.Warn("用户信息参数无效")
		return &appuser.UpdateUserInfoResp{
			Code:    constants.CodeInvalidParams,
			Message: constants.GetMessage(constants.CodeInvalidParams),
		}, nil
	}

	// 检查用户是否存在
	existUser, err := l.svcCtx.AppUserModel.FindOne(l.ctx, uint64(in.UserInfo.Id))
	if err != nil {
		if err == model.ErrNotFound {
			log.Warn("用户不存在")
			return &appuser.UpdateUserInfoResp{
				Code:    constants.CodeUserNotFound,
				Message: constants.GetMessage(constants.CodeUserNotFound),
			}, nil
		}
		log.WithError(err).Error("查询用户失败")
		return &appuser.UpdateUserInfoResp{
			Code:    constants.CodeInternalError,
			Message: constants.GetMessage(constants.CodeInternalError),
		}, nil
	}

	// 更新用户信息
	updatedUser := &model.AppUsers{
		Id:         uint64(in.UserInfo.Id),
		Phone:      existUser.Phone,    // 手机号不允许修改
		Password:   existUser.Password, // 密码不在此处修改
		Name:       in.UserInfo.Name,
		Nickname:   in.UserInfo.Nickname,
		Age:        uint64(in.UserInfo.Age),
		Gender:     uint64(in.UserInfo.Gender),
		Occupation: in.UserInfo.Occupation,
		Address:    in.UserInfo.Address,
		Income:     in.UserInfo.Income,
		Status:     existUser.Status,    // 状态管理需要单独接口
		CreatedAt:  existUser.CreatedAt, // 创建时间不变
		UpdatedAt:  time.Now(),
	}

	err = l.svcCtx.AppUserModel.Update(l.ctx, updatedUser)
	if err != nil {
		log.WithError(err).Error("更新用户信息失败")
		return &appuser.UpdateUserInfoResp{
			Code:    constants.CodeInternalError,
			Message: constants.GetMessage(constants.CodeInternalError),
		}, nil
	}

	// 返回更新后的用户信息
	userInfo := &appuser.UserInfo{
		Id:         int64(updatedUser.Id),
		Phone:      updatedUser.Phone,
		Name:       updatedUser.Name,
		Nickname:   updatedUser.Nickname,
		Age:        int32(updatedUser.Age),
		Gender:     int32(updatedUser.Gender),
		Occupation: updatedUser.Occupation,
		Address:    updatedUser.Address,
		Income:     updatedUser.Income,
		Status:     int32(updatedUser.Status),
		CreatedAt:  updatedUser.CreatedAt.Unix(),
		UpdatedAt:  updatedUser.UpdatedAt.Unix(),
	}

	log.Info("更新用户信息成功")
	return &appuser.UpdateUserInfoResp{
		Code:     constants.CodeSuccess,
		Message:  constants.GetMessage(constants.CodeSuccess),
		UserInfo: userInfo,
	}, nil
}
