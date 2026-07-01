package logger

import (
	"os"
	"strings"
	"time"

	"github.com/natefinch/lumberjack"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type zapLogger struct {
	zap *zap.Logger
}

func NewZapLogger(logPath string, level string) Logger {
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
	return &zapLogger{zap: zap.New(core, zap.AddCaller(), zap.AddCallerSkip(2))}
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
	encoderConfig.LineEnding = zapcore.DefaultLineEnding
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

func toZapFields(field Field) zap.Field {
	switch v := field.Value.(type) {
	case int:
		return zap.Int(field.Key, v)
	case int64:
		return zap.Int64(field.Key, v)
	case string:
		return zap.String(field.Key, v)
	case bool:
		return zap.Bool(field.Key, v)
	case error:
		return zap.Error(v)
	case time.Duration:
		return zap.Duration(field.Key, v)
	case stackMarker:
		return zap.Stack(field.Key)
	default:
		return zap.Any(field.Key, v)
	}
}
func toZapFieldsSlice(fields []Field) []zap.Field {
	zapFields := make([]zap.Field, 0, len(fields))
	for _, field := range fields {
		zapFields = append(zapFields, toZapFields(field))
	}
	return zapFields
}

func (z *zapLogger) Debug(msg string, fields ...Field) {
	z.zap.Debug(msg, toZapFieldsSlice(fields)...)
}
func (z *zapLogger) Info(msg string, fields ...Field) {
	z.zap.Info(msg, toZapFieldsSlice(fields)...)
}
func (z *zapLogger) Warn(msg string, fields ...Field) {
	z.zap.Warn(msg, toZapFieldsSlice(fields)...)
}
func (z *zapLogger) Error(msg string, fields ...Field) {
	z.zap.Error(msg, toZapFieldsSlice(fields)...)
}
func (z *zapLogger) With(fields ...Field) Logger {
	return &zapLogger{zap: z.zap.With(toZapFieldsSlice(fields)...)}
}
