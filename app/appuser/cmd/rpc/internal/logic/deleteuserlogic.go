package logic

import (
	"context"

	"model"
	"rpc/appuser"
	"rpc/internal/pkg/constants"
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
	l.Infof("删除用户请求, phone: %s", in.Phone)

	// 参数验证
	if in.Phone == "" {
		l.Infof("手机号不能为空")
		return nil, constants.ErrInvalidParams
	}

	// 查找用户是否存在
	user, err := l.svcCtx.AppUserModel.FindOneByPhone(l.ctx, in.Phone)
	if err != nil {
		if err == model.ErrNotFound {
			l.Infof("用户不存在")
			return nil, constants.ErrUserNotFound
		}
		l.Errorf("查询用户失败: %v", err)
		return nil, constants.ErrInternalError
	}

	// 检查用户状态，避免重复删除
	if user.Status == 4 { // 假设状态 4 表示已删除
		l.Infof("用户已经被删除")
		return nil, constants.ErrUserAlreadyDeleted
	}

	// 软删除：更新用户状态为删除状态
	user.Status = 4 // 状态 4 表示已删除
	err = l.svcCtx.AppUserModel.Update(l.ctx, user)
	if err != nil {
		l.Errorf("更新用户状态失败: %v", err)
		return nil, constants.ErrInternalError
	}

	l.Infof("用户删除成功, user_id: %d", user.Id)
	return &appuser.DeleteUserResp{}, nil
}
