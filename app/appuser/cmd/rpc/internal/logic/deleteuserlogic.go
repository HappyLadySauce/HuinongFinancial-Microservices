package logic

import (
	"context"

	"model"
	"rpc/appuser"
	"rpc/internal/pkg/constants"
	"rpc/internal/pkg/logger"
	"rpc/internal/svc"

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

// 软删除用户（更新状态为删除状态）
func (l *DeleteUserLogic) DeleteUser(in *appuser.DeleteUserReq) (*appuser.DeleteUserResp, error) {
	log := logger.WithContext(l.ctx).WithField("phone", in.Phone)
	log.Info("删除用户请求")

	// 参数验证
	if in.Phone == "" {
		log.Warn("手机号不能为空")
		return &appuser.DeleteUserResp{
			Code:    constants.CodeInvalidParams,
			Message: constants.GetMessage(constants.CodeInvalidParams),
		}, nil
	}

	// 查找用户是否存在
	user, err := l.svcCtx.AppUserModel.FindOneByPhone(l.ctx, in.Phone)
	if err != nil {
		if err == model.ErrNotFound {
			log.Warn("用户不存在")
			return &appuser.DeleteUserResp{
				Code:    constants.CodeUserNotFound,
				Message: constants.GetMessage(constants.CodeUserNotFound),
			}, nil
		}
		log.WithError(err).Error("查询用户失败")
		return &appuser.DeleteUserResp{
			Code:    constants.CodeInternalError,
			Message: constants.GetMessage(constants.CodeInternalError),
		}, nil
	}

	// 检查用户状态，避免重复删除
	if user.Status == 4 { // 假设状态 4 表示已删除
		log.Warn("用户已经被删除")
		return &appuser.DeleteUserResp{
			Code:    constants.CodeUserAlreadyDeleted,
			Message: constants.GetMessage(constants.CodeUserAlreadyDeleted),
		}, nil
	}

	// 软删除：更新用户状态为删除状态
	user.Status = 4 // 状态 4 表示已删除
	err = l.svcCtx.AppUserModel.Update(l.ctx, user)
	if err != nil {
		log.WithError(err).Error("更新用户状态失败")
		return &appuser.DeleteUserResp{
			Code:    constants.CodeInternalError,
			Message: constants.GetMessage(constants.CodeInternalError),
		}, nil
	}

	log.WithField("user_id", user.Id).Info("用户删除成功")
	return &appuser.DeleteUserResp{
		Code:    constants.CodeSuccess,
		Message: constants.GetMessage(constants.CodeSuccess),
	}, nil
}
