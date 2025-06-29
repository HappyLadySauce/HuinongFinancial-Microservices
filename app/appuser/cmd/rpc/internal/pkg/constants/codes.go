package constants

// 响应代码常量
const (
	// 成功
	CodeSuccess = 200

	// 通用错误
	CodeInvalidParams      = 400 // 参数错误
	CodeUnauthorized       = 401 // 未授权
	CodeForbidden          = 403 // 禁止访问
	CodeNotFound           = 404 // 未找到
	CodeInternalError      = 500 // 内部错误
	CodeServiceUnavailable = 503 // 服务不可用

	// 用户相关错误码
	CodeUserNotFound       = 1001 // 用户不存在
	CodeUserExists         = 1002 // 用户已存在
	CodePasswordError      = 1003 // 密码错误
	CodeUserDisabled       = 1004 // 用户被禁用
	CodeUserFrozen         = 1005 // 用户被冻结
	CodePhoneInvalid       = 1006 // 手机号格式无效
	CodeTokenInvalid       = 1007 // Token 无效
	CodeTokenExpired       = 1008 // Token 过期
	CodeUserAlreadyDeleted = 1009 // 用户已被删除

	// 权限相关错误码
	CodeRoleInvalid      = 2001 // 角色无效
	CodePermissionDenied = 2002 // 权限不足
)

// 错误消息映射
var CodeMessages = map[int32]string{
	CodeSuccess:            "success",
	CodeInvalidParams:      "参数错误",
	CodeUnauthorized:       "未授权",
	CodeForbidden:          "禁止访问",
	CodeNotFound:           "未找到",
	CodeInternalError:      "内部错误",
	CodeServiceUnavailable: "服务不可用",

	CodeUserNotFound:       "用户不存在",
	CodeUserExists:         "用户已存在",
	CodePasswordError:      "密码错误",
	CodeUserDisabled:       "用户被禁用",
	CodeUserFrozen:         "用户被冻结",
	CodePhoneInvalid:       "手机号格式无效",
	CodeTokenInvalid:       "Token 无效",
	CodeTokenExpired:       "Token 过期",
	CodeUserAlreadyDeleted: "用户已被删除",

	CodeRoleInvalid:      "角色无效",
	CodePermissionDenied: "权限不足",
}

// GetMessage 获取错误码对应的消息
func GetMessage(code int32) string {
	if msg, ok := CodeMessages[code]; ok {
		return msg
	}
	return "未知错误"
}
