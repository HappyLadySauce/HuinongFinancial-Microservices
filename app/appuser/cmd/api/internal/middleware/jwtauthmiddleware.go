package middleware

import (
	"net/http"
)

type JwtAuthMiddleware struct {
}

func NewJwtAuthMiddleware() *JwtAuthMiddleware {
	return &JwtAuthMiddleware{}
}

func (m *JwtAuthMiddleware) Handle(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// go-zero的JWT认证由框架自动处理
		// 通过rest.WithJwt()配置后，框架会自动：
		// 1. 验证JWT token的有效性
		// 2. 解析JWT中的claims数据
		// 3. 将解析后的数据注入到request context中
		// 4. 如果JWT验证失败，会自动返回401错误
		// 5. 如果验证成功，用户信息会自动注入到context中，可以通过ctx.Value()获取

		// 这里可以添加额外的业务逻辑，比如：
		// - 记录访问日志
		// - 检查用户权限
		// - 其他自定义验证逻辑

		next(w, r)
	}
}
