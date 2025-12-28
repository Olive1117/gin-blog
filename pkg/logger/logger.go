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

var Log *zap.Logger

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
	Log = zap.New(core, zap.AddCaller())
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

// FromContext 从 ctx 中提取 trace_id 并返回一个带字段的 zap.Logger
func FromContext(ctx context.Context) *zap.Logger {
	if ctx == nil {
		return Log
	}
	traceID := contextutil.GetTraceID(ctx)
	if traceID != "" {
		return Log.With(zap.String(TraceIDKey, traceID))
	}
	return Log
}
