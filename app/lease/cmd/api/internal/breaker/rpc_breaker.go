package breaker

import (
	"context"
	"errors"
	"strings"
	"time"

	"github.com/zeromicro/go-zero/core/breaker"
	"github.com/zeromicro/go-zero/core/logx"
)

// RPC调用熔断器封装
type RpcBreakerClient struct {
	serviceName string
	logger      logx.Logger
}

// 创建RPC熔断器客户端
func NewRpcBreakerClient(serviceName string) *RpcBreakerClient {
	return &RpcBreakerClient{
		serviceName: serviceName,
		logger:      logx.WithContext(context.Background()),
	}
}

// 使用熔断器执行RPC调用 - 默认错误判断
func (r *RpcBreakerClient) DoWithBreaker(ctx context.Context, fn func() error) error {
	return breaker.DoWithFallback(r.serviceName, func() error {
		return fn()
	}, func(err error) error {
		r.logger.WithContext(ctx).Errorf("RPC调用失败，熔断器生效，服务: %s, 错误: %v", r.serviceName, err)
		return errors.New("服务暂时不可用，请稍后重试")
	})
}

// 使用熔断器执行RPC调用 - 自定义错误判断
func (r *RpcBreakerClient) DoWithBreakerAcceptable(ctx context.Context, fn func() error, acceptable func(err error) bool) error {
	return breaker.DoWithFallbackAcceptable(r.serviceName, func() error {
		return fn()
	}, func(err error) error {
		r.logger.WithContext(ctx).Errorf("RPC调用失败，熔断器生效，服务: %s, 错误: %v", r.serviceName, err)
		return errors.New("服务暂时不可用，请稍后重试")
	}, acceptable)
}

// 使用熔断器执行RPC调用并返回结果 - 泛型版本
func DoWithBreakerResult[T any](ctx context.Context, serviceName string, fn func() (T, error)) (T, error) {
	var result T
	var resultErr error

	err := breaker.DoWithFallback(serviceName, func() error {
		res, err := fn()
		if err != nil {
			return err
		}
		result = res
		return nil
	}, func(err error) error {
		logx.WithContext(ctx).Errorf("RPC调用失败，熔断器生效，服务: %s, 错误: %v", serviceName, err)
		resultErr = errors.New("服务暂时不可用，请稍后重试")
		return resultErr
	})

	if err != nil {
		return result, err
	}
	if resultErr != nil {
		return result, resultErr
	}

	return result, nil
}

// 使用熔断器执行RPC调用并返回结果 - 自定义错误判断版本
func DoWithBreakerResultAcceptable[T any](ctx context.Context, serviceName string, fn func() (T, error), acceptable func(err error) bool) (T, error) {
	var result T
	var resultErr error

	err := breaker.DoWithFallbackAcceptable(serviceName, func() error {
		res, err := fn()
		if err != nil {
			return err
		}
		result = res
		return nil
	}, func(err error) error {
		logx.WithContext(ctx).Errorf("RPC调用失败，熔断器生效，服务: %s, 错误: %v", serviceName, err)
		resultErr = errors.New("服务暂时不可用，请稍后重试")
		return resultErr
	}, acceptable)

	if err != nil {
		return result, err
	}
	if resultErr != nil {
		return result, resultErr
	}

	return result, nil
}

// 带重试的熔断器调用 - 用于有返回值的RPC调用
func DoWithBreakerAndRetry[T any](ctx context.Context, serviceName string, fn func() (T, error), maxRetries int, retryDelay time.Duration) (T, error) {
	var result T
	var lastErr error

	for attempt := 0; attempt <= maxRetries; attempt++ {
		if attempt > 0 {
			// 等待重试间隔
			select {
			case <-ctx.Done():
				return result, ctx.Err()
			case <-time.After(retryDelay):
			}
			logx.WithContext(ctx).Infof("重试第%d次调用服务: %s", attempt, serviceName)
		}

		result, lastErr = DoWithBreakerResultAcceptable(ctx, serviceName, fn, IsAcceptableError)
		if lastErr == nil {
			return result, nil
		}

		// 如果是业务错误，不重试
		if IsAcceptableError(lastErr) {
			return result, lastErr
		}

		logx.WithContext(ctx).Error("调用服务%s失败，错误: %v", serviceName, lastErr)
	}

	logx.WithContext(ctx).Errorf("调用服务%s重试%d次后仍失败，最后错误: %v", serviceName, maxRetries, lastErr)
	return result, lastErr
}

// 判断是否为可接受的错误（不计入熔断统计）
func IsAcceptableError(err error) bool {
	if err == nil {
		return true
	}

	errMsg := err.Error()

	// 业务错误关键词 - 不触发熔断，也不重试
	businessErrors := []string{
		"参数错误",
		"权限不足",
		"用户不存在",
		"产品不存在",
		"申请不存在",
		"状态错误",
		"重复提交",
		"余额不足",
		"库存不足",
		"密码错误",
		"用户已存在",
		"手机号已注册",
		"验证码错误",
		"token无效",
		"未找到",
		"已存在",
		"无权限",
		"数据不存在",
	}

	// 系统错误关键词 - 触发熔断，可以重试
	systemErrors := []string{
		"连接超时",
		"网络错误",
		"服务不可用",
		"连接失败",
		"超时",
		"网络异常",
		"服务异常",
		"系统异常",
		"internal error",
		"connection",
		"timeout",
		"unavailable",
	}

	// 先检查是否为业务错误
	for _, keyword := range businessErrors {
		if strings.Contains(strings.ToLower(errMsg), strings.ToLower(keyword)) {
			return true // 业务错误：不计入熔断统计，不重试
		}
	}

	// 检查是否为明确的系统错误
	for _, keyword := range systemErrors {
		if strings.Contains(strings.ToLower(errMsg), strings.ToLower(keyword)) {
			return false // 系统错误：计入熔断统计，可以重试
		}
	}

	// 对于未分类的错误，默认作为系统错误处理
	return false
}

// 错误类型检查辅助函数
func IsBusinessError(err error) bool {
	return IsAcceptableError(err)
}

func IsSystemError(err error) bool {
	return !IsAcceptableError(err)
}

// 创建友好的错误信息
func CreateFriendlyErrorMessage(serviceName string, err error) error {
	if err == nil {
		return nil
	}

	if IsBusinessError(err) {
		// 业务错误直接返回原始错误信息
		return err
	}

	// 系统错误返回友好提示
	switch serviceName {
	case "appuser-rpc":
		return errors.New("用户服务暂时不可用，请稍后重试")
	case "oauser-rpc":
		return errors.New("管理员服务暂时不可用，请稍后重试")
	case "loan-rpc":
		return errors.New("贷款服务暂时不可用，请稍后重试")
	case "lease-rpc":
		return errors.New("租赁服务暂时不可用，请稍后重试")
	case "loanproduct-rpc":
		return errors.New("贷款产品服务暂时不可用，请稍后重试")
	case "leaseproduct-rpc":
		return errors.New("租赁产品服务暂时不可用，请稍后重试")
	default:
		return errors.New("服务暂时不可用，请稍后重试")
	}
}
