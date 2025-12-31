package logger

import (
	"context"
	"os"
	"strings"

	"github.com/Olive1117/gin-blog/pkg/contextutil"
	"github.com/natefinch/lumberjack"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var L *zap.Logger

const TraceIDKey = "trace_id"

func NewLogger(logPath string, level string) {
	writer := zapcore.AddSync(&lumberjack.Logger{
		Filename:   logPath,
		MaxSize:    100,
		MaxBackups: 5,
		MaxAge:     30,
		Compress:   true,
	})
	encoderConfig := zap.NewProductionEncoderConfig()
	encoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	JSONEncoder := zapcore.NewJSONEncoder(encoderConfig)
	encoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
	consoleEncoder := zapcore.NewConsoleEncoder(encoderConfig)

	core := zapcore.NewTee(
		zapcore.NewCore(JSONEncoder, writer, getLevel(level)),
		zapcore.NewCore(consoleEncoder, zapcore.AddSync(os.Stdout), getLevel(level)),
	)
	L = zap.New(core, zap.AddCaller())
	zap.RedirectStdLog(L)
}

func getLevel(level string) zapcore.Level {
	switch strings.ToLower(level) {
	case "debug":
		return zapcore.DebugLevel
	case "info":
		return zapcore.InfoLevel
	case "warn":
		return zapcore.WarnLevel
	case "error":
		return zapcore.ErrorLevel
	case "dpanic":
		return zapcore.DPanicLevel
	case "panic":
		return zapcore.PanicLevel
	case "fatal":
		return zapcore.FatalLevel
	default:
		return zapcore.InfoLevel // 默认级别
	}
}

// FromContext 从 c 中提取 trace_id 并返回一个带字段的 zap.Logger
func FromContext(c context.Context) *zap.Logger {
	if c == nil {
		return L
	}
	// 先从context获取*zap.logger实例，没有才进行创建逻辑
	logger := contextutil.GetContextLoggerKey(c)
	if logger != nil {
		return logger
	}
	traceID := contextutil.GetTraceID(c)
	if traceID != "" {
		return L.With(zap.String(TraceIDKey, traceID))
	}
	return L
}

// func Debug(ctx context.Context, msg string, fields ...zap.Field) {
// 	FromContext(ctx).Debug(msg, fields...)
// }

// func Info(ctx context.Context, msg string, fields ...zap.Field) {
// 	FromContext(ctx).Info(msg, fields...)
// }

// func Error(ctx context.Context, msg string, fields ...zap.Field) {
// 	FromContext(ctx).Error(msg, fields...)
// }

// func Warn(ctx context.Context, msg string, fields ...zap.Field) {
// 	FromContext(ctx).Warn(msg, fields...)
// }
