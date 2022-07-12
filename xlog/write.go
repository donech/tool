package xlog

import (
	"context"

	"github.com/donech/tool/xtrace"

	"go.uber.org/zap"
)

var serviceKey = "service"
var serviceName = "service"
var internalTraceId = "internal"

//S SugaredLogger with xtrace-id field
func S(ctx context.Context) *zap.SugaredLogger {
	return zap.S().With(TraceIDField(ctx), ServiceField())
}

//Logger with xtrace-id field
func L(ctx context.Context) *zap.Logger {
	return zap.L().With(TraceIDField(ctx), ServiceField())
}

//SS SugaredLogger with traceID = system
func SS() *zap.SugaredLogger {
	return zap.S().With(zap.String(string(xtrace.KeyName), internalTraceId), ServiceField())
}

//Logger with with traceID = system
func SL() *zap.Logger {
	return zap.L().With(zap.String(string(xtrace.KeyName), internalTraceId), ServiceField())
}

//TraceIDField TraceIDField
func TraceIDField(ctx context.Context) zap.Field {
	traceID := xtrace.GetTraceIDFromContext(ctx)
	if traceID != "" {
		return zap.String(string(xtrace.KeyName), traceID)
	}
	return zap.Skip()
}

//ServiceField ServiceField
func ServiceField() zap.Field {
	return zap.String(serviceKey, serviceName)
}
