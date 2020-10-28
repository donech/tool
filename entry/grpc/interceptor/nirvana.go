package interceptor

import (
	"context"

	"github.com/donech/tool/xtrace"

	"github.com/donech/tool/xlog"

	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

type NirvanaInterceptor struct {
}

func (s *NirvanaInterceptor) Serve(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
	defer func() {
		err1, ok := recover().(error)
		if ok {
			xlog.S(ctx).Errorf("get panic error %v", err1)
			err = err1
		}
	}()
	var traceId string
	// Read metadata from client.
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		// Create and send header.
		traceId = xtrace.NewTraceID()
		ctx = context.WithValue(ctx, xtrace.KeyName, traceId)
		header := metadata.New(map[string]string{xtrace.KeyName: traceId})
		err = grpc.SendHeader(ctx, header)
		if err != nil {
			xlog.S(ctx).Errorf("grpc send header error %v", err)
		}
	} else {
		if t, ok := md[xtrace.KeyName]; ok {
			traceId = t[0]
		} else {
			// Create and send header.
			traceId = xtrace.NewTraceID()
			ctx = context.WithValue(ctx, xtrace.KeyName, traceId)
			header := metadata.New(map[string]string{xtrace.KeyName: traceId})
			err = grpc.SendHeader(ctx, header)
			if err != nil {
				xlog.S(ctx).Errorf("grpc send header error %v", err)
			}
		}
	}
	return handler(ctx, req)
}
