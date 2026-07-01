package logger

import (
	"context"
	"time"
)

var std Logger

type ctxTraceIDKey struct{}
type ctxLoggerKey struct{}

var kTraceID = ctxTraceIDKey{}
var kContextLogger = ctxLoggerKey{}

const traceIDKey = "trace_id"

type Logger interface {
	Debug(msg string, fields ...Field)
	Info(msg string, fields ...Field)
	Warn(msg string, fields ...Field)
	Error(msg string, fields ...Field)
	With(fields ...Field) Logger
}

func SetLogger(l Logger) {
	std = l
}

func Debug(msg string, fields ...Field) { std.Debug(msg, fields...) }
func Info(msg string, fields ...Field)  { std.Info(msg, fields...) }
func Warn(msg string, fields ...Field)  { std.Warn(msg, fields...) }
func Error(msg string, fields ...Field) { std.Error(msg, fields...) }
func DebugContext(ctx context.Context, msg string, fields ...Field) {
	fromContext(ctx).Debug(msg, fields...)
}
func InfoContext(ctx context.Context, msg string, fields ...Field) {
	fromContext(ctx).Info(msg, fields...)
}
func WarnContext(ctx context.Context, msg string, fields ...Field) {
	fromContext(ctx).Warn(msg, fields...)
}
func ErrorContext(ctx context.Context, msg string, fields ...Field) {
	fromContext(ctx).Error(msg, fields...)
}

func With(fields ...Field) Logger { return std.With(fields...) }

type Field struct {
	Key   string
	Value any
}
type stackMarker struct{}

func Int(key string, val int) Field                { return Field{key, val} }
func Int64(key string, val int64) Field            { return Field{key, val} }
func String(key, val string) Field                 { return Field{key, val} }
func Bool(key string, val bool) Field              { return Field{key, val} }
func Err(err error) Field                          { return Field{"error", err} }
func Duration(key string, val time.Duration) Field { return Field{key, val} }
func Any(key string, val any) Field                { return Field{key, val} }
func Stack(key string) Field                       { return Field{key, stackMarker{}} }

func fromContext(c context.Context) Logger {
	if c == nil {
		return std
	}
	// 先从context获取*zap.logger实例，没有才进行创建逻辑
	logger, ok := c.Value(kContextLogger).(Logger)
	if ok {
		return logger
	}
	traceID, ok := c.Value(kTraceID).(string)
	if ok && traceID != "" {
		return std.With(String(traceIDKey, traceID))
	}
	return std
}

func SetTraceIDCtx(c context.Context, traceID string) context.Context {
	return context.WithValue(c, kTraceID, traceID)
}
func SetLoggerCtx(c context.Context, logger Logger) context.Context {
	return context.WithValue(c, kContextLogger, logger)
}
