package xlog

import (
	"context"

	"github.com/donech/tool/xtrace"

	"go.uber.org/zap"
)

var SystemKey = "system"
var SystemName = "system"
var SystemTraceName = "system"

//S SugaredLogger with xtrace-id field
func S(ctx context.Context) *zap.SugaredLogger {
	return zap.S().With(TraceIDField(ctx), SystemTraceIDField())
}

//Logger with xtrace-id field
func L(ctx context.Context) *zap.Logger {
	return zap.L().With(TraceIDField(ctx), SystemTraceIDField())
}

//SS SugaredLogger with traceID = system
func SS() *zap.SugaredLogger {
	return zap.S().With(zap.String(string(xtrace.KeyName), SystemName), SystemTraceIDField())
}

//Logger with with traceID = system
func SL() *zap.Logger {
	return zap.L().With(zap.String(string(xtrace.KeyName), SystemName), SystemTraceIDField())
}

//TraceIDField TraceIDField
func TraceIDField(ctx context.Context) zap.Field {
	traceID := xtrace.GetTraceIDFromContext(ctx)
	if traceID != "" {
		return zap.String(string(xtrace.KeyName), traceID)
	}
	return zap.Skip()
}

//SystemTraceIDField SystemTraceIDField
func SystemTraceIDField() zap.Field {
	return zap.String(SystemKey, SystemTraceName)
}
