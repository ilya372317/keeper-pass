package app

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"os"

	"github.com/ilya372317/pass-keeper/internal/client/adapter/grpcclient/localstore/tokenstore"
	"github.com/ilya372317/pass-keeper/internal/client/adapter/grpcclient/passcleint"
	"github.com/ilya372317/pass-keeper/internal/client/config"
	"github.com/ilya372317/pass-keeper/internal/client/interceptor"
	"github.com/ilya372317/pass-keeper/internal/client/service/passkeeper"
	pb "github.com/ilya372317/pass-keeper/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/credentials/insecure"
)

type Container struct {
	tokenFile *os.File
	conf      config.Config
}

func NewContainer(conf config.Config, tokenFile *os.File) *Container {
	return &Container{
		conf:      conf,
		tokenFile: tokenFile,
	}
}

func (c *Container) GetDefaultPassKeeperService() *passkeeper.Service {
	return passkeeper.New(c.GetPassClient(), c.GetUserCredStorage())
}

func (c *Container) GetAuthInterceptor() *interceptor.AuthInterceptor {
	return interceptor.NewAuthInterceptor(
		c.GetUserCredStorage(),
		c.conf.GRPC.OpenMethods,
	)
}

func (c *Container) GetUserCredStorage() *tokenstore.TokenStorage {
	return tokenstore.New(c.tokenFile)
}

func (c *Container) GetPassClient() *passcleint.PassClient {
	return passcleint.New(c.GetGRPCClientService())
}

func (c *Container) GetGRPCClientService() pb.PassServiceClient {
	return pb.NewPassServiceClient(c.MustGetGRPCConnection())
}

func (c *Container) MustGetGRPCConnection() *grpc.ClientConn {
	tlsCredentials, err := c.loadTLSCredentials()
	if err != nil {
		panic(err)
	}

	authInterceptor := c.GetAuthInterceptor()
	conn, err := grpc.Dial(c.conf.GRPC.Host,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithChainUnaryInterceptor(authInterceptor.Unary()),
		grpc.WithTransportCredentials(tlsCredentials))
	if err != nil {
		panic(err)
	}

	return conn
}

func (c *Container) loadTLSCredentials() (credentials.TransportCredentials, error) {
	if len(c.conf.GRPC.TLSCertPath) == 0 {
		return credentials.NewClientTLSFromCert(nil, ""), nil
	}

	pemServerCA, err := os.ReadFile(c.conf.GRPC.TLSCertPath)
	if err != nil {
		return nil,
			fmt.Errorf(
				"failed open tls cert. if you connected to server with trusted cert just "+
					"remove cert path in config: %w",
				err,
			)
	}

	certPool := x509.NewCertPool()
	if !certPool.AppendCertsFromPEM(pemServerCA) {
		return nil, fmt.Errorf("failed append certs: %w", err)
	}

	tlsConf := &tls.Config{
		RootCAs:    certPool,
		MinVersion: tls.VersionTLS12,
	}

	return credentials.NewTLS(tlsConf), nil
}
