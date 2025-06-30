package logic

import (
	"context"

	"rpc/internal/pkg/constants"
	"rpc/internal/svc"
	"rpc/oauser"

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

func (l *LogoutLogic) Logout(in *oauser.LogoutReq) (*oauser.LogoutResp, error) {
	l.Infof("后台用户注销请求")

	// 参数验证
	if in.Token == "" {
		l.Infof("token参数为空")
		return &oauser.LogoutResp{
			Code:    constants.CodeInvalidParams,
			Message: constants.GetMessage(constants.CodeInvalidParams),
		}, nil
	}

	// 验证 token 的有效性
	claims, err := l.svcCtx.JwtUtils.ValidateAndGetClaims(in.Token)
	if err != nil {
		l.Infof("无效的token")
		return &oauser.LogoutResp{
			Code:    constants.CodeUnauthorized,
			Message: constants.GetMessage(constants.CodeUnauthorized),
		}, nil
	}

	// 在实际项目中，这里可以将 token 加入黑名单
	// 由于这里是简化版本，我们只记录注销操作
	l.Infof("后台用户注销成功, user_id: %d, phone: %s", claims.UserID, claims.Phone)

	return &oauser.LogoutResp{
		Code:    constants.CodeSuccess,
		Message: constants.GetMessage(constants.CodeSuccess),
	}, nil
}
