package logic

import (
	"context"

	"model"
	"rpc/internal/pkg/constants"
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
	l.Infof("删除后台用户请求, phone: %s", in.Phone)

	// 参数验证
	if in.Phone == "" {
		l.Infof("手机号参数为空")
		return &oauser.DeleteUserResp{
			Code:    constants.CodeInvalidParams,
			Message: constants.GetMessage(constants.CodeInvalidParams),
		}, nil
	}

	// 验证调用者权限（如果提供了 caller_token）
	if in.CallerToken != "" {
		claims, err := l.svcCtx.JwtUtils.ValidateAndGetClaims(in.CallerToken)
		if err != nil {
			l.Infof("无效的调用者token")
			return &oauser.DeleteUserResp{
				Code:    constants.CodeUnauthorized,
				Message: constants.GetMessage(constants.CodeUnauthorized),
			}, nil
		}

		// 检查是否为管理员权限
		if claims.Role != constants.RoleAdmin {
			l.Infof("权限不足，只有管理员可以删除用户")
			return &oauser.DeleteUserResp{
				Code:    constants.CodePermissionDenied,
				Message: "权限不足，只有管理员可以删除用户",
			}, nil
		}
	}

	// 查找要删除的用户
	user, err := l.svcCtx.OaUserModel.FindOneByPhone(l.ctx, in.Phone)
	if err != nil {
		if err == model.ErrNotFound {
			l.Infof("用户不存在")
			return &oauser.DeleteUserResp{
				Code:    constants.CodeUserNotFound,
				Message: constants.GetMessage(constants.CodeUserNotFound),
			}, nil
		}
		l.Errorf("查询用户失败: %v", err)
		return &oauser.DeleteUserResp{
			Code:    constants.CodeInternalError,
			Message: constants.GetMessage(constants.CodeInternalError),
		}, nil
	}

	// 执行软删除或硬删除
	// 这里采用软删除的方式，将用户状态设置为禁用
	user.Status = constants.UserStatusDisabled
	err = l.svcCtx.OaUserModel.Update(l.ctx, user)
	if err != nil {
		l.Errorf("删除用户失败: %v", err)
		return &oauser.DeleteUserResp{
			Code:    constants.CodeInternalError,
			Message: constants.GetMessage(constants.CodeInternalError),
		}, nil
	}

	// 如果需要硬删除，可以使用以下代码替换上面的软删除逻辑
	// err = l.svcCtx.OaUserModel.Delete(l.ctx, user.Id)
	// if err != nil {
	//     l.Errorf("删除用户失败: %v", err)
	//     return &oauser.DeleteUserResp{
	//         Code:    constants.CodeInternalError,
	//         Message: constants.GetMessage(constants.CodeInternalError),
	//     }, nil
	// }

	l.Infof("删除后台用户成功, user_id: %d", user.Id)

	return &oauser.DeleteUserResp{
		Code:    constants.CodeSuccess,
		Message: constants.GetMessage(constants.CodeSuccess),
	}, nil
}
