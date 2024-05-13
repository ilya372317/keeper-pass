package app

import (
	"context"
	"crypto/tls"
	"fmt"

	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/logging"
	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/recovery"
	mygrpc "github.com/ilya372317/pass-keeper/internal/server/handler/grpc"
	v1 "github.com/ilya372317/pass-keeper/internal/server/handler/grpc/v1"
	"github.com/ilya372317/pass-keeper/pkg/logger"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

// StartGRPCServer starting grpc server. Block gorutine before grpc server will stopped by ctx parameter.
func (a *App) StartGRPCServer(ctx context.Context) error {
	tlsCredentials, err := a.loadTLSCredentials()
	if err != nil {
		return fmt.Errorf("cannot load TLS credentials: %w", err)
	}

	server := mygrpc.New(a.conf, a.getUnaryInterceptors(), a.getStreamInterceptors(), tlsCredentials)
	server.RegisterHandler(
		v1.New(a.c.GetDefaultAuthService(),
			a.c.GetDefaultLoginPassService(),
			a.c.GetDefaultCreditCardService(),
			a.c.GetDefaultTextService(),
			a.c.GetDefaultBinaryService(),
			a.c.GetDefaultGeneralDataService(),
		),
	)
	err = server.StartAndListen(ctx)
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

func (a *App) loadTLSCredentials() (credentials.TransportCredentials, error) {
	serverCert, err := tls.LoadX509KeyPair(a.conf.GRPC.TLSCertPath, a.conf.GRPC.TLSKeyPath)
	if err != nil {
		return nil, fmt.Errorf("failed load server cert pair: %w", err)
	}

	config := &tls.Config{
		Certificates: []tls.Certificate{serverCert},
		ClientAuth:   tls.NoClientCert,
		MinVersion:   tls.VersionTLS12,
	}

	return credentials.NewTLS(config), nil
}
