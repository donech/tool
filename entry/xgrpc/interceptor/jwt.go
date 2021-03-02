package interceptor

import (
	"context"

	"github.com/donech/tool/xlog"

	"github.com/grpc-ecosystem/go-grpc-middleware/util/metautils"

	"github.com/donech/tool/xjwt"

	"google.golang.org/grpc"
)

var (
	headerAuthorize = "authorization"
)

type JwtInterceptor struct {
	jwtFactory *xjwt.JWTFactory
}

func NewJwtInterceptor(jwtFactory *xjwt.JWTFactory) *JwtInterceptor {
	return &JwtInterceptor{jwtFactory: jwtFactory}
}

func (i *JwtInterceptor) Serve(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
	token := metautils.ExtractIncoming(ctx).Get(headerAuthorize)
	if token != "" {
		claims, err := i.jwtFactory.GetClaims(token)
		if err != nil {
			xlog.S(ctx).Warnf("jwt GetClaims error, %+v", err)
		}
		newCtx := xjwt.SetClaimsToCtx(ctx, claims)
		return handler(newCtx, req)
	}
	return handler(ctx, req)
}
