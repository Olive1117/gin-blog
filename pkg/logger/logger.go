package logger

import (
	"context"
	"os"
	"strings"

	"github.com/natefinch/lumberjack"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type traceIDKey struct{}
type contextLoggerKey struct{}

var kTraceID = traceIDKey{}
var kContextLogger = contextLoggerKey{}
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
	core := zapcore.NewTee(
		zapcore.NewCore(NewJSONEncoder(), writer, getLevel(level)),
		zapcore.NewCore(NewConsoleEncoder(), zapcore.AddSync(os.Stdout), getLevel(level)),
	)
	L = zap.New(core, zap.AddCaller())
	zap.RedirectStdLog(L)
}
func NewJSONEncoder() zapcore.Encoder {
	encoderConfig := zap.NewProductionEncoderConfig()
	encoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	return zapcore.NewJSONEncoder(encoderConfig)
}
func NewConsoleEncoder() zapcore.Encoder {
	encoderConfig := zap.NewDevelopmentEncoderConfig()
	encoderConfig.EncodeTime = zapcore.TimeEncoderOfLayout("2006-01-02 15:04:05")
	encoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
	encoderConfig.EncodeCaller = zapcore.ShortCallerEncoder
	return zapcore.NewConsoleEncoder(encoderConfig)
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
	logger, ok := c.Value(kContextLogger).(*zap.Logger)
	if ok {
		return logger
	}
	traceID, ok := c.Value(kTraceID).(string)
	if ok && traceID != "" {
		return L.With(zap.String(TraceIDKey, traceID))
	}
	return L
}

func SetTraceID(c context.Context, traceID string) context.Context {
	return context.WithValue(c, kTraceID, traceID)
}
func SetCurrentUser(c context.Context, logger *zap.Logger) context.Context {
	return context.WithValue(c, kContextLogger, logger)
}
