package interceptor

import (
	"context"

	"go.uber.org/zap"

	"github.com/donech/tool/xlog"

	"google.golang.org/grpc"
)

type LogInterceptor struct{}

func (i *LogInterceptor) Serve(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
	xlog.L(ctx).Info("incoming grpc req", zap.Reflect("req", req))
	resp, err = handler(ctx, req)
	xlog.L(ctx).Info("output grpc resp", zap.Reflect("resp", resp), zap.Error(err))
	return resp, err
}
