package constants

// 错误码定义
const (
	// 成功
	CodeSuccess = 200

	// 客户端错误 4xx
	CodeInvalidParams    = 400   // 参数错误
	CodeUnauthorized     = 401   // 未授权
	CodeForbidden        = 403   // 禁止访问
	CodeNotFound         = 404   // 资源不存在
	CodePhoneInvalid     = 40001 // 手机号格式无效
	CodePasswordError    = 40002 // 密码错误
	CodeUserNotFound     = 40003 // 用户不存在
	CodeUserExists       = 40004 // 用户已存在
	CodeUserDisabled     = 40005 // 用户被禁用
	CodeTokenInvalid     = 40006 // Token无效
	CodeTokenExpired     = 40007 // Token过期
	CodePermissionDenied = 40008 // 权限不足

	// 服务器错误 5xx
	CodeInternalError = 500   // 服务器内部错误
	CodeDatabaseError = 50001 // 数据库错误
	CodeCacheError    = 50002 // 缓存错误
)

// 错误消息映射
var messageMap = map[int32]string{
	CodeSuccess:          "操作成功",
	CodeInvalidParams:    "参数错误",
	CodeUnauthorized:     "未授权访问",
	CodeForbidden:        "禁止访问",
	CodeNotFound:         "资源不存在",
	CodePhoneInvalid:     "手机号格式无效",
	CodePasswordError:    "密码错误",
	CodeUserNotFound:     "用户不存在",
	CodeUserExists:       "用户已存在",
	CodeUserDisabled:     "用户被禁用",
	CodeTokenInvalid:     "Token无效",
	CodeTokenExpired:     "Token已过期",
	CodePermissionDenied: "权限不足",
	CodeInternalError:    "服务器内部错误",
	CodeDatabaseError:    "数据库操作失败",
	CodeCacheError:       "缓存操作失败",
}

// GetMessage 根据错误码获取错误消息
func GetMessage(code int32) string {
	if msg, ok := messageMap[code]; ok {
		return msg
	}
	return "未知错误"
}

// 用户状态常量
const (
	UserStatusNormal   = 1 // 正常
	UserStatusDisabled = 2 // 禁用
)

// 用户角色常量
const (
	RoleAdmin    = "admin"    // 管理员
	RoleOperator = "operator" // 普通操作员
)

// 性别常量
const (
	GenderUnknown = 0 // 未知
	GenderMale    = 1 // 男
	GenderFemale  = 2 // 女
)
