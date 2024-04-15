package v1

import (
	"context"
	"fmt"
	"log"
	"net"
	"sync"
	"testing"

	pb "github.com/ilya372317/pass-keeper/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/test/bufconn"
)

var (
	client pb.PassServiceClient
	lis    *bufconn.Listener
)

func TestMain(m *testing.M) {
	ctx := context.Background()
	lis = bufconn.Listen(1024 * 1024)
	conn, err := grpc.DialContext(ctx, "bufnet",
		grpc.WithContextDialer(bufDialer),
		grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		panic(err)
	}
	client = pb.NewPassServiceClient(conn)
	defer func() {
		err = conn.Close()
		if err != nil {
			log.Printf("failed close grpc client conn: %v", err)
		}
	}()

	grpcServer := grpc.NewServer()
	pb.RegisterPassServiceServer(grpcServer, &Handler{})

	wg := sync.WaitGroup{}
	wg.Add(1)
	go func() {
		defer wg.Done()
		err := grpcServer.Serve(lis)
		if err != nil {
			return
		}
	}()
	m.Run()
	grpcServer.GracefulStop()
	wg.Wait()
}

func bufDialer(context.Context, string) (net.Conn, error) {
	conn, err := lis.Dial()
	if err != nil {
		return nil, fmt.Errorf("failed make test connection for grpc: %w", err)
	}

	return conn, nil
}

func TestHandler_Register(t *testing.T) {
}

func TestHandler_Auth(t *testing.T) {
}
