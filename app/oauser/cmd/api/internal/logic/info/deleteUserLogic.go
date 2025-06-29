package info

import (
	"context"
	"net/http"
	"strings"

	"api/internal/svc"
	"api/internal/types"
	"rpc/oauserclient"

	"github.com/zeromicro/go-zero/core/logx"
)

type DeleteUserLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
	r      *http.Request
}

func NewDeleteUserLogic(ctx context.Context, svcCtx *svc.ServiceContext, r *http.Request) *DeleteUserLogic {
	return &DeleteUserLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
		r:      r,
	}
}

func (l *DeleteUserLogic) DeleteUser(req *types.DeleteUserReq) (resp *types.DeleteUserResp, err error) {
	// 从请求头中获取 token
	authHeader := l.r.Header.Get("Authorization")
	if authHeader == "" {
		l.Logger.Error("缺少 Authorization 头")
		return &types.DeleteUserResp{
			Code:    401,
			Message: "缺少身份验证信息",
		}, nil
	}

	// 提取 Bearer token
	token := ""
	if strings.HasPrefix(authHeader, "Bearer ") {
		token = authHeader[7:]
	} else {
		l.Logger.Error("Authorization 头格式错误")
		return &types.DeleteUserResp{
			Code:    401,
			Message: "身份验证信息格式错误",
		}, nil
	}

	// 调用 RPC 服务删除用户，传递调用者 token
	deleteResp, err := l.svcCtx.OaUserRpc.DeleteUser(l.ctx, &oauserclient.DeleteUserReq{
		Phone:       req.Phone,
		CallerToken: token,
	})
	if err != nil {
		l.Logger.Errorf("RPC DeleteUser failed: %v", err)
		return &types.DeleteUserResp{
			Code:    500,
			Message: "服务器内部错误",
		}, nil
	}

	// 转换响应格式
	return &types.DeleteUserResp{
		Code:    deleteResp.Code,
		Message: deleteResp.Message,
	}, nil
}
