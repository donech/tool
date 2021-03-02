package xgrpc

import (
	"context"
	"net"

	"github.com/donech/tool/xjwt"
	"github.com/donech/tool/xlog"

	"google.golang.org/grpc/reflection"

	"github.com/donech/tool/entry/xgrpc/interceptor"

	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	grpcopentracing "github.com/grpc-ecosystem/go-grpc-middleware/tracing/opentracing"
	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	"github.com/opentracing/opentracing-go"

	"google.golang.org/grpc"
)

func New(config Config, options ...Option) *Entry {
	en := &Entry{
		config: config,
	}
	for _, option := range options {
		option(en)
	}
	return en
}

type Option func(entry *Entry)

func WithRegisteServer(server RegisteServer) Option {
	return func(entry *Entry) {
		entry.registeServer = server
	}
}

func WithRegisteWebHandler(handler RegisteWebHandler) Option {
	return func(entry *Entry) {
		entry.registeWebHandler = handler
	}
}

func WithJwtFactory(jwtFactory *xjwt.JWTFactory) Option {
	return func(entry *Entry) {
		entry.jwtFactory = jwtFactory
	}
}

type RegisteServer func(server *grpc.Server)
type RegisteWebHandler func(ctx context.Context, mux *runtime.ServeMux, conn *grpc.ClientConn) error

type Entry struct {
	config            Config
	srv               *grpc.Server
	registeServer     RegisteServer
	registeWebHandler RegisteWebHandler
	jwtFactory        *xjwt.JWTFactory
}

func (e *Entry) Run() error {
	traceIdInterceptor := interceptor.TraceIdInterceptor{}
	logInterceptor := interceptor.LogInterceptor{}
	JwtInterceptor := interceptor.NewJwtInterceptor(e.jwtFactory)
	srv := grpc.NewServer(
		grpc.UnaryInterceptor(grpc_middleware.ChainUnaryServer(
			traceIdInterceptor.Serve,
			grpcopentracing.UnaryServerInterceptor(grpcopentracing.WithTracer(opentracing.GlobalTracer())),
			logInterceptor.Serve,
			JwtInterceptor.Serve,
		)))
	if e.registeServer != nil {
		e.registeServer(srv)
		e.srv = srv
	}
	if e.config.EnableReflect {
		xlog.SS().Info("enable grpc reflect")
		reflection.Register(srv)
	}
	listen, err := net.Listen("tcp", e.config.Port)
	if err != nil {
		xlog.SS().Errorf("listen tcp error: %s", err)
		return err
	}
	xlog.SS().Infof("listening tcp port: %s", e.config.Port)

	go func() {
		xlog.SS().Info("start grpc server at ", e.config.Port)
		if err = srv.Serve(listen); err != nil {
			xlog.SS().Fatalf("grpc serve listen error: %s", err)
		}
	}()
	return nil
}

func (e *Entry) Stop(ctx context.Context) error {
	e.srv.GracefulStop()
	return nil
}
