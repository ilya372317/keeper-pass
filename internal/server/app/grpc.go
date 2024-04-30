package app

import (
	"context"
	"fmt"

	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/logging"
	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/recovery"
	mygrpc "github.com/ilya372317/pass-keeper/internal/server/handler/grpc"
	v1 "github.com/ilya372317/pass-keeper/internal/server/handler/grpc/v1"
	"github.com/ilya372317/pass-keeper/internal/server/logger"
	"google.golang.org/grpc"
)

// StartGRPCServer starting grpc server. Block gorutine before grpc server will stopped by ctx parameter.
func (a *App) StartGRPCServer(ctx context.Context) error {
	server := mygrpc.New(a.conf, a.getUnaryInterceptors(), a.getStreamInterceptors())
	server.RegisterHandler(v1.New(a.c.GetDefaultAuthService(), a.c.GetDefaultDataService()))
	err := server.StartAndListen(ctx)
	if err != nil {
		return fmt.Errorf("failed start grpc server: %w", err)
	}
	return nil
}

func (a *App) getUnaryInterceptors() []grpc.UnaryServerInterceptor {
	return []grpc.UnaryServerInterceptor{
		logging.UnaryServerInterceptor(logger.InterceptorLogger()),
		recovery.UnaryServerInterceptor(),
		a.c.GetAuthInterceptor().Unary(a.c.conf.GRPC.OpenMethods),
	}
}

func (a *App) getStreamInterceptors() []grpc.StreamServerInterceptor {
	return []grpc.StreamServerInterceptor{
		logging.StreamServerInterceptor(logger.InterceptorLogger()),
		recovery.StreamServerInterceptor(),
		a.c.GetAuthInterceptor().Stream(a.c.conf.GRPC.OpenMethods),
	}
}
