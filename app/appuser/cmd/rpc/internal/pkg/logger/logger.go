package logger

import (
	"context"
	"os"
	"path/filepath"

	"github.com/sirupsen/logrus"
	"gopkg.in/natefinch/lumberjack.v2"
)

var (
	// GlobalLogger 全局日志实例
	GlobalLogger *logrus.Logger
)

// Config 日志配置
type Config struct {
	ServiceName string `json:"serviceName"`
	Mode        string `json:"mode"`       // file, console
	Path        string `json:"path"`       // 日志文件路径
	Level       string `json:"level"`      // debug, info, warn, error
	KeepDays    int    `json:"keepDays"`   // 日志保留天数
	Compress    bool   `json:"compress"`   // 是否压缩
	MaxSize     int    `json:"maxSize"`    // 单个文件最大尺寸（MB）
	MaxBackups  int    `json:"maxBackups"` // 保留文件数量
}

// InitLogger 初始化日志
func InitLogger(cfg Config) {
	GlobalLogger = logrus.New()

	// 设置日志级别
	level, err := logrus.ParseLevel(cfg.Level)
	if err != nil {
		level = logrus.InfoLevel
	}
	GlobalLogger.SetLevel(level)

	// 设置日志格式
	GlobalLogger.SetFormatter(&logrus.JSONFormatter{
		TimestampFormat: "2006-01-02 15:04:05.000",
		FieldMap: logrus.FieldMap{
			logrus.FieldKeyTime:  "timestamp",
			logrus.FieldKeyLevel: "level",
			logrus.FieldKeyMsg:   "message",
		},
	})

	// 添加服务名称字段
	GlobalLogger = GlobalLogger.WithField("service", cfg.ServiceName).Logger

	// 配置输出
	if cfg.Mode == "file" {
		// 确保日志目录存在
		if cfg.Path != "" {
			os.MkdirAll(cfg.Path, 0755)
		}

		// 文件输出配置
		logFile := &lumberjack.Logger{
			Filename:   filepath.Join(cfg.Path, cfg.ServiceName+".log"),
			MaxSize:    cfg.MaxSize,
			MaxBackups: cfg.MaxBackups,
			MaxAge:     cfg.KeepDays,
			Compress:   cfg.Compress,
			LocalTime:  true,
		}
		GlobalLogger.SetOutput(logFile)
	} else {
		// 控制台输出
		GlobalLogger.SetOutput(os.Stdout)
	}
}

// WithContext 创建带上下文的日志实例
func WithContext(ctx context.Context) *logrus.Entry {
	entry := GlobalLogger.WithContext(ctx)

	// 如果有 SkyWalking trace 信息，添加到日志中
	if traceID := getTraceID(ctx); traceID != "" {
		entry = entry.WithField("traceId", traceID)
	}

	return entry
}

// WithFields 创建带字段的日志实例
func WithFields(fields logrus.Fields) *logrus.Entry {
	return GlobalLogger.WithFields(fields)
}

// getTraceID 从上下文中获取 SkyWalking TraceID
func getTraceID(ctx context.Context) string {
	// SkyWalking Go Agent 会在 context 中设置 trace 信息
	// 这里预留接口，后续集成时完善
	return ""
}
