package contextutil

import (
	"context"
)

type contextKey string

const (
	traceIDKey     contextKey = "trace_id"
	currentUserKey contextKey = "current_user"
)

func SetTraceID(ctx context.Context, traceID string) context.Context {
	return context.WithValue(ctx, traceIDKey, traceID)
}

func GetTraceID(ctx context.Context) string {
	if ctx == nil {
		return ""
	}
	if traceID, ok := ctx.Value(traceIDKey).(string); ok {
		return traceID
	}
	return ""
}

func SetCurrentUser(ctx context.Context, userID uint) context.Context {
	return context.WithValue(ctx, currentUserKey, userID)
}

func GetCurrentUser(ctx context.Context) uint {
	if ctx == nil {
		return 0
	}
	if id, ok := ctx.Value(currentUserKey).(uint); ok {
		return id
	}
	return 0
}
