package interceptor

import (
	"context"

	"github.com/pkg/errors"

	"github.com/donech/tool/xtrace"

	"github.com/donech/tool/xlog"

	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

type TraceIdInterceptor struct{}

func (s *TraceIdInterceptor) Serve(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
	defer func() {
		err1, ok := recover().(error)
		if ok {
			xlog.S(ctx).Errorf("get panic error %+v", errors.WithStack(err1))
			err = err1
		}
	}()

	var traceId string
	// Read metadata from client.
	md, ok := metadata.FromIncomingContext(ctx)
	if ok {
		xlog.S(ctx).Debugf("incoming md is: %v", md)
		traceId = xtrace.GetTraceIDFromGrpcMetadata(md)
		if traceId != "" {
			xlog.S(ctx).Infof("grpc incoming header with trace id %s", traceId)
		}
	}
	if traceId == "" {
		xlog.S(ctx).Debug("trace-id not found, so generate it")
		traceId = xtrace.NewTraceID()
	}
	ctx = context.WithValue(ctx, xtrace.KeyName, traceId)
	header := metadata.New(map[string]string{string(xtrace.KeyName): traceId})
	err = grpc.SendHeader(ctx, header)
	if err != nil {
		xlog.S(ctx).Errorf("grpc send header error %v", err)
	}
	return handler(ctx, req)
}
