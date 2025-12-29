package contextutil

import (
	"context"

	"go.uber.org/zap"
)

type contextKey string

const (
	traceIDKey       contextKey = "trace_id"
	currentUserIDKey contextKey = "current_user"
	contextLoggerKey contextKey = "context_logger"
)

func SetTraceID(c context.Context, traceID string) context.Context {
	return context.WithValue(c, traceIDKey, traceID)
}
func GetTraceID(c context.Context) string {
	if c == nil {
		return ""
	}
	if traceID, ok := c.Value(traceIDKey).(string); ok {
		return traceID
	}
	return ""
}

func SetCurrentUser(c context.Context, userID uint) context.Context {
	return context.WithValue(c, currentUserIDKey, userID)
}
func GetCurrentUser(c context.Context) uint {
	if c == nil {
		return 0
	}
	if id, ok := c.Value(currentUserIDKey).(uint); ok {
		return id
	}
	return 0
}

func SetContextLoggerKey(c context.Context, logger *zap.Logger) context.Context {
	return context.WithValue(c, contextLoggerKey, logger)
}
func GetContextLoggerKey(c context.Context) *zap.Logger {
	if c == nil {
		return nil
	}
	if logger, ok := c.Value(contextLoggerKey).(*zap.Logger); ok {
		return logger
	}
	return nil
}
