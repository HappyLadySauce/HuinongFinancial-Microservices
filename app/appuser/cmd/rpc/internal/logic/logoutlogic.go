package logic

import (
	"context"

	"rpc/appuser"
	"rpc/internal/pkg/constants"
	"rpc/internal/pkg/utils"
	"rpc/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
)

type LogoutLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewLogoutLogic(ctx context.Context, svcCtx *svc.ServiceContext) *LogoutLogic {
	return &LogoutLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 用户登出
func (l *LogoutLogic) Logout(in *appuser.LogoutReq) (*appuser.LogoutResp, error) {
	l.Info("用户登出请求")

	// 参数验证
	if in.Token == "" {
		l.Infof("Token不能为空")
		return &appuser.LogoutResp{
			Code:    constants.CodeInvalidParams,
			Message: constants.GetMessage(constants.CodeInvalidParams),
		}, nil
	}

	// 验证 Token 有效性
	claims, err := utils.ParseToken(in.Token, l.svcCtx.Config.JwtAuth.AccessSecret)
	if err != nil {
		l.Errorf("Token解析失败: %v", err)
		return &appuser.LogoutResp{
			Code:    constants.CodeTokenInvalid,
			Message: constants.GetMessage(constants.CodeTokenInvalid),
		}, nil
	}

	// 获取用户类型
	userType, ok := claims["user_type"].(string)
	if !ok {
		l.Infof("Token中缺少用户类型信息")
		return &appuser.LogoutResp{
			Code:    constants.CodeTokenInvalid,
			Message: constants.GetMessage(constants.CodeTokenInvalid),
		}, nil
	}

	// 检查用户类型，确保是 app 用户
	if userType != "app" {
		l.Infof("用户类型不匹配, user_type: %s", userType)
		return &appuser.LogoutResp{
			Code:    constants.CodeTokenInvalid,
			Message: constants.GetMessage(constants.CodeTokenInvalid),
		}, nil
	}

	// 获取用户ID和手机号
	userID, _ := claims["user_id"].(float64) // JWT中数字通常是float64
	phone, _ := claims["phone"].(string)

	// 记录登出操作
	l.Infof("用户登出成功, user_id: %.0f, phone: %s", userID, phone)

	// 在真实场景中，可以考虑以下操作：
	// 1. 将 token 添加到黑名单 (Redis)
	// 2. 更新用户最后登出时间
	// 3. 清理用户会话相关缓存

	// 目前简单实现，仅记录日志
	return &appuser.LogoutResp{
		Code:    constants.CodeSuccess,
		Message: constants.GetMessage(constants.CodeSuccess),
	}, nil
}
