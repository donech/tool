package grpc

import (
	"context"
	"log"
	"net"

	"github.com/donech/tool/xlog"

	"google.golang.org/grpc/reflection"

	"github.com/donech/tool/entry/grpc/interceptor"

	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	grpcopentracing "github.com/grpc-ecosystem/go-grpc-middleware/tracing/opentracing"
	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	"github.com/opentracing/opentracing-go"

	"google.golang.org/grpc"

	"go.uber.org/zap"
)

func New(config Config, logger *zap.Logger, server RegisteServer) *Entry {
	return &Entry{
		config:        config,
		logger:        logger,
		registeServer: server,
	}
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

type RegisteServer func(server *grpc.Server)
type RegisteWebHandler func(ctx context.Context, mux *runtime.ServeMux, conn *grpc.ClientConn) error

type Entry struct {
	config            Config
	logger            *zap.Logger
	srv               *grpc.Server
	registeServer     RegisteServer
	registeWebHandler RegisteWebHandler
}

func (e *Entry) Run() error {
	nirvana := interceptor.NirvanaInterceptor{}
	srv := grpc.NewServer(
		grpc.UnaryInterceptor(grpc_middleware.ChainUnaryServer(
			grpcopentracing.UnaryServerInterceptor(grpcopentracing.WithTracer(opentracing.GlobalTracer())),
			nirvana.Serve,
		)))
	if e.registeServer != nil {
		e.registeServer(srv)
		e.srv = srv
	}
	if e.config.EnableReflect {
		log.Print("enable grpc reflect")
		reflection.Register(srv)
	}
	listen, err := net.Listen("tcp", e.config.Port)
	if err != nil {
		log.Printf("listen tcp error: %#v", err)
		return err
	}
	xlog.S(context.Background()).Infof("listening tcp port: %s", e.config.Port)

	go func() {
		log.Println("start grpc server at ", e.config.Port)
		if err = srv.Serve(listen); err != nil {
			log.Fatalf("grpc serve listen error: %#v", err)
		}
	}()
	return nil
}

func (e *Entry) Stop(ctx context.Context) error {
	e.srv.GracefulStop()
	return nil
}
