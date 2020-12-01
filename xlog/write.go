package xlog

import (
	"context"

	"github.com/donech/tool/xtrace"

	"go.uber.org/zap"
)

//S SugaredLogger with xtrace-id field
func S(ctx context.Context) *zap.SugaredLogger {
	return zap.S().With(TraceIDField(ctx))
}

//Logger with xtrace-id field
func L(ctx context.Context) *zap.Logger {
	return zap.L().With(TraceIDField(ctx))
}

//TraceIDField TraceIDField
func TraceIDField(ctx context.Context) zap.Field {
	traceID := xtrace.GetTraceIDFromContext(ctx)
	if traceID != "" {
		return zap.String(string(xtrace.KeyName), traceID)
	}
	return zap.Skip()
}
