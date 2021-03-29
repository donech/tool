package interceptor

import (
	"context"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/donech/tool/xlog"

	"github.com/grpc-ecosystem/go-grpc-middleware/util/metautils"

	"github.com/donech/tool/xjwt"

	"google.golang.org/grpc"
)

var (
	headerAuthorize = "authorization"
)

type JwtInterceptor struct {
	jwtFactory  *xjwt.JWTFactory
	jumpMethods map[string]bool
}

func NewJwtInterceptor(jwtFactory *xjwt.JWTFactory, jumps map[string]bool) *JwtInterceptor {
	return &JwtInterceptor{jwtFactory: jwtFactory, jumpMethods: jumps}
}

func (i *JwtInterceptor) Serve(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
	if i.jumpMethods == nil {
		return handler(ctx, req)
	}
	if i.jumpMethods[info.FullMethod] {
		return handler(ctx, req)
	}
	token := metautils.ExtractIncoming(ctx).Get(headerAuthorize)
	if token == "" {
		xlog.S(ctx).Error("no token found")
		return nil, status.Error(codes.Unauthenticated, "Unauthorized")
	}
	claims, err := i.jwtFactory.GetClaims(token)
	if err != nil {
		xlog.S(ctx).Error("jwt GetClaims error, %+v", err)
		return nil, status.Error(codes.Unauthenticated, "Unauthorized")
	}
	newCtx := xjwt.SetClaimsToCtx(ctx, claims)
	return handler(newCtx, req)

}
